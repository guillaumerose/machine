package host

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/code-ready/machine/drivers/none"
	"github.com/code-ready/machine/libmachine/version"
)

var (
	errUnexpectedConfigVersion = errors.New("unexpected config version")
)

type RawDataDriver struct {
	*none.Driver
	Data []byte // passed directly back when invoking json.Marshal on this type
}

func (r *RawDataDriver) MarshalJSON() ([]byte, error) {
	return r.Data, nil
}

func (r *RawDataDriver) UnmarshalJSON(data []byte) error {
	r.Data = data
	return nil
}

func (r *RawDataDriver) UpdateConfigRaw(rawData []byte) error {
	return r.UnmarshalJSON(rawData)
}

func getMigratedHostMetadata(data []byte) (*Metadata, error) {
	// HostMetadata is for a "first pass" so we can then load the driver
	var (
		hostMetadata *MetadataV0
	)

	if err := json.Unmarshal(data, &hostMetadata); err != nil {
		return &Metadata{}, err
	}

	migratedHostMetadata := MigrateHostMetadataV0ToHostMetadataV1(hostMetadata)

	return migratedHostMetadata, nil
}

func MigrateHost(name string, data []byte) (*Host, error) {
	migratedHostMetadata, err := getMigratedHostMetadata(data)
	if err != nil {
		return nil, err
	}

	if migratedHostMetadata.ConfigVersion != version.ConfigVersion {
		return nil, errUnexpectedConfigVersion
	}

	globalStorePath := filepath.Dir(filepath.Dir(migratedHostMetadata.HostOptions.AuthOptions.StorePath))
	driver := &RawDataDriver{none.NewDriver(name, globalStorePath), nil}
	h := Host{
		Name:   name,
		Driver: driver,
	}
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, fmt.Errorf("Error unmarshalling most recent host version: %s", err)
	}
	h.RawDriver = driver.Data
	return &h, nil
}
