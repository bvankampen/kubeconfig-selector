package ui

import (
	"github.com/gdamore/tcell/v2"
)

func (ui *UI) configureInput() {
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		frontPageName, _ := ui.pages.GetFrontPage()
		switch frontPageName {
		case "help":
			switch event.Rune() {
			case 'q':
				ui.pages.HidePage("help")
				ui.pages.RemovePage("help")
			}
		case "main":
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
			case '?':
				ui.help()
			}
		}
		return event
	})
}
