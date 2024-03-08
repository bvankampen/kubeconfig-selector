package selector

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func loadActiveKubeConfig(dir string, file string) api.Config {
	dir, _ = homedir.Expand(dir)
	config, err := clientcmd.LoadFromFile(filepath.Join(dir, file))
	if err != nil {
		logrus.Errorf("Error loading kubeConfig %v", err)
		return api.Config{}
	}
	return *config
}

func loadKubeConfig(dir string, file string) api.Config {
	config, err := clientcmd.LoadFromFile(filepath.Join(dir, file))
	if err != nil {
		logrus.Fatalf("Error loading KubeConfig %v", err)
	}
	return *config
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
				apiConfigs = append(apiConfigs, loadKubeConfig(dir, file.Name()))
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

func saveKubeConfig(config *api.Config, context string, dir string, file string) {
	dir, _ = homedir.Expand(dir)
	path := filepath.Join(dir, file)
	// path := "/tmp/config"
	config.CurrentContext = context
	err := clientcmd.WriteToFile(*config, path)
	if err != nil {
		logrus.Fatalf("Unable to write %s Error: %v", path, err)
	}

}
