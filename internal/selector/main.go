package selector

import (
	"github.com/bvankampen/kubeconfig-selector/internal/config"
	"github.com/bvankampen/kubeconfig-selector/internal/ui"
	"github.com/urfave/cli"
)

func New(ctx cli.Context) (*Selector, error) {
	appconfig := config.LoadAppConfig()

	config.WriteAppConfig(appconfig) // Write appconfig to update new values.

	return &Selector{
		ctx:       &ctx,
		appConfig: *appconfig,
		debug:     ctx.GlobalBool("debug"),
	}, nil
}

func (s *Selector) Run() error {
	var ui ui.UI
	if err := ui.Init(s.ctx, s.appConfig, s.debug); err != nil {
		return err
	}
	return ui.Run()
}
