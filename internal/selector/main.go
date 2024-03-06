package selector

import (
	"context"
	"github.com/rivo/tview"
	"k8s.io/client-go/tools/clientcmd/api"
)

type TableListItem struct {
	Field string
	Value string
}

type ConfigList struct {
	Rows           []TableListItem
	Config         api.Config
	RedactedConfig api.Config
	Context        *api.Context
}

type Selector struct {
	ctx          context.Context
	appConfig    AppConfig
	kubeConfigs  []api.Config
	activeConfig api.Config
	app          *tview.Application
	list         *tview.List
	table        *tview.Table
	configView   *tview.TextView
	tableRow     int
	tableColumn  int
	configList   []ConfigList
}

func New(ctx context.Context) (*Selector, error) {

	appconfig := loadAppConfig()
	kubeconfigs, activeconfig := loadKubeConfigs(appconfig)

	return &Selector{
		ctx:          ctx,
		appConfig:    *appconfig,
		kubeConfigs:  kubeconfigs,
		activeConfig: activeconfig,
	}, nil

}

func (s *Selector) Run() error {
	s.selectKubeconfig()

	return nil

}
