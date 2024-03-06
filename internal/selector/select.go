package selector

import (
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"strings"
)

func redactConfig(config api.Config) api.Config {
	for _, cluster := range config.Clusters {
		if len(cluster.CertificateAuthorityData) > 0 {
			cluster.CertificateAuthorityData = []byte("REDACTED")
		}
	}
	for _, authInfo := range config.AuthInfos {
		if len(authInfo.ClientCertificateData) > 0 {
			authInfo.ClientCertificateData = []byte("REDACTED")
		}
		if len(authInfo.ClientKeyData) > 0 {
			authInfo.ClientKeyData = []byte("REDACTED")
		}
		if len(authInfo.Token) > 0 {
			authInfo.Token = "[REDACTED]"
		}
	}
	return config
}

func (s *Selector) AddtoTable(field string, value string) {
	s.table.SetCell(s.tableRow, s.tableColumn, tview.NewTableCell(field))
	s.tableColumn += 1
	s.table.SetCell(s.tableRow, s.tableColumn, tview.NewTableCell(value))
	s.tableRow += 1
	s.tableColumn = 0
}

func (s *Selector) AddtoTableList(tableList *ConfigList, field string, value string) {
	tableList.Rows = append(tableList.Rows, TableListItem{
		Field: field,
		Value: value,
	})
}

func (s *Selector) createContextList() {
	s.list.ShowSecondaryText(false)
	s.list.SetBorder(true).SetTitle("Context")
	for _, config := range s.kubeConfigs {
		var tableList ConfigList
		for _, configContext := range config.Contexts {

			s.AddtoTableList(&tableList, "Cluster", configContext.Cluster)
			s.AddtoTableList(&tableList, "File", configContext.LocationOfOrigin)
			s.list.AddItem(configContext.Cluster, "", 0, func() {
				// insert select kubeconfig code
			})
			tableList.Context = configContext
			tableList.Config = config
			tableList.RedactedConfig = redactConfig(config)
		}
		s.configList = append(s.configList, tableList)
	}
	s.list.AddItem("Quit", "", 'q', func() { s.app.Stop() })

	s.list.SetChangedFunc(func(index int, mainText string, secondayText string, shortcut rune) {
		s.tableRow = 0
		s.tableColumn = 0
		if index < len(s.configList) {
			for _, item := range s.configList[index].Rows {
				s.AddtoTable(item.Field, item.Value)
			}

			configBytes, _ := clientcmd.Write(s.configList[index].RedactedConfig)
			s.configView.SetText(strings.Replace(string(configBytes), "UkVEQUNURUQ=", "[REDACTED]", -1))

		} else {
			s.table.Clear()
			s.configView.Clear()
		}

	})

	for _, item := range s.configList[0].Rows {
		s.AddtoTable(item.Field, item.Value)
	}

}

func (s *Selector) selectKubeconfig() {
	s.app = tview.NewApplication()

	s.list = tview.NewList()
	s.configView = tview.NewTextView()

	s.table = tview.NewTable()

	s.table.SetBorder(true).SetTitle("Cluster")
	s.configView.SetBorder(true).SetTitle("Kubeconfig")

	s.createContextList()

	flexSub := tview.NewFlex().SetDirection(tview.FlexRow)
	flexSub.AddItem(s.table, 0, 1, false)
	if s.appConfig.ShowKubeConfig {
		flexSub.AddItem(s.configView, 0, 2, false)
	}
	flexMain := tview.NewFlex()
	flexMain.AddItem(s.list, 0, 1, true)
	flexMain.AddItem(flexSub, 0, 2, false)

	pages := tview.NewPages().AddPage("selectorPage", flexMain, true, true)
	err := s.app.SetRoot(pages, true).Run()
	if err != nil {
		logrus.Panicf("Error: %v", err)
	}

}
