package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bvankampen/kubeconfig-selector/internal/kubeconfig"
	"github.com/gdamore/tcell/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/rivo/tview"
	"k8s.io/client-go/tools/clientcmd/api"
)

func addtoTable(table *tview.Table, field string, value string) {
	row := table.GetRowCount()
	table.SetCell(row, 0, tview.NewTableCell(field).SetTextColor(tcell.ColorOrange))
	table.SetCell(row, 1, tview.NewTableCell(value))
}

type listEntry struct {
	name         string
	sourceFile   string
	prefixSymbol rune
}

func (ui *UI) buildSortedEntries() []listEntry {
	seen := make(map[string]bool)
	var entries []listEntry
	kubeDir, _ := homedir.Expand(ui.appConfig.KubeconfigDir)

	for _, cfg := range ui.kubeConfigs {
		for name, cfgContext := range cfg.Contexts {
			if seen[name] {
				continue
			}
			seen[name] = true

			var prefixSymbol rune
			if !strings.HasPrefix(cfgContext.LocationOfOrigin, kubeDir) {
				prefixSymbol = '*'
			} else if cluster, ok := cfg.Clusters[cfgContext.Cluster]; ok {
				if strings.HasSuffix(cluster.Server, "local") {
					prefixSymbol = 'r'
				}
			}

			entries = append(entries, listEntry{name: name, sourceFile: cfgContext.LocationOfOrigin, prefixSymbol: prefixSymbol})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].name < entries[j].name
	})
	return entries
}

func (ui *UI) createList() int {
	ui.list = tview.NewList()
	index := 0
	currentIndex := 0
	ui.list.ShowSecondaryText(false)
	ui.list.SetBorder(true).SetTitle("Context").SetBorderColor(tcell.ColorBlue)
	ui.list.SetHighlightFullLine(true)

	ui.listEntries = ui.buildSortedEntries()
	for _, entry := range ui.listEntries {
		ui.list.AddItem(entry.name, "", entry.prefixSymbol, nil)

		if ui.activeConfig.CurrentContext != "" {
			for _, cfg := range ui.kubeConfigs {
				if cfgContext, ok := cfg.Contexts[entry.name]; ok {
					activeConfigContext := ui.activeConfig.Contexts[ui.activeConfig.CurrentContext]
					activeConfigCluster := activeConfigContext.Cluster
					activeConfigServer := ui.activeConfig.Clusters[activeConfigContext.Cluster].Server

					if cfgContext.Cluster == activeConfigCluster &&
						cfg.Clusters[cfgContext.Cluster].Server == activeConfigServer &&
						entry.name == ui.activeConfig.CurrentContext {
						currentIndex = index
					}
					break
				}
			}
		}

		index++
	}

	ui.list.SetChangedFunc(func(index int, mainText string, secondayText string, shortcut rune) {
		ui.redrawAppMain()
	})

	ui.list.SetSelectedFunc(func(index int, mainText string, secondayText string, shortcut rune) {
		ui.selectKubeConfig(index)
	})

	return currentIndex
}

func (ui *UI) redrawList() {
	index := 0
	currentIndex := 0
	ui.list.Clear()

	ui.listEntries = ui.buildSortedEntries()
	for _, entry := range ui.listEntries {
		ui.list.AddItem(entry.name, "", entry.prefixSymbol, nil)

		if ui.activeConfig.CurrentContext != "" {
			for _, cfg := range ui.kubeConfigs {
				if cfgContext, ok := cfg.Contexts[entry.name]; ok {
					activeConfigContext := ui.activeConfig.Contexts[ui.activeConfig.CurrentContext]
					activeConfigCluster := activeConfigContext.Cluster
					activeConfigServer := ui.activeConfig.Clusters[activeConfigContext.Cluster].Server

					if cfgContext.Cluster == activeConfigCluster &&
						cfg.Clusters[cfgContext.Cluster].Server == activeConfigServer &&
						entry.name == ui.activeConfig.CurrentContext {
						currentIndex = index
					}
					break
				}
			}
		}

		index++
	}
	ui.list.SetCurrentItem(currentIndex)
}

func (ui *UI) selectKubeConfig(index int) {
	name, config, _ := ui.getConfigByIndex(index)
	err := kubeconfig.SaveKubeConfig(
		config.DeepCopy(),
		name,
		ui.appConfig.KubeconfigDir,
		ui.appConfig.KubeconfigFile,
		true,
		ui.appConfig.CreateLink,
		false)
	if err != nil {
		ui.ErrorMessage(err.Error())
	} else {
		ui.app.Stop()
	}
}

func (ui *UI) deleteConfigByIndex(index int) {
	if index >= len(ui.listEntries) {
		return
	}
	entry := ui.listEntries[index]
	for i, config := range ui.kubeConfigs {
		if ctx, ok := config.Contexts[entry.name]; ok {
			if ctx.LocationOfOrigin == entry.sourceFile {
				ui.kubeConfigs = append(ui.kubeConfigs[:i], ui.kubeConfigs[i+1:]...)
				return
			}
		}
	}
}

func (ui *UI) getConfigByIndex(index int) (string, api.Config, *api.Context) {
	if index >= len(ui.listEntries) {
		return "", api.Config{}, &api.Context{}
	}
	entry := ui.listEntries[index]
	for _, config := range ui.kubeConfigs {
		if ctx, ok := config.Contexts[entry.name]; ok {
			if ctx.LocationOfOrigin == entry.sourceFile {
				return entry.name, config, ctx
			}
		}
	}
	return "", api.Config{}, &api.Context{}
}

