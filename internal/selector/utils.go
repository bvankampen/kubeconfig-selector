package selector

import (
	sha256 "crypto/sha256"
	b64 "encoding/base64"
	hex "encoding/hex"
	"k8s.io/client-go/tools/clientcmd"
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

func getHash(config api.Config) string {
	configBytes, _ := clientcmd.Write(config)
	hash := sha256.New()
	hash.Write(configBytes)
	return hex.EncodeToString(hash.Sum(nil))
}
