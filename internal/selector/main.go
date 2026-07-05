package selector

import (
	"github.com/bvankampen/kubeconfig-selector/internal/config"
	"github.com/bvankampen/kubeconfig-selector/internal/kubeconfig"
	"github.com/urfave/cli/v3"
	"k8s.io/client-go/tools/clientcmd/api"
)

func New(cmd *cli.Command) (*Selector, error) {
	appconfig, err := config.LoadAppConfig()
	if err != nil {
		return nil, err
	}

	if err := config.WriteAppConfig(appconfig); err != nil {
		return nil, err
	}

	kubeConfigs, activeConfig, err := kubeconfig.LoadKubeConfigs(*appconfig)
	if err != nil {
		return nil, err
	}

	return &Selector{
		cmd:          cmd,
		appConfig:    *appconfig,
		kubeConfigs:  kubeConfigs,
		activeConfig: activeConfig,
		debug:        cmd.Bool("debug"),
	}, nil
}

func (s *Selector) Cmd() *cli.Command {
	return s.cmd
}

func (s *Selector) AppConfig() config.AppConfig {
	return s.appConfig
}

func (s *Selector) KubeConfigs() []api.Config {
	return s.kubeConfigs
}

func (s *Selector) ActiveConfig() api.Config {
	return s.activeConfig
}

func (s *Selector) Debug() bool {
	return s.debug
}

func (s *Selector) ReloadKubeConfigs() error {
	kubeConfigs, activeConfig, err := kubeconfig.LoadKubeConfigs(s.appConfig)
	if err != nil {
		return err
	}
	s.kubeConfigs = kubeConfigs
	s.activeConfig = activeConfig
	return nil
}
