package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (ui *UI) getHelpText() string {
	return "KS Help \n" +
		"Kubernetes Selector (ks) Version " + ui.ctx.App.Version + "\n\n" +
		"  [yellow]q:[white]       Quit \n" +
		"  [yellow]<enter>:[white] Use Kubeconfig \n" +
		"  [yellow]m:[white]       Move Kubeconfig to " + ui.appConfig.KubeconfigDir + " and use it " + "\n" +
		"  [yellow]d:[white]       Delete File \n" +
		"  [yellow]k:[white]       Toggle Kubeconfig \n" +
		"  [yellow]r:[white]	   Rename Context \n" +
		"  [yellow]x:[white]       Set kubeconfig as Rancher Manager kubeconfig\n" +
		"  [yellow](*):[white]     File not in " + ui.appConfig.KubeconfigDir + "\n" +
		"  [yellow]?:[white]       Help \n" +
		"\n\n\n(press q to close this screen)"
}

func (ui *UI) help() {
	help := tview.NewTextView()
	help.SetDynamicColors(true)
	help.SetBorderColor(tcell.ColorTeal)
	help.SetBorder(true)
	help.SetTitle("Help")
	help.SetText(ui.getHelpText())
	ui.pages.AddPage("help", help, true, true)
}
