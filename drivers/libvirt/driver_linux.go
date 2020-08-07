package libvirt

import (
	"github.com/code-ready/machine/libmachine/drivers"
)

type Driver struct {
	*drivers.BaseDriver

	// Driver specific configuration
	Memory      int
	CPU         int
	Network     string
	DiskPath    string
	CacheMode   string
	IOMode      string
	DiskPathURL string
	SSHKeyPath  string
}

const (
	defaultMemory    = 8192
	defaultCPU       = 4
	defaultCacheMode = "default"
	defaultIOMode    = "threads"
)

func NewDriver(hostName, storePath string) *Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
		Memory:    defaultMemory,
		CPU:       defaultCPU,
		CacheMode: defaultCacheMode,
		IOMode:    defaultIOMode,
	}
}
