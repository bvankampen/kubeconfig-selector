package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
)

const (
	appConfigFilename = ".app.config"
)

var (
	defaultConf = AppConfig{
		KubeconfigFolder: "~/.kube",
		KubeconfigFile:   "config",
	}
)

type AppConfig struct {
	KubeconfigFolder string `yaml:"kubeconfigFolder"`
	KubeconfigFile   string `yaml:"kubeconfigFile"`
}

func getFilePath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, appConfigFilename)
}

func (c *AppConfig) ConfigLoad() {
	filename := getFilePath()

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		*c = defaultConf
		data, _ := yaml.Marshal(c)
		os.WriteFile(getFilePath(), data, 0600)
	} else {
		yamlConfig, err := os.ReadFile(getFilePath())
		if err != nil {
			fmt.Println(err.Error())
		}
		yaml.Unmarshal(yamlConfig, &c)
	}

}
