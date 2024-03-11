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
		if event.Rune() == 'q' {
			s.app.Stop()
		}
		if event.Rune() == 'k' {
			if s.appConfig.ShowKubeConfig {
				s.appConfig.ShowKubeConfig = false
			} else {
				s.appConfig.ShowKubeConfig = true
			}
			s.reloadScreen()
		}
		if event.Rune() == 'd' {
			if s.debug {
				s.debug = false
			} else {
				s.debug = true
			}
			s.reloadScreen()
		}
		if event.Rune() == 'm' {
			s.moveKubeconfig()
			s.app.Stop()
		}
		return event
	})
}

func (s *Selector) Run() error {
	s.app = tview.NewApplication()
	s.configureInputKeys()
	pages := s.setupPages()
	err := s.app.SetRoot(pages, true).Run()
	if err != nil {
		logrus.Panicf("Error: %v", err)
	}
	return nil
}
