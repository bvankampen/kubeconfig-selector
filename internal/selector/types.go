package selector

import (
	"github.com/bvankampen/kubeconfig-selector/internal/config"
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
	ctx          *cli.Context
	appConfig    config.AppConfig
	kubeConfigs  []api.Config
	activeConfig api.Config
	debug        bool
}
