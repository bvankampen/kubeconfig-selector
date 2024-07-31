package ui

import (
	"github.com/bvankampen/kubeconfig-selector/internal/config"
	"github.com/rivo/tview"
	"github.com/urfave/cli"
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

type UI struct {
	ctx          *cli.Context
	debug        bool
	app          *tview.Application
	list         *tview.List
	views        *tview.Flex
	pages        *tview.Pages
	mainFlex     *tview.Flex
	kubeConfigs  []api.Config
	activeConfig api.Config
	appConfig    config.AppConfig
	debugView    *tview.TextView
}
