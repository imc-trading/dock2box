// +build darwin

package memory

import (
	"github.com/mickep76/hwinfo/common"
	"strconv"
)

// Get information about system memory.
func Get() (Memory, error) {
	m := Memory{}

	o, err := common.ExecCmdFields("/usr/sbin/sysctl", []string{"-a"}, ":", []string{
		"hw.memsize",
	})
	if err != nil {
		return Memory{}, err
	}

	m.TotalGB, err = strconv.Atoi(o["hw.memsize"])
	if err != nil {
		return Memory{}, err
	}
	m.TotalGB = m.TotalGB / 1024 / 1024 / 1024

	return m, nil
}
