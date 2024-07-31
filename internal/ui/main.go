package ui

import (
	"github.com/bvankampen/kubeconfig-selector/internal/config"
	"github.com/rivo/tview"
	"github.com/urfave/cli"
	"k8s.io/client-go/tools/clientcmd/api"
)

func (ui *UI) Init(ctx *cli.Context, appConfig config.AppConfig, kubeConfigs []api.Config, activeConfig api.Config, debug bool) {
	ui.ctx = ctx
	ui.debug = debug
	ui.app = tview.NewApplication()
	ui.pages = tview.NewPages()
	ui.appConfig = appConfig
	ui.kubeConfigs = kubeConfigs
	ui.activeConfig = activeConfig

	ui.configureInput()
	ui.pages.AddPage("main", ui.appPage(), true, true)
	ui.createAppMain()
}

func (ui *UI) Run() error {
	ui.app.SetRoot(ui.pages, true)
	err := ui.app.Run()
	return err
}
