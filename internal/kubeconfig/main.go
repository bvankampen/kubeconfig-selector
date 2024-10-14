package kubeconfig

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/bvankampen/kubeconfig-selector/internal/config"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	MARK = "# Written by rs"
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

func checkSuffixes(fileName string) bool {
	suffixes := []string{".yaml", ".config", ".conf"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(fileName, suffix) {
			return true
		}
	}
	return false
}

func loadKubeConfigsFromDirectory(dir string) []api.Config {
	var apiConfigs []api.Config
	dir, _ = homedir.Expand(dir)
	files, err := os.ReadDir(dir)
	if err != nil {
		logrus.Errorf("Error reading directory: %s (%v)", dir, err)
		return nil
	}
	for _, file := range files {
		if !file.IsDir() {
			// if strings.HasSuffix(file.Name(), ".yaml") {
			if checkSuffixes(file.Name()) {
				config, err := loadKubeConfig(dir, file.Name())
				if err == nil {
					apiConfigs = append(apiConfigs, config)
				}
			}
		}
	}
	return apiConfigs
}

func LoadKubeConfigs(appconfig config.AppConfig) ([]api.Config, api.Config) {
	kubeConfigDir, _ := homedir.Expand(appconfig.KubeconfigDir)
	_, err := os.ReadDir(kubeConfigDir)
	if err != nil {
		logrus.Fatalf("Error reading kubeconfig directory: %s (%v)", kubeConfigDir, err)
	}
	apiConfigs := loadKubeConfigsFromDirectory(appconfig.KubeconfigDir)
	for _, dir := range appconfig.ExtraKubeconfigDirs {
		apiConfigs = append(apiConfigs, loadKubeConfigsFromDirectory(dir)...)
	}

	activeConfig := loadActiveKubeConfig(appconfig.KubeconfigDir, appconfig.KubeconfigFile)

	return apiConfigs, activeConfig
}

func markKubeConfig(path string) {
	data, _ := os.ReadFile(path)
	config := MARK + "\n" + string(data)
	os.WriteFile(path, []byte(config), 0600)
}

func checkMark(path string) bool {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return true
	}
	data, _ := os.ReadFile(path)
	if strings.HasPrefix(string(data), MARK) {
		return true
	}
	return false
}

func SaveKubeConfig(config *api.Config, context string, dir string, file string, doCheckMark bool, createLink bool, isMove bool) error {
	dir, _ = homedir.Expand(dir)
	path := filepath.Join(dir, file)
	config.CurrentContext = context
	if createLink {
		kubeConfigLocation := config.Contexts[context].LocationOfOrigin
		err := clientcmd.WriteToFile(*config, kubeConfigLocation)
		if err != nil {
			return errors.New("Unable to write " + path + " Error: " + err.Error())
		}
		_, err = os.Stat(path)
		if err == nil {
			fileInfo, _ := os.Lstat(path)
			if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
				os.Remove(path)
			} else {
				return errors.New("File: " + path + " is not a symlink, please remove/rename this file first.")
			}
		}
		if isMove {
			// is kubeconfig is moved, then use new location instead of original location
			kubeConfigLocation = filepath.Join(dir, filepath.Base(kubeConfigLocation))
		}
		err = os.Symlink(kubeConfigLocation, path)
		if err != nil {
			return errors.New("Unable to create Symlink " + kubeConfigLocation + "->" + path + "Error: " + err.Error())
		}
	} else {
		if !checkMark(path) && doCheckMark {
			return errors.New("Kubeconfig (" + path + ") is not managed by rs. Remove/rename this file first.")
		}
		err := clientcmd.WriteToFile(*config, path)
		if err != nil {
			return errors.New("Unable to write " + path + " Error: " + err.Error())
		}
		markKubeConfig(path)
	}
	return nil
}

func MoveKubeConfig(config *api.Config, context string, kubeConfigDir string) error {
	originalKubeConfigLocation := config.Contexts[context].LocationOfOrigin
	filename := filepath.Base(originalKubeConfigLocation)
	dir, _ := homedir.Expand(kubeConfigDir)
	err := os.Rename(originalKubeConfigLocation, filepath.Join(dir, filename))
	if err != nil {
		return err
	}
	os.Chmod(filepath.Join(dir, filename), 0600)
	return nil
}

func DeleteKubeConfig(config *api.Config, context string, dir string, file string, createLink bool, activeContext bool) error {
	if createLink && activeContext {
		dir, _ = homedir.Expand(dir)
		path := filepath.Join(dir, file)
		os.Remove(path)
	}
	originalKubeConfigLocation := config.Contexts[context].LocationOfOrigin
	os.Remove(originalKubeConfigLocation)
	return nil
}
