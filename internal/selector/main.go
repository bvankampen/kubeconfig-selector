package selector

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	// "github.com/davecgh/go-spew/spew"
)

func New(ctx cli.Context) (*Selector, error) {
	appconfig := loadAppConfig()
	kubeconfigs, activeconfig := loadKubeConfigs(appconfig)

	return &Selector{
		ctx:          &ctx,
		appConfig:    *appconfig,
		kubeConfigs:  kubeconfigs,
		activeConfig: activeconfig,
		debug:        ctx.GlobalBool("debug"),
	}, nil
}

func (s *Selector) configureInputKeys() {
	s.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if s.errorMessage.HasFocus() {
			return event
		}

		frontPageName, _ := s.pages.GetFrontPage()
		switch frontPageName {
		case "help":
			switch event.Rune() {
			case 'q':
				s.pages.HidePage("help")
			}
		case "selectorPage":
			switch event.Rune() {
			case 'q':
				s.app.Stop()
			case 'k':
				if s.appConfig.ShowKubeConfig {
					s.appConfig.ShowKubeConfig = false
				} else {
					s.appConfig.ShowKubeConfig = true
				}
				s.reloadScreen()
			case 'd':
				s.deleteCurrentItem()
			case 'v':
				if s.debug {
					s.debug = false
				} else {
					s.debug = true
				}
				s.reloadScreen()
			case 'm':
				s.moveKubeconfig()
				s.app.Stop()
			case '?':
				s.showHelpMessage()
			}
		}
		return event
	})
}

func (s *Selector) Run() error {
	s.app = tview.NewApplication()
	s.configureInputKeys()
	s.setupPages()
	err := s.app.SetRoot(s.pages, true).Run()
	if err != nil {
		logrus.Panicf("Error: %v", err)
	}
	return nil
}
