package ui

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func redactConfigToString(config api.Config) string {
	redactedString := "REDACTED"

	for _, cluster := range config.Clusters {
		cluster.CertificateAuthorityData = nil
	}
	for _, authInfo := range config.AuthInfos {
		authInfo.ClientCertificateData = nil
		authInfo.ClientKeyData = nil
		if authInfo.Token != "" {
			authInfo.Token = redactedString
		}
		if authInfo.Password != "" {
			authInfo.Password = redactedString
		}
	}
	configBytes, _ := clientcmd.Write(config)
	return string(configBytes)
}
