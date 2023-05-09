package app

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"strings"
)

type listItem struct {
	Context       string
	ContextLabel  string
	ShortFilename string
	FullFilename  string
	Server        string
	User          string
	IsCurrent     bool
	IsExit        bool
}

type Application struct {
	KubeConfigFile string
	NewConfig      listItem
	CurrentConfig  listItem
}

func (a *Application) Run() {
	var k kubeConfigs
	var appConf AppConfig
	var list []listItem
	home, _ := homedir.Dir()
	a.KubeConfigFile = filepath.Join(home, ".kube/config")

	appConf.ConfigLoad()

	k.ParseKubeConfigs(appConf.KubeconfigFolder, appConf.KubeconfigFile)

	current := k.GetCurrent()
	a.CurrentConfig = listItem{
		Context:       current.CurrentContext,
		ShortFilename: current.ShortFilename,
		FullFilename:  current.FullFilename,
		User:          current.GetUser(),
		Server:        current.GetServer(),
	}

	for _, config := range k.list {
		if config.CurrentContext != "" {
			isCurrent := false
			if config.CurrentContext == a.CurrentConfig.Context {
				isCurrent = true
			}

			var contextLabel string
			contextLabel = config.CurrentContext
			if contextLabel == "default" {
				contextLabel = fmt.Sprintf("%s [%s]", strings.Replace(config.ShortFilename, ".yaml", "", -1), config.CurrentContext)
			}

			list = append(list, listItem{
				Context:       config.CurrentContext,
				ContextLabel:  contextLabel,
				ShortFilename: config.ShortFilename,
				FullFilename:  config.FullFilename,
				User:          config.GetUser(),
				Server:        config.GetServer(),
				IsCurrent:     isCurrent,
				IsExit:        false,
			})
		}
	}

	list = append(list, listItem{
		ContextLabel: "Exit",
		IsExit:       true,
	})

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "> {{if .IsCurrent}}{{ .ContextLabel | blue }} (current){{else}}{{ .ContextLabel | cyan }}{{end}}",
		Inactive: "  {{if .IsCurrent}}{{ .ContextLabel | blue }} (current){{else}}{{ .ContextLabel }}{{end}}",
		Selected: "{{if not .IsExit}}Change context to: {{ .Context }}{{else}}No changes{{end}}",
		Details: `{{if not .IsExit}}
{{ "Context       :" | faint }} {{ .Context }}
{{ "ShortFilename :" | faint }} {{ .ShortFilename }}
{{ "Server        :" | faint }} {{ .Server }}
{{ "User          :" | faint }} {{ .User }}{{end}}`,
	}

	searcher := func(input string, index int) bool {
		item := list[index]
		name := strings.Replace(strings.ToLower(item.Context), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select Cluster",
		Items:     list,
		Templates: templates,
		Searcher:  searcher,
		Size:      20,
		HideHelp:  true,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	a.NewConfig = list[i]

	if !a.NewConfig.IsExit && !a.NewConfig.IsCurrent {
		a.switchToKubeConfig()
	}
}

func (a *Application) switchToKubeConfig() {
	_, errStat := os.Stat(a.KubeConfigFile)
	if errStat == nil {
		file, errEvalSymLinks := filepath.EvalSymlinks(a.KubeConfigFile)
		if errEvalSymLinks == nil {
			if file == a.KubeConfigFile {
				fmt.Printf("Error: generating symlink %s -> %s\n", a.NewConfig.FullFilename, a.KubeConfigFile)
				return
			} else {
				errRemove := os.Remove(a.KubeConfigFile)
				if errRemove != nil {
					fmt.Printf("Error: removing %s\n", a.KubeConfigFile)
				}
			}
		}
	}

	errSymlink := os.Symlink(a.NewConfig.FullFilename, a.KubeConfigFile)
	if errSymlink != nil {
		fmt.Printf("Error: generating symlink %s -> %s\n", a.NewConfig.FullFilename, a.KubeConfigFile)
	}

}
