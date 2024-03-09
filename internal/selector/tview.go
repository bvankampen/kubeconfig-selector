package selector

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
	// "github.com/davecgh/go-spew/spew"
)

func (s *Selector) addtoTable(field string, value string) {
	s.table.SetCell(s.tableRow, s.tableColumn, tview.NewTableCell(field).SetTextColor(tcell.ColorOrange))
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
	for _, item := range s.configList[index].Rows {
		s.addtoTable(item.Field, item.Value)
	}

	configBytes, _ := clientcmd.Write(s.configList[index].RedactedConfig)
	s.configView.SetText(string(configBytes))
}

func (s *Selector) createContextList() {
	activeConfigHash := getHash(s.activeConfig)
	index := 0
	currentIndex := 0
	s.list.ShowSecondaryText(false)
	s.list.SetBorder(true).SetTitle("Context").SetBorderColor(tcell.ColorBlue)
	s.list.SetHighlightFullLine(true)
	for _, config := range s.kubeConfigs {
		for name, configContext := range config.Contexts {
			var tableList ConfigList
			s.addtoTableList(&tableList, "Context", name)
			s.addtoTableList(&tableList, "Cluster", configContext.Cluster)
			s.addtoTableList(&tableList, "User", configContext.AuthInfo)
			s.addtoTableList(&tableList, "Server", config.Clusters[configContext.Cluster].Server)
			s.addtoTableList(&tableList, "File", configContext.LocationOfOrigin)

			kubeDir, _ := homedir.Expand(s.appConfig.KubeconfigDir)

			var star rune
			star = 0
			if !strings.HasPrefix(configContext.LocationOfOrigin, kubeDir) {
				star = '*'
			}

			s.list.AddItem(name, "", star, nil)
			tableList.Context = configContext
			tableList.Config = *config.DeepCopy()
			tableList.RedactedConfig = redactConfig(*config.DeepCopy())
			s.configList = append(s.configList, tableList)
			configHash := getHash(config)

			if configHash == activeConfigHash {
				currentIndex = index
			}
			index++
		}
	}

	s.list.SetChangedFunc(func(index int, mainText string, secondayText string, shortcut rune) {
		s.updateScreen(index)
	})

	s.list.SetSelectedFunc(func(index int, mainText string, secondayText string, shortcut rune) {
		saveKubeConfig(s.configList[index].Config.DeepCopy(), mainText, s.appConfig.KubeconfigDir, s.appConfig.KubeconfigFile)
		s.app.Stop()
	})

	s.updateScreen(0)
	s.list.SetCurrentItem(currentIndex)

}

func (s *Selector) createHelpView() {
	s.helpView.SetBorder(false)
	// s.helpView.SetRegions(true)
	s.helpView.SetDynamicColors(true)
	helpText := "[yellow]q:[white] Quit " +
		"[yellow]<enter>:[white] Use Kubeconfig " +
		"[yellow]m:[white] Move Kubeconfig to " + s.appConfig.KubeconfigDir + " and use it " +
		"[yellow]k:[white] Toggle Kubeconfig " +
		"[yellow](*):[white] File not in " + s.appConfig.KubeconfigDir
	s.helpView.SetText(helpText)
}

func (s *Selector) setupPages() *tview.Pages {
	s.list = tview.NewList()
	s.configView = tview.NewTextView()
	s.debugView = tview.NewTextView()
	s.table = tview.NewTable()
	s.helpView = tview.NewTextView()

	s.table.SetBorder(true).SetTitle("Cluster")
	s.configView.SetBorder(true).SetTitle("Kubeconfig")
	s.debugView.SetBorder(true).SetTitle("Debug")

	s.createContextList()
	s.createHelpView()

	title := tview.NewTextView()
	title.SetBackgroundColor(tcell.ColorDarkCyan)
	title.SetTextColor(tcell.ColorBlack)
	title.SetText(fmt.Sprintf(" Kubeconfig Selector %s https://github.com/bvankampen/kubeconfig-selector", s.ctx.App.Version))

	flexViews := tview.NewFlex().SetDirection(tview.FlexRow)
	tableSize := 0
	if s.appConfig.ShowKubeConfig {
		tableSize = 7
	}
	flexViews.AddItem(s.table, tableSize, 1, false)
	if s.appConfig.ShowKubeConfig {
		flexViews.AddItem(s.configView, 0, 2, false)
	}
	if s.debug {
		flexViews.AddItem(s.debugView, 0, 3, false)
	}
	flexMain := tview.NewFlex()
	flexMain.AddItem(s.list, 0, 1, true)
	flexMain.AddItem(flexViews, 0, 2, false)

	flexApp := tview.NewFlex().SetDirection(tview.FlexRow)
	flexApp.AddItem(title, 1, 1, false)
	flexApp.AddItem(flexMain, 0, 1, true)
	flexApp.AddItem(s.helpView, 1, 1, false)

	pages := tview.NewPages().AddPage("selectorPage", flexApp, true, true)

	return pages

}

func (s *Selector) moveKubeconfig() {
	index := s.list.GetCurrentItem()
	config := s.configList[index].Config
	context, _ := s.list.GetItemText(index)
	saveKubeConfig(config.DeepCopy(), context, s.appConfig.KubeconfigDir, s.appConfig.KubeconfigFile)
	orgKubeConfig := config.Contexts[context].LocationOfOrigin
	filename := filepath.Base(orgKubeConfig)
	dir, _ := homedir.Expand(s.appConfig.KubeconfigDir)
	err := os.Rename(orgKubeConfig, filepath.Join(dir, filename))
	if err != nil {
		logrus.Errorf("Unable to move file %v", err)
	}
	os.Chmod(filepath.Join(dir, filename), 0600)
}
