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

	clusters, err := rancher.FetchDownstreamClusters(server, token, false)
	if err != nil {
		if rancher.IsCertError(err) {
			ui.showCertConfirm(func(insecure bool) {
				if !insecure {
					ui.ShowInfoMessage(err.Error())
					return
				}
				clusters, err := rancher.FetchDownstreamClusters(server, token, true)
				if err != nil {
					ui.ShowInfoMessage(err.Error())
					return
				}
				ui.showDownList(clusters)
			})
			return
		}
		ui.ShowInfoMessage(err.Error())
		return
	}

	ui.showDownList(clusters)
}

func (ui *UI) showDownList(clusters []rancher.DownstreamCluster) {
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
			ui.pages.HidePage(pageDownstream)
			ui.pages.RemovePage(pageDownstream)
			return
		}

		for _, c := range clusters {
			if c.Name == mainText {
				ui.downloadDownstreamKubeConfig(c)
				return
			}
		}
	})

	_, _, width, height := ui.pages.GetRect()
	listHeight := len(clusters) + 3
	x := (width - 50) / 2
	y := (height - listHeight) / 2
	list.SetRect(x, y, 50, listHeight)

	ui.pages.AddPage(pageDownstream, list, false, true)
}

func (ui *UI) downloadDownstreamKubeConfig(cluster rancher.DownstreamCluster) {
	index := ui.list.GetCurrentItem()
	_, config, ctx := selector.GetConfigByIndex(ui.kubeConfigs, ui.listEntries, index)
	server := config.Clusters[ctx.Cluster].Server
	token := config.AuthInfos[ctx.AuthInfo].Token

	err := selector.DownloadDownstreamKubeConfig(server, token, cluster, ui.appConfig.KubeconfigDir, false)
	if err != nil {
		if rancher.IsCertError(err) {
			ui.showCertConfirm(func(insecure bool) {
				if !insecure {
					ui.ErrorMessage(err.Error())
					return
				}
				err := selector.DownloadDownstreamKubeConfig(server, token, cluster, ui.appConfig.KubeconfigDir, true)
				if err != nil {
					ui.ErrorMessage(err.Error())
					return
				}
				ui.pages.HidePage(pageDownstream)
				ui.pages.RemovePage(pageDownstream)
				ui.redrawLists()
			})
			return
		}
		ui.ErrorMessage(err.Error())
		return
	}

	ui.pages.HidePage(pageDownstream)
	ui.pages.RemovePage(pageDownstream)
	ui.redrawLists()
}

func (ui *UI) showCertConfirm(onAccept func(bool)) {
	modal := tview.NewModal()
	modal.SetText("Self-signed or invalid certificate detected.\nAccept connection anyway?")
	modal.AddButtons([]string{"Accept", "Cancel"})
	modal.SetBackgroundColor(tcell.ColorOrange)
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		ui.pages.HidePage(pageCertConfirm)
		ui.pages.RemovePage(pageCertConfirm)
		if buttonLabel == "Accept" {
			onAccept(true)
		} else {
			onAccept(false)
		}
	})
	modal.SetTitle("Certificate Error")
	ui.pages.AddPage(pageCertConfirm, modal, false, true)
}
