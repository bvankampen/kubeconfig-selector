package selector

import (
	b64 "encoding/base64"

	"k8s.io/client-go/tools/clientcmd/api"
)

func redactConfig(config api.Config) api.Config {
	redacted, _ := b64.StdEncoding.DecodeString("REDACTED")
	for _, cluster := range config.Clusters {
		if len(cluster.CertificateAuthorityData) > 0 {
			cluster.CertificateAuthorityData = redacted
		}
	}
	for _, authInfo := range config.AuthInfos {
		if len(authInfo.ClientCertificateData) > 0 {
			authInfo.ClientCertificateData = redacted
		}
		if len(authInfo.ClientKeyData) > 0 {
			authInfo.ClientKeyData = redacted
		}
		if len(authInfo.Token) > 0 {
			authInfo.Token = "REDACTED"
		}
	}
	return config
}
