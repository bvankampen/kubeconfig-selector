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
	ui.pages.AddPage("error", errorMessage, false, true)
}
