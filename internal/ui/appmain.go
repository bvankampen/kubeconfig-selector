package ui

import (
	"fmt"

	"github.com/bvankampen/kubeconfig-selector/internal/selector"
	"github.com/gdamore/tcell/v2"
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
	ui.list.ShowSecondaryText(false)
	ui.list.SetBorder(true).SetTitle("Context").SetBorderColor(tcell.ColorBlue)
	ui.list.SetHighlightFullLine(true)

	ui.listEntries = selector.BuildSortedEntries(ui.kubeConfigs, ui.appConfig.KubeconfigDir)
	for _, entry := range ui.listEntries {
		ui.list.AddItem(entry.Name, "", entry.PrefixSymbol, nil)
	}

	currentIndex := selector.FindActiveIndex(ui.kubeConfigs, ui.listEntries, ui.activeConfig)

	ui.list.SetChangedFunc(func(index int, mainText string, secondayText string, shortcut rune) {
		ui.redrawAppMain()
	})

	ui.list.SetSelectedFunc(func(index int, mainText string, secondayText string, shortcut rune) {
		ui.selectKubeConfig(index)
	})

	return currentIndex
}

func (ui *UI) redrawList() {
	ui.list.Clear()

	ui.listEntries = selector.BuildSortedEntries(ui.kubeConfigs, ui.appConfig.KubeconfigDir)
	for _, entry := range ui.listEntries {
		ui.list.AddItem(entry.Name, "", entry.PrefixSymbol, nil)
	}

	currentIndex := selector.FindActiveIndex(ui.kubeConfigs, ui.listEntries, ui.activeConfig)
	ui.list.SetCurrentItem(currentIndex)
}

func (ui *UI) selectKubeConfig(index int) {
	err := selector.SelectConfig(ui.kubeConfigs, ui.listEntries, index, ui.appConfig)
	if err != nil {
		ui.ErrorMessage(err.Error())
	} else {
		ui.app.Stop()
	}
}

func (ui *UI) getConfigByIndex(index int) (string, api.Config, *api.Context) {
	return selector.GetConfigByIndex(ui.kubeConfigs, ui.listEntries, index)
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
	err := selector.MoveConfig(ui.kubeConfigs, ui.listEntries, index, ui.appConfig)
	if err != nil {
		ui.ErrorMessage(err.Error())
	} else {
		ui.app.Stop()
	}
}

func (ui *UI) renameKubeConfigContext(config api.Config, contextName string, newContextName string) {
	err := selector.RenameContext(ui.kubeConfigs, config, contextName, newContextName, ui.appConfig.KubeconfigDir, ui.appConfig.KubeconfigFile, ui.appConfig.CreateLink)
	if err != nil {
		ui.ErrorMessage(err.Error())
		return
	}
}
