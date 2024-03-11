package selector

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	mark = "# Written by rs"
)

func loadActiveKubeConfig(dir string, file string) api.Config {
	dir, _ = homedir.Expand(dir)
	config, err := clientcmd.LoadFromFile(filepath.Join(dir, file))
	if err != nil {
		logrus.Debugf("Error loading kubeConfig %s/%s \nError:%v", dir, file, err)
		return api.Config{}
	}
	return *config
}

func loadKubeConfig(dir string, file string) (api.Config, error) {
	config, err := clientcmd.LoadFromFile(filepath.Join(dir, file))
	if err != nil {
		logrus.Debugf("Error loading kubeConfig %s/%s \nError:%v", dir, file, err)
		return api.Config{}, err
	}
	return *config, nil
}

func loadKubeConfigsFromDirectory(dir string) []api.Config {
	var apiConfigs []api.Config
	dir, _ = homedir.Expand(dir)
	files, err := os.ReadDir(dir)
	if err != nil {
		logrus.Fatalf("Error reading directory: %s (%v)", dir, err)
	}
	for _, file := range files {
		if !file.IsDir() {
			if strings.HasSuffix(file.Name(), ".yaml") {
				config, err := loadKubeConfig(dir, file.Name())
				if err == nil {
					apiConfigs = append(apiConfigs, config)
				}
			}
		}
	}
	return apiConfigs
}

func loadKubeConfigs(appconfig *AppConfig) ([]api.Config, api.Config) {
	apiConfigs := loadKubeConfigsFromDirectory(appconfig.KubeconfigDir)
	for _, dir := range appconfig.ExtraKubeconfigDirs {
		apiConfigs = append(apiConfigs, loadKubeConfigsFromDirectory(dir)...)
	}

	activeConfig := loadActiveKubeConfig(appconfig.KubeconfigDir, appconfig.KubeconfigFile)

	return apiConfigs, activeConfig
}

func markKubeConfig(path string) {
	data, _ := os.ReadFile(path)
	config := mark + "\n" + string(data)
	os.WriteFile(path, []byte(config), 0600)
}

func checkMark(path string) bool {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return true
	}
	data, _ := os.ReadFile(path)
	if strings.HasPrefix(string(data), mark) {
		return true
	}
	return false
}

func saveKubeConfig(config *api.Config, context string, dir string, file string) error {
	dir, _ = homedir.Expand(dir)
	path := filepath.Join(dir, file)
	config.CurrentContext = context
	if !checkMark(path) {
		return errors.New("Kubeconfig (" + path + ") is not managed by rs. Remove/rename this file first.")
	}
	err := clientcmd.WriteToFile(*config, path)
	markKubeConfig(path)
	if err != nil {
		return errors.New("Unable to write " + path + " Error: " + err.Error())
	}
	return nil
}
