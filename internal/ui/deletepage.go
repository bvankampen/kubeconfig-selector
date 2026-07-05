package ui

import (
	"fmt"

	"github.com/bvankampen/kubeconfig-selector/internal/selector"
	"github.com/rivo/tview"
)

func (ui *UI) deleteCurrentItem() {
	index := ui.list.GetCurrentItem()
	name, _, _ := selector.GetConfigByIndex(ui.kubeConfigs, ui.listEntries, index)
	deleteMessage := tview.NewModal()
	deleteMessage.SetText(fmt.Sprintf("Do you want to delete kubeconfig file for context: %s?", name))
	deleteMessage.AddButtons([]string{"Yes", "No"})
	deleteMessage.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		switch buttonLabel {
		case "Yes":
			err := selector.DeleteConfig(ui.kubeConfigs, ui.listEntries, index, ui.appConfig, ui.activeConfig)
			if err != nil {
				ui.ErrorMessage(fmt.Sprintf("Error deleting kubeconfig: %v", err))
				return
			}
			ui.kubeConfigs = selector.DeleteConfigByIndex(ui.kubeConfigs, ui.listEntries, index)
			ui.redrawLists()
		}
		ui.pages.
			HidePage(pageDelete).
			RemovePage(pageDelete)
	})
	ui.pages.AddPage(pageDelete, deleteMessage, false, true)
}
