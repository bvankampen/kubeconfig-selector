package ui

import (
	"github.com/bvankampen/kubeconfig-selector/internal/config"
	"github.com/bvankampen/kubeconfig-selector/internal/selector"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v3"
	"k8s.io/client-go/tools/clientcmd/api"
)

type UI struct {
	cmd          *cli.Command
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
	listEntries  []selector.ListEntry
}
