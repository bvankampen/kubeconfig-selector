package selector

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bvankampen/kubeconfig-selector/internal/config"
	"github.com/bvankampen/kubeconfig-selector/internal/kubeconfig"
	"github.com/bvankampen/kubeconfig-selector/internal/rancher"
	"github.com/mitchellh/go-homedir"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type ListEntry struct {
	Name         string
	SourceFile   string
	PrefixSymbol rune
}

func BuildSortedEntries(kubeConfigs []api.Config, kubeconfigDir string) []ListEntry {
	seen := make(map[string]bool)
	var entries []ListEntry
	kubeDir, _ := homedir.Expand(kubeconfigDir)

	for _, cfg := range kubeConfigs {
		for name, cfgContext := range cfg.Contexts {
			if seen[name] {
				continue
			}
			seen[name] = true

			var prefixSymbol rune
			if !strings.HasPrefix(cfgContext.LocationOfOrigin, kubeDir) {
				prefixSymbol = '*'
			} else if cluster, ok := cfg.Clusters[cfgContext.Cluster]; ok {
				if strings.HasSuffix(cluster.Server, "local") {
					prefixSymbol = 'r'
				}
			}

			entries = append(entries, ListEntry{Name: name, SourceFile: cfgContext.LocationOfOrigin, PrefixSymbol: prefixSymbol})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	return entries
}

func GetConfigByIndex(kubeConfigs []api.Config, listEntries []ListEntry, index int) (string, api.Config, *api.Context) {
	if index >= len(listEntries) {
		return "", api.Config{}, &api.Context{}
	}
	entry := listEntries[index]
	for _, config := range kubeConfigs {
		if ctx, ok := config.Contexts[entry.Name]; ok {
			if ctx.LocationOfOrigin == entry.SourceFile {
				return entry.Name, config, ctx
			}
		}
	}
	return "", api.Config{}, &api.Context{}
}

func DeleteConfigByIndex(kubeConfigs []api.Config, listEntries []ListEntry, index int) []api.Config {
	if index >= len(listEntries) {
		return kubeConfigs
	}
	entry := listEntries[index]
	for i, config := range kubeConfigs {
		if ctx, ok := config.Contexts[entry.Name]; ok {
			if ctx.LocationOfOrigin == entry.SourceFile {
				return append(kubeConfigs[:i], kubeConfigs[i+1:]...)
			}
		}
	}
	return kubeConfigs
}

func FindActiveIndex(kubeConfigs []api.Config, listEntries []ListEntry, activeConfig api.Config) int {
	if activeConfig.CurrentContext == "" {
		return 0
	}
	for i, entry := range listEntries {
		for _, cfg := range kubeConfigs {
			if cfgContext, ok := cfg.Contexts[entry.Name]; ok {
				activeConfigContext := activeConfig.Contexts[activeConfig.CurrentContext]
				activeConfigCluster := activeConfigContext.Cluster
				activeConfigServer := activeConfig.Clusters[activeConfigContext.Cluster].Server

				if cfgContext.Cluster == activeConfigCluster &&
					cfg.Clusters[cfgContext.Cluster].Server == activeConfigServer &&
					entry.Name == activeConfig.CurrentContext {
					return i
				}
				break
			}
		}
	}
	return 0
}

func SelectConfig(kubeConfigs []api.Config, listEntries []ListEntry, index int, appConfig config.AppConfig) error {
	name, config, _ := GetConfigByIndex(kubeConfigs, listEntries, index)
	return kubeconfig.SaveKubeConfig(
		config.DeepCopy(),
		name,
		appConfig.KubeconfigDir,
		appConfig.KubeconfigFile,
		true,
		appConfig.CreateLink,
		false)
}

func MoveConfig(kubeConfigs []api.Config, listEntries []ListEntry, index int, appConfig config.AppConfig) error {
	name, config, _ := GetConfigByIndex(kubeConfigs, listEntries, index)
	return kubeconfig.SaveKubeConfig(
		config.DeepCopy(),
		name,
		appConfig.KubeconfigDir,
		appConfig.KubeconfigFile,
		true,
		appConfig.CreateLink,
		true)
}

func DeleteConfig(kubeConfigs []api.Config, listEntries []ListEntry, index int, appConfig config.AppConfig, activeConfig api.Config) error {
	name, config, _ := GetConfigByIndex(kubeConfigs, listEntries, index)
	activeContext := activeConfig.CurrentContext == name
	return kubeconfig.DeleteKubeConfig(
		config.DeepCopy(),
		name,
		appConfig.KubeconfigDir,
		appConfig.KubeconfigFile,
		appConfig.CreateLink,
		activeContext)
}

func RenameContext(kubeConfigs []api.Config, config api.Config, contextName string, newContextName string, kubeconfigDir string, kubeconfigFile string, createLink bool) error {
	if contextName == newContextName {
		return nil
	}
	for _, cfg := range kubeConfigs {
		if _, exists := cfg.Contexts[newContextName]; exists {
			return fmt.Errorf("Context %q already exists.", newContextName)
		}
	}
	for name, context := range config.Contexts {
		if name == contextName {
			kubeConfigPath := filepath.Dir(context.LocationOfOrigin)
			kubeConfigFilename := filepath.Base(context.LocationOfOrigin)
			config.Contexts[newContextName] = context.DeepCopy()
			config.CurrentContext = newContextName
			delete(config.Contexts, contextName)
			err := kubeconfig.SaveKubeConfigFile(
				config.DeepCopy(),
				newContextName,
				kubeConfigPath,
				kubeConfigFilename,
			)
			if err != nil {
				return err
			}

			ext := filepath.Ext(kubeConfigFilename)
			newFilename := newContextName + ext
			oldPath := filepath.Join(kubeConfigPath, kubeConfigFilename)
			newPath := filepath.Join(kubeConfigPath, newFilename)
			err = os.Rename(oldPath, newPath)
			if err != nil {
				return err
			}
			context.LocationOfOrigin = newPath

			if createLink {
				activePath, _ := homedir.Expand(filepath.Join(kubeconfigDir, kubeconfigFile))
				fi, err := os.Lstat(activePath)
				if err == nil && fi.Mode()&os.ModeSymlink != 0 {
					linkTarget, err := os.Readlink(activePath)
					if err == nil && linkTarget == oldPath {
						os.Remove(activePath)
						os.Symlink(newPath, activePath)
					}
				}
			}

			for i, cfg := range kubeConfigs {
				for n := range cfg.Contexts {
					if n == contextName {
						kubeConfigs[i] = config
						return nil
					}
				}
			}
			return nil
		}
	}
	return nil
}

func DownloadDownstreamKubeConfig(server, token string, cluster rancher.DownstreamCluster, kubeconfigDir string, insecure bool) error {
	data, err := rancher.FetchClusterKubeConfig(server, token, cluster.ID, insecure)
	if err != nil {
		return err
	}

	cfg, err := clientcmd.Load(data)
	if err != nil {
		return fmt.Errorf("Failed to parse kubeconfig for %s: %v", cluster.Name, err)
	}

	kubeDir, _ := homedir.Expand(kubeconfigDir)
	filePath := filepath.Join(kubeDir, cluster.Name+".yaml")

	err = clientcmd.WriteToFile(*cfg, filePath)
	if err != nil {
		return fmt.Errorf("Failed to save kubeconfig: %v", err)
	}

	os.Chmod(filePath, 0o600)
	return nil
}
