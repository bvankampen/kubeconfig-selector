package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (ui *UI) ErrorMessage(errorText string) {
	errorMessage := tview.NewModal()
	errorMessage.SetText(fmt.Sprintf("Error: \n\n%s", errorText))
	errorMessage.AddButtons([]string{"Quit"})
	errorMessage.SetBackgroundColor(tcell.ColorRed)
	errorMessage.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Quit" {
			ui.app.Stop()
		}
	})
	errorMessage.SetTitle("Error")
	ui.pages.AddPage(pageError, errorMessage, false, true)
}

func (ui *UI) ShowInfoMessage(infoText string) {
	msg := tview.NewModal()
	msg.SetText(fmt.Sprintf("%s", infoText))
	msg.AddButtons([]string{"OK"})
	msg.SetBackgroundColor(tcell.ColorRed)
	msg.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		ui.pages.HidePage(pageInfo)
		ui.pages.RemovePage(pageInfo)
	})
	msg.SetTitle("Error")
	ui.pages.AddPage(pageInfo, msg, false, true)
}
