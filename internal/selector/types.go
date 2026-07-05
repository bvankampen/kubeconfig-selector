package selector

import (
	"github.com/bvankampen/kubeconfig-selector/internal/config"
	"github.com/urfave/cli/v3"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Selector struct {
	cmd          *cli.Command
	appConfig    config.AppConfig
	kubeConfigs  []api.Config
	activeConfig api.Config
	debug        bool
}
