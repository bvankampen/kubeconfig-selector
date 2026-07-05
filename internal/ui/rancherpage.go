package ui

import (
	"strings"

	"github.com/bvankampen/kubeconfig-selector/internal/rancher"
	"github.com/bvankampen/kubeconfig-selector/internal/selector"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (ui *UI) showDownstreamClusters() {
	index := ui.list.GetCurrentItem()
	_, config, ctx := selector.GetConfigByIndex(ui.kubeConfigs, ui.listEntries, index)

	server := config.Clusters[ctx.Cluster].Server

	if !strings.HasSuffix(server, "local") {
		return
	}

	token := config.AuthInfos[ctx.AuthInfo].Token

	if server == "" || token == "" {
		ui.ShowInfoMessage("Selected context has no valid Rancher server or token.")
		return
	}

	clusters, err := rancher.FetchDownstreamClusters(server, token)
	if err != nil {
		ui.ShowInfoMessage(err.Error())
		return
	}

	if len(clusters) == 0 {
		ui.ShowInfoMessage("No downstream clusters found.")
		return
	}

	list := tview.NewList()
	list.ShowSecondaryText(false)
	list.SetBorder(true).SetTitle("Downstream Cluster List").SetBorderColor(tcell.ColorTeal)
	list.SetHighlightFullLine(true)

	for _, c := range clusters {
		list.AddItem(c.Name, "", 0, nil)
	}

	list.AddItem("Close", "", 'q', nil)

	list.SetSelectedFunc(func(i int, mainText string, secondaryText string, shortcut rune) {
		if mainText == "Close" {
			ui.pages.HidePage("downstream")
			ui.pages.RemovePage("downstream")
			return
		}

		for _, c := range clusters {
			if c.Name == mainText {
				ui.downloadDownstreamKubeConfig(server, token, c)
				return
			}
		}
	})

	_, _, width, height := ui.pages.GetRect()
	listHeight := len(clusters) + 3
	x := (width - 50) / 2
	y := (height - listHeight) / 2
	list.SetRect(x, y, 50, listHeight)

	ui.pages.AddPage("downstream", list, false, true)
}

func (ui *UI) downloadDownstreamKubeConfig(server, token string, cluster rancher.DownstreamCluster) {
	err := selector.DownloadDownstreamKubeConfig(server, token, cluster, ui.appConfig.KubeconfigDir)
	if err != nil {
		ui.ErrorMessage(err.Error())
		return
	}

	ui.pages.HidePage("downstream")
	ui.pages.RemovePage("downstream")
	ui.redrawLists()
}
