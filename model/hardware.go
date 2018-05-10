package model

import "github.com/mickep76/hwinfo"

type Hardware struct {
	PCICards   hwinfo.PCICards   `json:"pciCards,omitempty"`
	Disks      hwinfo.Disks      `json:"disks,omitempty"`
	Interfaces hwinfo.Interfaces `json:"interfaces,omitempty"`
	Routes     hwinfo.Routes     `json:"routes,omitempty"`

	*hwinfo.Distro
	*hwinfo.CPU
	*hwinfo.Memory
	*hwinfo.DMI
	*hwinfo.IPMI
	*hwinfo.DNS
}
