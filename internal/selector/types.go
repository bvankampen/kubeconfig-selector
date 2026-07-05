package selector

import (
	"github.com/bvankampen/kubeconfig-selector/internal/config"
	"github.com/urfave/cli/v3"
)

type Selector struct {
	cmd       *cli.Command
	appConfig config.AppConfig
	debug     bool
}
