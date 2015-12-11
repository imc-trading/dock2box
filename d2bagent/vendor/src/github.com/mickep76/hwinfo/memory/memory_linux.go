// +build linux

package memory

import (
	"github.com/mickep76/hwinfo/common"
	"strconv"
	"strings"
)

// Get information about system memory.
func Get() (Memory, error) {
	m := Memory{}

	o, err := common.LoadFileFields("/proc/meminfo", ":", []string{
		"MemTotal",
	})
	if err != nil {
		return Memory{}, err
	}

	m.TotalGB, err = strconv.Atoi(strings.TrimRight(o["MemTotal"], " kB"))
	if err != nil {
		return Memory{}, err
	}
	m.TotalGB = m.TotalGB / 1024 / 1024

	return m, nil
}
