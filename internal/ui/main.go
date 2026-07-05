package ui

import (
	"github.com/bvankampen/kubeconfig-selector/internal/config"
	"github.com/bvankampen/kubeconfig-selector/internal/kubeconfig"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v3"
	"k8s.io/client-go/tools/clientcmd/api"
)

func (ui *UI) Init(cmd *cli.Command, appConfig config.AppConfig, kubeConfigs []api.Config, activeConfig api.Config, debug bool) error {
	ui.cmd = cmd
	ui.debug = debug
	ui.app = tview.NewApplication()
	ui.pages = tview.NewPages()
	ui.appConfig = appConfig
	ui.kubeConfigs = kubeConfigs
	ui.activeConfig = activeConfig

	ui.configureInput()
	ui.pages.AddPage(pageMain, ui.appPage(), true, true)
	ui.createAppMain()
	return nil
}

func (ui *UI) ReloadKubeConfigs() error {
	kubeConfigs, activeConfig, err := kubeconfig.LoadKubeConfigs(ui.appConfig)
	if err != nil {
		return err
	}
	ui.kubeConfigs = kubeConfigs
	ui.activeConfig = activeConfig
	return nil
}

func (ui *UI) Run() error {
	ui.app.SetRoot(ui.pages, true)
	err := ui.app.Run()
	return err
}
