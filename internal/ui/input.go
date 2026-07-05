package ui

import (
	"github.com/gdamore/tcell/v2"
)

func (ui *UI) configureInput() {
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		frontPageName, _ := ui.pages.GetFrontPage()
		switch frontPageName {
		case pageHelp:
			switch event.Rune() {
			case 'q':
				ui.pages.HidePage(pageHelp)
				ui.pages.RemovePage(pageHelp)
			}
		case pageDownstream:
			switch event.Rune() {
			case 'q':
				ui.pages.HidePage(pageDownstream)
				ui.pages.RemovePage(pageDownstream)
			}
		case pageMain:
			switch event.Rune() {
			case 'q':
				ui.app.Stop()
			case 'r':
				ui.renameCurrentItem()
				return nil
			case 'k':
				if ui.appConfig.ShowKubeConfig {
					ui.appConfig.ShowKubeConfig = false
				} else {
					ui.appConfig.ShowKubeConfig = true
				}
				ui.redrawAppMain()
			case 'd':
				ui.deleteCurrentItem()
			case 'm':
				ui.moveKubeConfig()
			case 'x':
				ui.showDownstreamClusters()
				return nil
			case '?':
				ui.help()
			}
			if event.Key() == tcell.KeyF5 {
				ui.redrawLists()
			}
		}
		return event
	})
}
