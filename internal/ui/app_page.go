package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createHeader(version string) *tview.TextView {
	header := tview.NewTextView()
	header.SetBackgroundColor(tcell.ColorDarkCyan)
	header.SetTextColor(tcell.ColorBlack)
	header.SetText(fmt.Sprintf(" Kubeconfig Selector %s https://github.com/bvankampen/kubeconfig-selector", version))
	return header
}

func createFooter(kubeconfigDir string) *tview.TextView {
	footer := tview.NewTextView() // 	s.helpView.SetBorder(false)
	footer.SetDynamicColors(true)
	footerText := "[yellow]q:[white] Quit " +
		"[yellow]<enter>:[white] Use Kubeconfig " +
		"[yellow]m:[white] Move Kubeconfig to " + kubeconfigDir + " and use it " +
		"[yellow]d:[white] Delete file " +
		"[yellow]r:[white] Rename context " +
		"[yellow](*):[white] File not in " + kubeconfigDir + " " +
		"[yellow]?:[white]help"
	footer.SetText(footerText)
	return footer
}

func (ui *UI) appPage() *tview.Flex {
	ui.mainFlex = tview.NewFlex()
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(createHeader(ui.ctx.App.Version), 1, 1, false)
	flex.AddItem(ui.mainFlex, 0, 1, true)
	flex.AddItem(createFooter(ui.appConfig.KubeconfigDir), 1, 1, false)

	return flex
}
