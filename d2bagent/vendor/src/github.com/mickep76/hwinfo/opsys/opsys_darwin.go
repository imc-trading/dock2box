// +build darwin

package opsys

import (
	"github.com/mickep76/hwinfo/common"
	"runtime"
)

// Get information about the operating system.
func Get() (OpSys, error) {
	opsys := OpSys{}

	o, err := common.ExecCmdFields("/usr/bin/sw_vers", []string{}, ":", []string{
		"ProductName",
		"ProductVersion",
	})
	if err != nil {
		return OpSys{}, err
	}

	opsys.Kernel = runtime.GOOS
	opsys.Product = o["ProductName"]
	opsys.ProductVersion = o["ProductVersion"]

	opsys.KernelVersion, err = common.ExecCmd("uname", []string{"-r"})
	if err != nil {
		return OpSys{}, err
	}

	return opsys, nil
}
