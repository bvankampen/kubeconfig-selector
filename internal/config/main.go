package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
)

const (
	appConfigFilename = "~/.config/ks.conf"
)

var defaultConf = AppConfig{
	KubeconfigDir:       "~/.kube",
	KubeconfigFile:      "config",
	ExtraKubeconfigDirs: []string{},
	ShowKubeConfig:      true,
	CreateLink:          true,
	RancherFix:          true,
}

type AppConfig struct {
	KubeconfigDir       string   `yaml:"kubeconfigDir"`
	KubeconfigFile      string   `yaml:"kubeconfigFile"`
	ExtraKubeconfigDirs []string `yaml:"extraKubeconfigDirs"`
	ShowKubeConfig      bool     `yaml:"showKubeconfig"`
	CreateLink          bool     `yaml:"createLink"`
	RancherFix          bool     `yaml:"rancherFix"`
}

func LoadAppConfig() (*AppConfig, error) {
	filename, _ := homedir.Expand(appConfigFilename)

	appconfig := AppConfig{}

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		appconfig = defaultConf
		data, err := yaml.Marshal(appconfig)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal default config: %w", err)
		}
		if err := os.WriteFile(filename, data, 0o600); err != nil {
			return nil, fmt.Errorf("unable to write configfile: %w", err)
		}
	} else {
		yamlConfig, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("unable to read configfile: %w", err)
		}
		err = yaml.Unmarshal(yamlConfig, &appconfig)
		if err != nil {
			return nil, fmt.Errorf("unable to load configfile: %w", err)
		}
	}
	return &appconfig, nil
}

func WriteAppConfig(appconfig *AppConfig) error {
	filename, _ := homedir.Expand(appConfigFilename)
	data, err := yaml.Marshal(appconfig)
	if err != nil {
		return fmt.Errorf("unable to marshal config: %w", err)
	}
	if err := os.WriteFile(filename, data, 0o600); err != nil {
		return fmt.Errorf("unable to write configfile: %w", err)
	}
	return nil
}
