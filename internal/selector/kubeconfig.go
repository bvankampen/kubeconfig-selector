package selector

import (
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
	"path/filepath"
	"strings"
)

func loadActiveKubeConfig(dir string, file string) api.Config {
	dir, _ = homedir.Expand(dir)
	logrus.Debugf("Load Active kubeConfig from: %s/%s", dir, file)
	config, err := clientcmd.LoadFromFile(filepath.Join(dir, file))
	if err != nil {
		logrus.Fatalf("Error loading kubeConfig %v", err)
	}
	return *config
}

func loadKubeConfig(dir string, file string) api.Config {
	logrus.Debugf("Loading kubeConfig file: %s/%s", dir, file)
	config, err := clientcmd.LoadFromFile(filepath.Join(dir, file))
	if err != nil {
		logrus.Fatalf("Error loading KubeConfig %v", err)
	}
	return *config
}

func loadKubeConfigsFromDirectory(dir string) []api.Config {
	var apiConfigs []api.Config
	dir, _ = homedir.Expand(dir)
	logrus.Debugf("Searching directory: %s", dir)
	files, err := os.ReadDir(dir)
	if err != nil {
		logrus.Fatalf("Error reading directory: %s (%v)", dir, err)
	}
	for _, file := range files {
		if !file.IsDir() {
			if strings.HasSuffix(file.Name(), ".yaml") {
				logrus.Debugf("Found file: %s", file.Name())
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

//type KubeConfig struct {
//	APIVersion     string      `yaml:"apiVersion"`
//	Clusters       []Clusters  `yaml:"clusters"`
//	Contexts       []Contexts  `yaml:"contexts"`
//	CurrentContext string      `yaml:"current-context"`
//	Kind           string      `yaml:"kind"`
//	Preferences    Preferences `yaml:"preferences"`
//	Users          []Users     `yaml:"users"`
//	ShortFilename  string
//	FullFilename   string
//}
//type Cluster struct {
//	CertificateAuthorityData string `yaml:"certificate-authority-data"`
//	Server                   string `yaml:"server"`
//}
//type Clusters struct {
//	Cluster Cluster `yaml:"cluster"`
//	Name    string  `yaml:"name"`
//}
//type Context struct {
//	Cluster string `yaml:"cluster"`
//	User    string `yaml:"user"`
//}
//type Contexts struct {
//	Context Context `yaml:"context"`
//	Name    string  `yaml:"name"`
//}
//type Preferences struct {
//}
//type User struct {
//	ClientCertificateData string `yaml:"client-certificate-data"`
//	ClientKeyData         string `yaml:"client-key-data"`
//}
//type Users struct {
//	Name string `yaml:"name"`
//	User User   `yaml:"user"`
//}
//
//type kubeConfigs struct {
//	current KubeConfig
//	list    []KubeConfig
//}
//
//func Parse(filename string, shortFilename string) KubeConfig {
//	yamlFile, _ := os.ReadFile(filename)
//	var kubeConfig KubeConfig
//	err := yaml.Unmarshal(yamlFile, &kubeConfig)
//	if err != nil {
//		return KubeConfig{}
//	}
//	kubeConfig.FullFilename = filename
//	kubeConfig.ShortFilename = shortFilename
//	return kubeConfig
//}
//
//func (k *kubeConfigs) ParseKubeConfigs(folder string, configFile string) {
//	if strings.HasPrefix(folder, "~/") {
//		folder, _ = homedir.Expand(folder)
//	}
//
//	files, err := os.ReadDir(folder)
//
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//
//	for _, file := range files {
//		filename := folder + "/" + file.Name()
//		if !file.IsDir() {
//			if file.Name() == configFile {
//				k.current = Parse(filename, file.Name())
//			}
//			if strings.HasSuffix(file.Name(), ".yaml") {
//				k.list = append(k.list, Parse(filename, file.Name()))
//			}
//		}
//	}
//}
//
//func (k *kubeConfigs) GetCurrent() KubeConfig {
//	return k.current
//}
//
//func (k *KubeConfig) GetUser() string {
//	for _, context := range k.Contexts {
//		if context.Name == k.CurrentContext {
//			return context.Context.User
//		}
//	}
//	return ""
//}
//
//func (k *KubeConfig) GetServer() string {
//	for _, context := range k.Contexts {
//		if context.Name == k.CurrentContext {
//			for _, cluster := range k.Clusters {
//				if cluster.Name == context.Context.Cluster {
//					return cluster.Cluster.Server
//				}
//			}
//		}
//	}
//	return ""
//}
