package selector

import (
	"github.com/bvankampen/kubeconfig-selector/internal/config"
	"github.com/bvankampen/kubeconfig-selector/internal/ui"
	"github.com/urfave/cli/v3"
)

func New(cmd *cli.Command) (*Selector, error) {
	appconfig, err := config.LoadAppConfig()
	if err != nil {
		return nil, err
	}

	if err := config.WriteAppConfig(appconfig); err != nil {
		return nil, err
	}

	return &Selector{
		cmd:       cmd,
		appConfig: *appconfig,
		debug:     cmd.Bool("debug"),
	}, nil
}

func (s *Selector) Run() error {
	var ui ui.UI
	if err := ui.Init(s.cmd, s.appConfig, s.debug); err != nil {
		return err
	}
	return ui.Run()
}
