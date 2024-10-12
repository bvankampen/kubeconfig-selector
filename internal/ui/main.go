package ui

import (
	"github.com/bvankampen/kubeconfig-selector/internal/config"
	"github.com/bvankampen/kubeconfig-selector/internal/kubeconfig"
	"github.com/rivo/tview"
	"github.com/urfave/cli"
)

func (ui *UI) Init(ctx *cli.Context, appConfig config.AppConfig, debug bool) {
	kubeConfigs, activeConfig := kubeconfig.LoadKubeConfigs(appConfig)
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

func (ui *UI) ReloadKubeConfigs() {
	kubeConfigs, activeConfig := kubeconfig.LoadKubeConfigs(ui.appConfig)
	ui.kubeConfigs = kubeConfigs
	ui.activeConfig = activeConfig
}

func (ui *UI) Run() error {
	ui.app.SetRoot(ui.pages, true)
	err := ui.app.Run()
	return err
}
