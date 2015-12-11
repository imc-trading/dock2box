// +build linux

package hwinfo

import (
	"os"
	"strings"
	"time"

	"github.com/mickep76/hwinfo/cpu"
	"github.com/mickep76/hwinfo/disks"
	"github.com/mickep76/hwinfo/dock2box"
	"github.com/mickep76/hwinfo/lvm"
	"github.com/mickep76/hwinfo/memory"
	"github.com/mickep76/hwinfo/mounts"
	"github.com/mickep76/hwinfo/network"
	"github.com/mickep76/hwinfo/opsys"
	"github.com/mickep76/hwinfo/pci"
	"github.com/mickep76/hwinfo/routes"
	"github.com/mickep76/hwinfo/sysctl"
	"github.com/mickep76/hwinfo/system"
)

// HWInfo information.
type HWInfo struct {
	Hostname      string             `json:"hostname"`
	ShortHostname string             `json:"short_hostname"`
	CPU           *cpu.CPU           `json:"cpu"`
	Memory        *memory.Memory     `json:"memory"`
	OpSys         *opsys.OpSys       `json:"opsys"`
	System        *system.System     `json:"system"`
	Network       *network.Network   `json:"network"`
	PCI           *[]pci.PCI         `json:"pci"`
	Disks         *[]disks.Disk      `json:"disks"`
	Routes        *[]routes.Route    `json:"routes"`
	Sysctl        *[]sysctl.Sysctl   `json:"sysctl"`
	LVM           *lvm.LVM           `json:"lvm"`
	Mounts        *[]mounts.Mount    `json:"mounts"`
	Dock2Box      *dock2box.Dock2Box `json:"dock2box"`

	cpuTTL      int
	memoryTTL   int
	opSysTTL    int
	systemTTL   int
	networkTTL  int
	pciTTL      int
	disksTTL    int
	routesTTL   int
	sysctlTTL   int
	lvmTTL      int
	mountsTTL   int
	dock2boxTTL int
	last        time.Time
}

func NewHWInfo() *HWInfo {
	return &HWInfo{
		cpuTTL:      24 * 60 * 60, // Every 24 hours
		memoryTTL:   24 * 60 * 60, // Every 24 hours
		opSysTTL:    60 * 60,      // Every hour
		systemTTL:   60 * 60,      // Every hour
		networkTTL:  60 * 60,      // Every hour
		pciTTL:      60 * 60 * 2,  // Every other hour
		disksTTL:    60 * 60,      // Every hour
		routesTTL:   60 * 60,      // Every hour
		sysctlTTL:   60 * 60,      // Every hour
		lvmTTL:      60 * 60,      // Every hour
		mountsTTL:   60 * 60,      // Every hour
		dock2boxTTL: 60 * 60,      // Every hour
	}
}

func (hwi *HWInfo) TTL(cpu int, memory int, opSys int, system int, network int, pci int, disks int, routes int, sysctl int, lvm int, mounts int, dock2box int) {
	hwi.cpuTTL = cpu
	hwi.memoryTTL = memory
	hwi.opSysTTL = opSys
	hwi.systemTTL = system
	hwi.networkTTL = network
	hwi.pciTTL = pci
	hwi.disksTTL = disks
	hwi.routesTTL = routes
	hwi.sysctlTTL = sysctl
	hwi.lvmTTL = lvm
	hwi.mountsTTL = mounts
	hwi.dock2boxTTL = dock2box
}

// Get information about a system.
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
		i, err := network.Get()
		if err != nil {
			return err
		}
		hwi.Network = &i
	}

	ttl = now
	ttl.Add(time.Duration(hwi.pciTTL) * time.Second)
	if hwi.PCI == nil || hwi.last.Before(ttl) {
		i, err := pci.Get()
		if err != nil {
			return err
		}
		hwi.PCI = &i
	}

	ttl = now
	ttl.Add(time.Duration(hwi.disksTTL) * time.Second)
	if hwi.Disks == nil || hwi.last.Before(ttl) {
		i, err := disks.Get()
		if err != nil {
			return err
		}
		hwi.Disks = &i
	}

	ttl = now
	ttl.Add(time.Duration(hwi.routesTTL) * time.Second)
	if hwi.Routes == nil || hwi.last.Before(ttl) {
		i, err := routes.Get()
		if err != nil {
			return err
		}
		hwi.Routes = &i
	}

	ttl = now
	ttl.Add(time.Duration(hwi.sysctlTTL) * time.Second)
	if hwi.Sysctl == nil || hwi.last.Before(ttl) {
		i, err := sysctl.Get()
		if err != nil {
			return err
		}
		hwi.Sysctl = &i
	}

	ttl = now
	ttl.Add(time.Duration(hwi.lvmTTL) * time.Second)
	if hwi.LVM == nil || hwi.last.Before(ttl) {
		i, err := lvm.Get()
		if err != nil {
			return err
		}
		hwi.LVM = &i
	}

	ttl = now
	ttl.Add(time.Duration(hwi.mountsTTL) * time.Second)
	if hwi.Mounts == nil || hwi.last.Before(ttl) {
		i, err := mounts.Get()
		if err != nil {
			return err
		}
		hwi.Mounts = &i
	}

	// Don't bail just because Dock2Box is not available,
	//  need a better way to handle this maybe use log package or return different error levels
	ttl = now
	ttl.Add(time.Duration(hwi.dock2boxTTL) * time.Second)
	if hwi.Dock2Box == nil || hwi.last.Before(ttl) {
		i, _ := dock2box.Get()
		hwi.Dock2Box = &i
	}

	return nil
}

// Get information about a system.
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

	i7, err := pci.Get()
	if err != nil {
		return err
	}
	hwi.PCI = &i7

	i8, err := disks.Get()
	if err != nil {
		return err
	}
	hwi.Disks = &i8

	i9, err := routes.Get()
	if err != nil {
		return err
	}
	hwi.Routes = &i9

	i10, err := sysctl.Get()
	if err != nil {
		return err
	}
	hwi.Sysctl = &i10

	i11, err := lvm.Get()
	if err != nil {
		return err
	}
	hwi.LVM = &i11

	i12, err := mounts.Get()
	if err != nil {
		return err
	}
	hwi.Mounts = &i12

	// Don't bail just because Dock2Box is not available,
	//  need a better way to handle this maybe use log package or return different error levels
	i13, _ := dock2box.Get()
	hwi.Dock2Box = &i13

	return nil
}
