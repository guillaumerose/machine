package host

import (
	"github.com/code-ready/machine/libmachine/auth"
	"github.com/code-ready/machine/libmachine/engine"
)

// MigrateHostMetadataV0ToHostMetadataV1 fills nested host metadata and modifies if needed
// this is used for configuration updates
func MigrateHostMetadataV0ToHostMetadataV1(m *MetadataV0) *Metadata {
	hostMetadata := &Metadata{}
	hostMetadata.DriverName = m.DriverName
	hostMetadata.HostOptions.EngineOptions = &engine.Options{}
	hostMetadata.HostOptions.AuthOptions = &auth.Options{
		StorePath:            m.StorePath,
		CaCertPath:           m.CaCertPath,
		CaCertRemotePath:     "",
		ServerCertPath:       m.ServerCertPath,
		ServerKeyPath:        m.ServerKeyPath,
		ClientKeyPath:        "",
		ServerCertRemotePath: "",
		ServerKeyRemotePath:  "",
		CaPrivateKeyPath:     m.PrivateKeyPath,
		ClientCertPath:       m.ClientCertPath,
	}

	hostMetadata.ConfigVersion = m.ConfigVersion

	return hostMetadata
}
