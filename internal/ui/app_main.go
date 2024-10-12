package ui

import (
	"fmt"
	"path/filepath"
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

func (ui *UI) createList() int {
	ui.list = tview.NewList()
	index := 0
	currentIndex := 0
	ui.list.ShowSecondaryText(false)
	ui.list.SetBorder(true).SetTitle("Context").SetBorderColor(tcell.ColorBlue)
	ui.list.SetHighlightFullLine(true)
	for _, config := range ui.kubeConfigs {
		for name, configContext := range config.Contexts {
			kubeDir, _ := homedir.Expand(ui.appConfig.KubeconfigDir)

			var star rune
			star = 0
			if !strings.HasPrefix(configContext.LocationOfOrigin, kubeDir) {
				star = '*'
			}

			ui.list.AddItem(name, "", star, nil)

			if ui.activeConfig.CurrentContext != "" {
				activeConfigContext := ui.activeConfig.Contexts[ui.activeConfig.CurrentContext]
				activeConfigCluster := activeConfigContext.Cluster
				activeConfigServer := ui.activeConfig.Clusters[activeConfigContext.Cluster].Server

				if configContext.Cluster == activeConfigCluster &&
					config.Clusters[configContext.Cluster].Server == activeConfigServer &&
					name == ui.activeConfig.CurrentContext {
					currentIndex = index
				}
			}

			index++
		}
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
	ui.list.Clear()
	for _, config := range ui.kubeConfigs {
		for name, configContext := range config.Contexts {
			kubeDir, _ := homedir.Expand(ui.appConfig.KubeconfigDir)

			var star rune
			star = 0
			if !strings.HasPrefix(configContext.LocationOfOrigin, kubeDir) {
				star = '*'
			}

			ui.list.AddItem(name, "", star, nil)

			if ui.activeConfig.CurrentContext != "" {
				activeConfigContext := ui.activeConfig.Contexts[ui.activeConfig.CurrentContext]
				activeConfigCluster := activeConfigContext.Cluster
				activeConfigServer := ui.activeConfig.Clusters[activeConfigContext.Cluster].Server

				if configContext.Cluster == activeConfigCluster &&
					config.Clusters[configContext.Cluster].Server == activeConfigServer &&
					name == ui.activeConfig.CurrentContext {
				}
			}

			index++
		}
	}
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
	ui.kubeConfigs = append(ui.kubeConfigs[:index], ui.kubeConfigs[index+1:]...)
}

func (ui *UI) getConfigByIndex(index int) (string, api.Config, *api.Context) {
	contextName, _ := ui.list.GetItemText(index)
	for _, config := range ui.kubeConfigs {
		for name, context := range config.Contexts {
			if name == contextName {
				return name, config, context
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
	addtoTable(infoTable, "Context", name)
	addtoTable(infoTable, "Cluster", context.Cluster)
	addtoTable(infoTable, "User", context.AuthInfo)

	addtoTable(infoTable, "Server", config.Clusters[context.Cluster].Server)
	addtoTable(infoTable, "File", context.LocationOfOrigin)
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

func (ui *UI) printDebug(debugMessage string) {
	if ui.debug {
		debugMessage = fmt.Sprintf("%s%s\n", ui.debugView.GetText(false), debugMessage)
		ui.debugView.SetText(debugMessage)
	}
}

func (ui *UI) createAppMain() {
	currentIndex := ui.createList()
	ui.createViews(false)
	ui.mainFlex.AddItem(ui.list, 0, 1, true)
	ui.mainFlex.AddItem(ui.views, 0, 2, false)
	ui.list.SetCurrentItem(currentIndex)
	ui.redrawAppMain()
}

func (ui *UI) redrawAppMain() {
	ui.createViews(true)
}

func (ui *UI) redrawLists() {
	ui.ReloadKubeConfigs()
	ui.redrawList()
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
	}
	err = kubeconfig.MoveKubeConfig(config.DeepCopy(), name, ui.appConfig.KubeconfigDir)
	if err != nil {
		ui.ErrorMessage(err.Error())
	}
	ui.app.Stop()
}

func (ui *UI) renameKubeConfigContext(index int, config api.Config, contextName string, newContextName string) {
	if contextName != newContextName {
		for name, context := range config.Contexts {
			if name == contextName {
				kubeConfigPath := filepath.Dir(context.LocationOfOrigin)
				kubeConfigFilename := filepath.Base(context.LocationOfOrigin)
				config.Contexts[newContextName] = context.DeepCopy()
				config.CurrentContext = newContextName
				delete(config.Contexts, contextName)
				ui.kubeConfigs[index] = config
				err := kubeconfig.SaveKubeConfig(
					config.DeepCopy(),
					newContextName,
					kubeConfigPath,
					kubeConfigFilename,
					false,
					ui.appConfig.CreateLink,
					false)
				if err != nil {
					ui.ErrorMessage(err.Error())
				}
				return
			}
		}
	}
}
