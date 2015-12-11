package network

import (
	"fmt"
	"github.com/mickep76/hwinfo/common"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

// Info structure for information about a systems network interfaces.
type Interface struct {
	Name            string   `json:"name"`
	MTU             int      `json:"mtu"`
	IPAddr          []string `json:"ipaddr"`
	HWAddr          string   `json:"hwaddr"`
	Flags           []string `json:"flags"`
	Driver          string   `json:"driver,omitempty"`
	DriverVersion   string   `json:"driver_version,omitempty"`
	FirmwareVersion string   `json:"firmware_version,omitempty"`
	PCIBus          string   `json:"pci_bus,omitempty"`
	PCIBusURL       string   `json:"pci_bus_url,omitempty"`
	SwChassisID     string   `json:"sw_chassis_id"`
	SwName          string   `json:"sw_name"`
	SwDescr         string   `json:"sw_descr"`
	SwPortID        string   `json:"sw_port_id"`
	SwPortDescr     string   `json:"sw_port_descr"`
	SwVLAN          string   `json:"sw_vlan"`
}

// Info structure for information about a systems network.
type Network struct {
	Interfaces    []Interface `json:"interfaces"`
	OnloadVersion string      `json:"onload_version,omitempty"`
}

// GetInfo return information about a systems memory.
func Get() (Network, error) {
	network := Network{}

	intfs, err := net.Interfaces()
	if err != nil {
		return Network{}, err
	}

	for _, intf := range intfs {
		// Skip loopback interfaces
		if intf.Flags&net.FlagLoopback != 0 {
			continue
		}
		/*
			// Skip interfaces that are down
			if intf.Flags&net.FlagUp == 0 {
				continue
			}
		*/

		addrs, err := intf.Addrs()
		if err != nil {
			return Network{}, err
		}

		nintf := Interface{
			Name:   intf.Name,
			HWAddr: intf.HardwareAddr.String(),
			MTU:    intf.MTU,
		}

		for _, addr := range addrs {
			nintf.IPAddr = append(nintf.IPAddr, addr.String())
		}

		if intf.Flags&net.FlagUp != 0 {
			nintf.Flags = append(nintf.Flags, "up")
		}
		if intf.Flags&net.FlagBroadcast != 0 {
			nintf.Flags = append(nintf.Flags, "broadcast")
		}
		if intf.Flags&net.FlagPointToPoint != 0 {
			nintf.Flags = append(nintf.Flags, "pointtopoint")
		}
		if intf.Flags&net.FlagMulticast != 0 {
			nintf.Flags = append(nintf.Flags, "multicast")
		}

		switch runtime.GOOS {
		case "linux":
			o, err := common.ExecCmdFields("/usr/sbin/ethtool", []string{"-i", intf.Name}, ":", []string{
				"driver",
				"version",
				"firmware-version",
				"bus-info",
			})
			if err != nil {
				return Network{}, err
			}

			nintf.Driver = o["driver"]
			nintf.DriverVersion = o["version"]
			nintf.FirmwareVersion = o["firmware-version"]
			if strings.HasPrefix(o["bus-info"], "0000:") {
				nintf.PCIBus = o["bus-info"]
				nintf.PCIBusURL = fmt.Sprintf("/pci/%v", o["bus-info"])
			}

			o2, err := common.ExecCmdFields("/usr/sbin/lldpctl", []string{intf.Name}, ":", []string{
				"ChassisID",
				"SysName",
				"SysDescr",
				"PortID",
				"PortDescr",
				"VLAN",
			})
			if err != nil {
				return Network{}, err
			}

			nintf.SwChassisID = o2["ChassisID"]
			nintf.SwName = o2["SysName"]
			nintf.SwDescr = o2["SysDescr"]
			nintf.SwPortID = o2["PortID"]
			nintf.SwPortDescr = o2["PortDescr"]
			nintf.SwVLAN = o2["VLAN"]
		}

		network.Interfaces = append(network.Interfaces, nintf)
	}

	switch runtime.GOOS {
	case "linux":
		_, err := exec.LookPath("onload")
		if err == nil {
			o, _ := common.ExecCmdFields("onload", []string{"--version"}, ":", []string{"Kernel module"})
			if err != nil {
				return Network{}, err
			}
			network.OnloadVersion = o["Kernel module"]
		}
	}

	return network, nil
}
