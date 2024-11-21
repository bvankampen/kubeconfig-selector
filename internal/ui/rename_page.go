package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (ui *UI) renameCurrentItem() {
	index := ui.list.GetCurrentItem()
	name, config, _ := ui.getConfigByIndex(index)
	renameForm := tview.NewForm()
	renameForm.SetBorder(true)
	renameForm.SetTitle("Rename")
	renameForm.AddInputField("Context", name, 27, nil, nil)
	renameForm.SetBackgroundColor(tcell.ColorDarkCyan)
	renameForm.SetFieldTextColor(tcell.ColorWhite)
	renameForm.SetFieldBackgroundColor(tcell.ColorBlack)
	renameForm.SetButtonBackgroundColor(tcell.ColorBlack)
	renameForm.SetButtonTextColor(tcell.ColorWhite)
	renameForm.SetButtonsAlign(tview.AlignCenter)
	renameForm.AddButton("Rename", func() {
		newName := renameForm.GetFormItemByLabel("Context").(*tview.InputField).GetText()
		ui.renameKubeConfigContext(
			index,
			*config.DeepCopy(),
			name,
			newName,
		)
		ui.list.SetItemText(index, newName, "")

		ui.redrawAppMain()

		ui.pages.
			HidePage("rename").
			RemovePage("rename")
	})
	renameForm.AddButton("Cancel", func() {
		ui.pages.
			HidePage("rename").
			RemovePage("rename")
	})

	// Center renameForm on screen
	_, _, width, _ := ui.pages.GetRect() // get screen width
	x := ((width) / 2) - 20
	renameForm.SetRect(x, 5, 40, 7)

	ui.pages.AddPage("rename", renameForm, false, true)
}
