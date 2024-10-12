package ui

import (
	"fmt"

	"github.com/bvankampen/kubeconfig-selector/internal/kubeconfig"
	"github.com/rivo/tview"
)

func (ui *UI) deleteCurrentItem() {
	index := ui.list.GetCurrentItem()
	name, config, _ := ui.getConfigByIndex(index)
	deleteMessage := tview.NewModal()
	deleteMessage.SetText(fmt.Sprintf("Do you want to delete context: %s", name))
	deleteMessage.AddButtons([]string{"Yes", "No"})
	deleteMessage.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		switch buttonLabel {
		case "Yes":
			kubeconfig.DeleteKubeConfig(
				config.DeepCopy(),
				name,
				ui.appConfig.KubeconfigDir,
				ui.appConfig.KubeconfigFile,
				ui.appConfig.CreateLink,
			)
			ui.deleteConfigByIndex(index)
			ui.redrawLists()
			// ui.list.RemoveItem(index)
		}
		ui.pages.
			HidePage("delete").
			RemovePage("delete")
	})
	ui.pages.AddPage("delete", deleteMessage, false, true)
}
