package selector

import (
	"context"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"strconv"
)

type TableListItem struct {
	Field string
	Value string
}

type ConfigList struct {
	Rows           []TableListItem
	Config         api.Config
	RedactedConfig api.Config
	Context        *api.Context
}

type Selector struct {
	ctx          context.Context
	appConfig    AppConfig
	kubeConfigs  []api.Config
	activeConfig api.Config
	app          *tview.Application
	list         *tview.List
	table        *tview.Table
	configView   *tview.TextView
	debugView    *tview.TextView
	tableRow     int
	tableColumn  int
	configList   []ConfigList
	debug        bool
}

func New(ctx context.Context, debug bool) (*Selector, error) {

	appconfig := loadAppConfig()
	kubeconfigs, activeconfig := loadKubeConfigs(appconfig)

	return &Selector{
		ctx:          ctx,
		appConfig:    *appconfig,
		kubeConfigs:  kubeconfigs,
		activeConfig: activeconfig,
		debug:        debug,
	}, nil

}

func (s *Selector) addtoTable(field string, value string) {
	s.table.SetCell(s.tableRow, s.tableColumn, tview.NewTableCell(field))
	s.tableColumn += 1
	s.table.SetCell(s.tableRow, s.tableColumn, tview.NewTableCell(value))
	s.tableRow += 1
	s.tableColumn = 0
}

func (s *Selector) addtoTableList(tableList *ConfigList, field string, value string) {
	tableList.Rows = append(tableList.Rows, TableListItem{
		Field: field,
		Value: value,
	})
}

func (s *Selector) printDebug(str string, addToText bool) {
	if s.debug {
		if addToText {
			currentStr := s.debugView.GetText(false)
			if currentStr != "" {
				str = currentStr + "\n" + str
			}
		}
		s.debugView.SetText(str)
		s.debugView.ScrollToEnd()
	}
}

func (s *Selector) updateScreen(index int) {
	s.tableRow = 0
	s.tableColumn = 0
	if index < len(s.configList) {
		for _, item := range s.configList[index].Rows {
			s.addtoTable(item.Field, item.Value)
		}

		configBytes, _ := clientcmd.Write(s.configList[index].Config)
		s.configView.SetText(string(configBytes))

	} else {
		s.table.Clear()
		s.configView.Clear()
	}
}

func (s *Selector) createContextList() {
	activeConfigHash := getHash(s.activeConfig)
	index := 0
	currentIndex := 0
	s.list.ShowSecondaryText(false)
	s.list.SetBorder(true).SetTitle("Context")
	s.list.SetHighlightFullLine(true)
	for _, config := range s.kubeConfigs {
		var tableList ConfigList
		for _, configContext := range config.Contexts {

			s.addtoTableList(&tableList, "Cluster", configContext.Cluster)
			s.addtoTableList(&tableList, "File", configContext.LocationOfOrigin)
			s.list.AddItem(configContext.Cluster, "", 0, func() {
				// insert select kubeconfig code
			})
			tableList.Context = configContext
			tableList.Config = *config.DeepCopy()
			tableList.RedactedConfig = redactConfig(*config.DeepCopy())
		}
		s.configList = append(s.configList, tableList)
		configHash := getHash(config)

		s.printDebug("c:"+activeConfigHash, true)
		s.printDebug(strconv.Itoa(index)+":"+configHash, true)

		if configHash == activeConfigHash {
			currentIndex = index
		}
		s.printDebug("currentIndex: "+strconv.Itoa(currentIndex), true)
		index++
	}
	s.list.AddItem("Quit", "", 0, func() { s.app.Stop() })

	s.list.SetChangedFunc(func(index int, mainText string, secondayText string, shortcut rune) {
		s.updateScreen(currentIndex)
	})

	//s.list.SetCurrentItem(index)

	s.updateScreen(0)

}

func (s *Selector) Run() error {
	s.app = tview.NewApplication()

	s.list = tview.NewList()
	s.configView = tview.NewTextView()
	s.debugView = tview.NewTextView()
	s.table = tview.NewTable()

	s.table.SetBorder(true).SetTitle("Cluster")
	s.configView.SetBorder(true).SetTitle("Kubeconfig")
	s.debugView.SetBorder(true).SetTitle("Debug")

	s.createContextList()

	flexSub := tview.NewFlex().SetDirection(tview.FlexRow)
	flexSub.AddItem(s.table, 0, 1, false)
	if s.appConfig.ShowKubeConfig {
		flexSub.AddItem(s.configView, 0, 2, false)
	}
	if s.debug {
		flexSub.AddItem(s.debugView, 0, 3, false)
	}
	flexMain := tview.NewFlex()
	flexMain.AddItem(s.list, 0, 1, true)
	flexMain.AddItem(flexSub, 0, 2, false)

	pages := tview.NewPages().AddPage("selectorPage", flexMain, true, true)
	err := s.app.SetRoot(pages, true).Run()
	if err != nil {
		logrus.Panicf("Error: %v", err)
	}
	return nil
}
