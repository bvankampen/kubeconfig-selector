package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (ui *UI) getHelpText() string {
	return "KS Help \n" +
		"Kubeconfig Selector (ks) Version " + ui.cmd.Version + "\n\n" +
		"  [yellow]q:[white]       Quit \n" +
		"  [yellow]<enter>:[white] Use Kubeconfig \n" +
		"  [yellow]m:[white]       Move Kubeconfig to " + ui.appConfig.KubeconfigDir + " and use it " + "\n" +
		"  [yellow]d:[white]       Delete File \n" +
		"  [yellow]k:[white]       Toggle Kubeconfig \n" +
		"  [yellow]r:[white]       Rename Context \n" +
		"  [yellow]x:[white]       Show Downstream Clusters \n" +
		"  [yellow]F5:[white]      Reload Kubeconfigs \n" +
		"  [yellow]?:[white]       Help \n\n" +
		"  Prefixes:\n" +
		"  [yellow](*):[white]     File not in " + ui.appConfig.KubeconfigDir + "\n" +
		"  [yellow](r):[white]     Rancher Manager context\n" +
		"\n\n(press q to close this screen)"
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