func (ui *UI) createInfoTable() *tview.Table {
	infoTable := tview.NewTable()
	infoTable.SetBorder(true)
	infoTable.SetTitle("Cluster")
	name, config, context := ui.getConfigByIndex(ui.list.GetCurrentItem())

	if name != "" {
		addtoTable(infoTable, "Context", name)
		cluster, found := config.Clusters[context.Cluster]
		if found {
			addtoTable(infoTable, "Cluster", context.Cluster)
			addtoTable(infoTable, "User", context.AuthInfo)
			addtoTable(infoTable, "Server", cluster.Server)
			addtoTable(infoTable, "File", context.LocationOfOrigin)
		} else {
			addtoTable(infoTable, "Cluster", "Cluster not found")
		}
	}
	return infoTable
}

func (ui *UI) createConfigTextView() *tview.TextView {
	configTextView := tview.NewTextView()
	configTextView.SetBorder(true)
	configTextView.SetTitle("Kubeconfig")
	_, config, _ := ui.getConfigByIndex(ui.list.GetCurrentItem())
	configTextView.SetText(redactConfigToString(*config.DeepCopy()))
	return configTextView
}

func (ui *UI) createViews(redraw bool) {
	if redraw {
		ui.views.Clear()
	} else {
		ui.views = tview.NewFlex().SetDirection(tview.FlexRow)
	}
	tableSize := 0
	if ui.appConfig.ShowKubeConfig {
		tableSize = 7
	}
	ui.views.AddItem(ui.createInfoTable(), tableSize, 1, false)
	if ui.appConfig.ShowKubeConfig {
		ui.views.AddItem(ui.createConfigTextView(), 0, 2, false)
	}
	if ui.debug {
		ui.debugView = tview.NewTextView()
		ui.debugView.SetBorder(true).SetTitle("Debug")
		ui.views.AddItem(ui.debugView, 10, 3, false)
	}
}

// func (ui *UI) printDebug(debugMessage string) {
// 	if ui.debug {
// 		debugMessage = fmt.Sprintf("%s%s\n", ui.debugView.GetText(false), debugMessage)
// 		ui.debugView.SetText(debugMessage)
// 	}
// }

func (ui *UI) createAppMain() {
	currentIndex := ui.createList()
	ui.createViews(false)
	ui.mainFlex.AddItem(ui.list, 0, 1, true)
	ui.mainFlex.AddItem(ui.views, 0, 2, false)
	ui.list.SetCurrentItem(currentIndex)
	ui.redrawAppMain()
	if ui.list.GetItemCount() == 0 {
		ui.ErrorMessage("No (other) configs found, nothing to choose from.")
	}
}

func (ui *UI) redrawAppMain() {
	ui.createViews(true)
}

func (ui *UI) redrawLists() {
	if err := ui.ReloadKubeConfigs(); err != nil {
		ui.ErrorMessage(fmt.Sprintf("Error reloading kubeconfigs: %v", err))
		return
	}
	ui.redrawList()
	ui.redrawAppMain()
}

func (ui *UI) moveKubeConfig() {
	index := ui.list.GetCurrentItem()
	name, config, _ := ui.getConfigByIndex(index)
	err := kubeconfig.SaveKubeConfig(
		config.DeepCopy(),
		name,
		ui.appConfig.KubeconfigDir,
		ui.appConfig.KubeconfigFile,
		true,
		ui.appConfig.CreateLink,
		true)
	if err != nil {
		ui.ErrorMessage(err.Error())
	} else {
		ui.app.Stop()
	}
}

func (ui *UI) renameKubeConfigContext(config api.Config, contextName string, newContextName string) {
	if contextName == newContextName {
		return
	}
	for _, cfg := range ui.kubeConfigs {
		if _, exists := cfg.Contexts[newContextName]; exists {
			ui.ShowInfoMessage(fmt.Sprintf("Context %q already exists.", newContextName))
			return
		}
	}
	for name, context := range config.Contexts {
		if name == contextName {
			kubeConfigPath := filepath.Dir(context.LocationOfOrigin)
			kubeConfigFilename := filepath.Base(context.LocationOfOrigin)
			config.Contexts[newContextName] = context.DeepCopy()
			config.CurrentContext = newContextName
			delete(config.Contexts, contextName)
			err := kubeconfig.SaveKubeConfigFile(
				config.DeepCopy(),
				newContextName,
				kubeConfigPath,
				kubeConfigFilename,
			)
			if err != nil {
				ui.ErrorMessage(err.Error())
				return
			}

			ext := filepath.Ext(kubeConfigFilename)
			newFilename := newContextName + ext
			oldPath := filepath.Join(kubeConfigPath, kubeConfigFilename)
			newPath := filepath.Join(kubeConfigPath, newFilename)
			err = os.Rename(oldPath, newPath)
			if err != nil {
				ui.ErrorMessage(err.Error())
				return
			}
			context.LocationOfOrigin = newPath

			for i, cfg := range ui.kubeConfigs {
				for n := range cfg.Contexts {
					if n == contextName {
						ui.kubeConfigs[i] = config
						return
					}
				}
			}
			return
		}
	}
}
