// +build darwin

package system

import (
	"github.com/mickep76/hwinfo/common"
)

// Get information about a system.
func Get() (System, error) {
	s := System{}
	s.Manufacturer = "Apple Inc."

	o, err := common.ExecCmdFields("/usr/sbin/system_profiler", []string{"SPHardwareDataType"}, ":", []string{
		"Model Name",
		"Model Identifier",
		"Serial Number",
		"Boot ROM Version",
		"SMC Version",
	})
	if err != nil {
		return System{}, err
	}

	s.Product = o["Model Name"]
	s.ProductVersion = o["Model Identifier"]
	s.SerialNumber = o["Serial Number"]
	s.BootROMVersion = o["Boot ROM Version"]
	s.SMCVersion = o["SMC Version"]

	return s, nil
}
