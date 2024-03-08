package selector

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	appConfigFilename = "~/.config/ks.conf"
)

var (
	defaultConf = AppConfig{
		KubeconfigDir:       "~/.kube",
		KubeconfigFile:      "config",
		ExtraKubeconfigDirs: []string{"~/Downloads"},
		ShowKubeConfig:      true,
	}
)

type AppConfig struct {
	KubeconfigDir       string   `yaml:"kubeconfigDir"`
	KubeconfigFile      string   `yaml:"kubeconfigFile"`
	ExtraKubeconfigDirs []string `yaml:"extraKubeconfigDirs"`
	ShowKubeConfig      bool     `yaml:"showKubeconfig"`
}

func loadAppConfig() *AppConfig {
	filename, _ := homedir.Expand(appConfigFilename)

	appconfig := AppConfig{}

	logrus.Debugf("Loading configfile: %s", filename)

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		appconfig = defaultConf
		data, _ := yaml.Marshal(appconfig)
		logrus.Debugf("Configfile doesn't exist creating default configfile")
		err := os.WriteFile(filename, data, 0600)
		if err != nil {
			logrus.Errorf("Unable to write configfile: %v", err)
		}
	} else {
		yamlConfig, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println(err.Error())
		}
		err = yaml.Unmarshal(yamlConfig, &appconfig)
		if err != nil {
			logrus.Errorf("Unable to load configfile: %v", err)
		}
	}

	return &appconfig

}
