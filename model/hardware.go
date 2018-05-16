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

func NewHardware() (*Hardware, []error) {
	h := &Hardware{}
	errs := []error{}

	var err error
	if h.Distro, err = hwinfo.GetDistro(); err != nil {
		errs = append(errs, err)
	}

	if h.CPU, err = hwinfo.GetCPU(); err != nil {
		errs = append(errs, err)
	}

	if h.Memory, err = hwinfo.GetMemory(); err != nil {
		errs = append(errs, err)
	}

	if h.DMI, err = hwinfo.GetDMI(); err != nil {
		errs = append(errs, err)
	}

	if h.IPMI, err = hwinfo.GetIPMI(); err != nil {
		errs = append(errs, err)
	}

	if h.DNS, err = hwinfo.GetDNS(); err != nil {
		errs = append(errs, err)
	}

	if h.PCICards, err = hwinfo.GetPCICards(); err != nil {
		errs = append(errs, err)
	}

	if h.Disks, err = hwinfo.GetDisks(); err != nil {
		errs = append(errs, err)
	}

	if h.Interfaces, err = hwinfo.GetInterfaces(); err != nil {
		errs = append(errs, err)
	}

	if h.Routes, err = hwinfo.GetRoutes(); err != nil {
		errs = append(errs, err)
	}

	return h, errs
}
