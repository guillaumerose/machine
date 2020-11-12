package auth

type Options struct {
	CertDir              string
	CaCertPath           string
	CaPrivateKeyPath     string
	CaCertRemotePath     string
	ServerCertPath       string
	ServerKeyPath        string
	ClientKeyPath        string
	ServerCertRemotePath string
	ServerKeyRemotePath  string
	ClientCertPath       string
	ServerCertSANs       []string
}
