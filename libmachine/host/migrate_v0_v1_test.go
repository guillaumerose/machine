package host

import (
	"reflect"
	"testing"

	"github.com/code-ready/machine/libmachine/auth"
	"github.com/code-ready/machine/libmachine/engine"
)

func TestMigrateHostMetadataV0ToV1(t *testing.T) {
	metadata := &MetadataV0{
		HostOptions: Options{
			EngineOptions: nil,
			AuthOptions:   nil,
		},
		StorePath:      "/tmp/store",
		CaCertPath:     "/tmp/store/certs/ca.pem",
		ServerCertPath: "/tmp/store/certs/server.pem",
	}
	expectedAuthOptions := &auth.Options{
		CaCertPath:     "/tmp/store/certs/ca.pem",
		ServerCertPath: "/tmp/store/certs/server.pem",
	}

	expectedMetadata := &Metadata{
		HostOptions: Options{
			EngineOptions: &engine.Options{},
			AuthOptions:   expectedAuthOptions,
		},
	}

	m := MigrateHostMetadataV0ToHostMetadataV1(metadata)

	if !reflect.DeepEqual(m, expectedMetadata) {
		t.Logf("\n%+v\n%+v", m, expectedMetadata)
		t.Fatal("Expected these structs to be equal, they were different")
	}
}
