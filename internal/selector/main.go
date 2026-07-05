package selector

import (
	"github.com/bvankampen/kubeconfig-selector/internal/config"
	"github.com/bvankampen/kubeconfig-selector/internal/ui"
	"github.com/urfave/cli/v3"
)

func New(cmd *cli.Command) (*Selector, error) {
	appconfig := config.LoadAppConfig()

	config.WriteAppConfig(appconfig) // Write appconfig to update new values.

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
