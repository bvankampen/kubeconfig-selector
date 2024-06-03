package selector

import (
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

type Selector struct {
	ctx           *cli.Context
	appConfig     AppConfig
	kubeConfigs   []api.Config
	activeConfig  api.Config
	app           *tview.Application
	list          *tview.List
	table         *tview.Table
	configView    *tview.TextView
	debugView     *tview.TextView
	helpView      *tview.TextView
	errorMessage  *tview.Modal
	deleteMessage *tview.Modal
	helpMessage   *tview.TextView
	pages         *tview.Pages
	tableRow      int
	tableColumn   int
	configList    []ConfigList
	debug         bool
}
