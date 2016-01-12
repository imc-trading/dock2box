// +build darwin

package cpu

import (
	"github.com/mickep76/hwinfo/common"
	"strconv"
	"strings"
)

// Get information about system CPU(s).
func Get() (CPU, error) {
	c := CPU{}

	o, err := common.ExecCmdFields("/usr/sbin/sysctl", []string{"-a"}, ":", []string{
		"machdep.cpu.core_count",
		"hw.physicalcpu_max",
		"hw.logicalcpu_max",
		"machdep.cpu.brand_string",
		"machdep.cpu.features",
	})
	if err != nil {
		return CPU{}, err
	}

	c.CoresPerSocket, err = strconv.Atoi(o["machdep.cpu.core_count"])
	if err != nil {
		return CPU{}, err
	}

	c.Physical, err = strconv.Atoi(o["hw.physicalcpu_max"])
	if err != nil {
		return CPU{}, err
	}

	c.Logical, err = strconv.Atoi(o["hw.logicalcpu_max"])
	if err != nil {
		return CPU{}, err
	}

	c.Sockets = c.Physical / c.CoresPerSocket
	c.ThreadsPerCore = c.Logical / c.Sockets / c.CoresPerSocket
	c.Model = o["machdep.cpu.brand_string"]
	c.Flags = strings.ToLower(o["machdep.cpu.features"])

	return c, nil
}
