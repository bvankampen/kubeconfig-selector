package selector

import (
	"github.com/bvankampen/kubeconfig-selector/internal/config"
	"github.com/bvankampen/kubeconfig-selector/internal/kubeconfig"
	"github.com/bvankampen/kubeconfig-selector/internal/ui"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	// "github.com/davecgh/go-spew/spew"
)

func New(ctx cli.Context) (*Selector, error) {
	appconfig := config.LoadAppConfig()
	kubeconfigs, activeconfig := kubeconfig.LoadKubeConfigs(appconfig)

	return &Selector{
		ctx:          &ctx,
		appConfig:    *appconfig,
		kubeConfigs:  kubeconfigs,
		activeConfig: activeconfig,
		debug:        ctx.GlobalBool("debug"),
	}, nil
}

func (s *Selector) Run() error {
	var ui ui.UI
	ui.Init(s.ctx, s.appConfig, s.kubeConfigs, s.activeConfig, s.debug)
	err := ui.Run()
	if err != nil {
		logrus.Panicf("Error: %v", err)
	}
	return nil
}
