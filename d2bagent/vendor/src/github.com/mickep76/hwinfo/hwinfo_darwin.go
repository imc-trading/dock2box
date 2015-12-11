package hwinfo

import (
	"os"
	"strings"
	"time"

	"github.com/mickep76/hwinfo/cpu"
	"github.com/mickep76/hwinfo/memory"
	"github.com/mickep76/hwinfo/network"
	"github.com/mickep76/hwinfo/opsys"
	"github.com/mickep76/hwinfo/system"
)

// HWInfo information.
type HWInfo struct {
	Hostname      string           `json:"hostname"`
	ShortHostname string           `json:"short_hostname"`
	CPU           *cpu.CPU         `json:"cpu"`
	Memory        *memory.Memory   `json:"memory"`
	OpSys         *opsys.OpSys     `json:"opsys"`
	System        *system.System   `json:"system"`
	Network       *network.Network `json:"network"`

	cpuTTL     int
	memoryTTL  int
	opSysTTL   int
	systemTTL  int
	networkTTL int
	last       time.Time
}

func NewHWInfo() *HWInfo {
	return &HWInfo{
		cpuTTL:     24 * 60 * 60, // Every 24 hours
		memoryTTL:  24 * 60 * 60, // Every 24 hours
		opSysTTL:   60 * 60,      // Every hour
		systemTTL:  60 * 60,      // Every hour
		networkTTL: 60 * 60,      // Every hour
	}
}

func (hwi *HWInfo) TTL(cpu int, memory int, opSys int, system int, network int) {
	hwi.cpuTTL = cpu
	hwi.memoryTTL = memory
	hwi.opSysTTL = opSys
	hwi.systemTTL = system
	hwi.networkTTL = network
}

// Get information about a system with no TTL.
func (hwi *HWInfo) Get() error {
	host, err := os.Hostname()
	if err != nil {
		return err
	}
	hwi.Hostname = host
	hwi.ShortHostname = strings.Split(host, ".")[0]

	i2, err := cpu.Get()
	if err != nil {
		return err
	}
	hwi.CPU = &i2

	i3, err := memory.Get()
	if err != nil {
		return err
	}
	hwi.Memory = &i3

	i4, err := opsys.Get()
	if err != nil {
		return err
	}
	hwi.OpSys = &i4

	i5, err := system.Get()
	if err != nil {
		return err
	}
	hwi.System = &i5

	i6, err := network.Get()
	if err != nil {
		return err
	}
	hwi.Network = &i6

	return nil
}

// Get information about a system with TTL.
func (hwi *HWInfo) GetTTL() error {
	host, err := os.Hostname()
	if err != nil {
		return err
	}
	hwi.Hostname = host
	hwi.ShortHostname = strings.Split(host, ".")[0]

	now := time.Now()
	ttl := now
	ttl.Add(time.Duration(hwi.cpuTTL) * time.Second)
	if hwi.CPU == nil || hwi.last.Before(ttl) {
		i, err := cpu.Get()
		if err != nil {
			return err
		}
		hwi.CPU = &i
	}

	ttl = now
	ttl.Add(time.Duration(hwi.memoryTTL) * time.Second)
	if hwi.Memory == nil || hwi.last.Before(ttl) {
		i, err := memory.Get()
		if err != nil {
			return err
		}
		hwi.Memory = &i
	}

	ttl = now
	ttl.Add(time.Duration(hwi.opSysTTL) * time.Second)
	if hwi.OpSys == nil || hwi.last.Before(ttl) {
		i, err := opsys.Get()
		if err != nil {
			return err
		}
		hwi.OpSys = &i
	}

	ttl = now
	ttl.Add(time.Duration(hwi.systemTTL) * time.Second)
	if hwi.System == nil || hwi.last.Before(ttl) {
		i, err := system.Get()
		if err != nil {
			return err
		}
		hwi.System = &i
	}

	ttl = now
	ttl.Add(time.Duration(hwi.networkTTL) * time.Second)
	if hwi.Network == nil || hwi.last.Before(ttl) {
		i6, err := network.Get()
		if err != nil {
			return err
		}
		hwi.Network = &i6
	}

	return nil
}
