package selector

import (
	"github.com/bvankampen/kubeconfig-selector/internal/config"
	"github.com/urfave/cli"
)

type Selector struct {
	ctx       *cli.Context
	appConfig config.AppConfig
	debug     bool
}
