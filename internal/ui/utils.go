package ui

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func redactConfigToString(config api.Config) string {
	redactedBytes := []byte("REDACTED")
	redactedString := "REDACTED"

	for _, cluster := range config.Clusters {
		if len(cluster.CertificateAuthorityData) > 0 {
			cluster.CertificateAuthorityData = redactedBytes
		}
	}
	for _, authInfo := range config.AuthInfos {
		if len(authInfo.ClientCertificateData) > 0 {
			authInfo.ClientCertificateData = redactedBytes
		}
		if len(authInfo.ClientKeyData) > 0 {
			authInfo.ClientKeyData = redactedBytes
		}
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
