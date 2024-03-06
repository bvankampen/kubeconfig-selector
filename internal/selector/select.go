package selector

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"strconv"
)

func (s *Selector) AddtoTable(field string, value string) {
	s.table.SetCell(s.tableRow, s.tableColumn, tview.NewTableCell(field))
	s.tableColumn += 1
	s.table.SetCell(s.tableRow, s.tableColumn, tview.NewTableCell(value))
	s.tableRow += 1
	s.tableColumn = 0
}

func (s *Selector) AddtoTableList(tableList []TableListItem, field string, value string) {
	tableList = append(tableList, TableListItem{
		Field: field,
		Value: value,
	})
}

func (s *Selector) createContextList() {
	s.list.ShowSecondaryText(false)
	s.list.SetBorder(true).SetTitle("Context")
	for _, config := range s.kubeConfigs {
		var tableList []TableListItem
		for _, ctx := range config.Contexts {

			s.AddtoTableList(tableList, "Cluster", ctx.Cluster)
			s.AddtoTableList(tableList, "File", ctx.LocationOfOrigin)
			s.list.AddItem(ctx.Cluster, "", 0, func() {
				// insert select kubeconfig code

				//s.tableRow = 0
				//s.tableColumn = 0
				//s.AddtoTable("Cluster", ctx.Cluster)
				//s.AddtoTable("File", ctx.LocationOfOrigin)
			})
		}
		s.tableList = append(s.tableList, tableList...)
	}
	s.list.AddItem("Quit", "", 'q', func() { s.app.Stop() })

	s.list.SetChangedFunc(func(index int, mainText string, secondayText string, shortcut rune) {
		s.tableRow = 0
		s.tableColumn = 0
		s.AddtoTable("index", strconv.Itoa(index))
		//s.AddtoTable("Cluster", ctx.Cluster)
		//s.AddtoTable("File", ctx.LocationOfOrigin)
	})

}

func (s *Selector) selectKubeconfig() {
	s.app = tview.NewApplication()
	s.flex = tview.NewFlex()
	s.list = tview.NewList()
	s.pages = tview.NewPages()
	s.table = tview.NewTable()

	s.table.SetBorder(true).SetTitle("Cluster")

	s.createContextList()

	spew.Dump(s.tableList)

	s.flex.AddItem(s.list, 0, 1, true)
	s.flex.AddItem(s.table, 0, 2, false)

	s.pages.AddPage("selectorPage", s.flex, true, true)
	err := s.app.SetRoot(s.pages, true).Run()
	if err != nil {
		logrus.Panicf("Error: %v", err)
	}

}
