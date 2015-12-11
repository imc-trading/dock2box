// +build linux

package opsys

import (
	"github.com/mickep76/hwinfo/common"
	"runtime"
)

// Get information about the operating system.
func Get() (OpSys, error) {
	opsys := OpSys{}

	o, err := common.ExecCmdFields("lsb_release", []string{"-a"}, ":", []string{
		"Distributor ID",
		"Release",
	})
	if err != nil {
		return OpSys{}, err
	}

	opsys.Kernel = runtime.GOOS
	opsys.Product = o["Distributor ID"]
	opsys.ProductVersion = o["Release"]

	opsys.KernelVersion, err = common.ExecCmd("uname", []string{"-r"})
	if err != nil {
		return OpSys{}, err
	}

	return opsys, nil
}
