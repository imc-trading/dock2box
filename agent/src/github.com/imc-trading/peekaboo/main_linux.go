// +build linux

package main

import (
	"github.com/Unknwon/macaron"
	"github.com/mickep76/hwinfo"
)

func routes(m *macaron.Macaron, hw *hwinfo.HWInfo) {
	// HTML endpoints
	m.Get("/", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "Peekaboo"
		ctx.Data["Kernel"] = hw.OpSys.Kernel
		ctx.Data["Version"] = Version
		ctx.Data["Hostname"] = hw.Hostname
		ctx.Data["ShortHostname"] = hw.ShortHostname
		ctx.Data["CPU"] = hw.CPU
		ctx.Data["Memory"] = hw.Memory
		ctx.Data["OpSys"] = hw.OpSys
		ctx.Data["System"] = hw.System

		ctx.HTML(200, "peekaboo")
	})

	m.Get("/network", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "Network"
		ctx.Data["Kernel"] = hw.OpSys.Kernel
		ctx.Data["ShortHostname"] = hw.ShortHostname
		ctx.HTML(200, "network")
	})

	m.Get("/storage", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "Storage"
		ctx.Data["Kernel"] = hw.OpSys.Kernel
		ctx.Data["ShortHostname"] = hw.ShortHostname
		ctx.HTML(200, "storage")
	})

	m.Get("/pci", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "PCI"
		ctx.Data["Kernel"] = hw.OpSys.Kernel
		ctx.Data["ShortHostname"] = hw.ShortHostname
		ctx.HTML(200, "pci")
	})

	m.Get("/sysctl", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "Sysctl"
		ctx.Data["Kernel"] = hw.OpSys.Kernel
		ctx.Data["ShortHostname"] = hw.ShortHostname
		ctx.HTML(200, "sysctl")
	})

	m.Get("/dock2box", func(ctx *macaron.Context) {
		ctx.Data["Title"] = "Dock2Box"
		ctx.Data["Kernel"] = hw.OpSys.Kernel
		ctx.Data["ShortHostname"] = hw.ShortHostname
		ctx.Data["Dock2Box"] = hw.Dock2Box
		ctx.HTML(200, "dock2box")
	})

	// JSON endpoints
	m.Get("/disks/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.Disks)
	})

	m.Get("/mounts/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.Mounts)
	})

	m.Get("/network/routes/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.Routes)
	})

	m.Get("/sysctl/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.Sysctl)
	})

	m.Get("/lvm/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.LVM)
	})

	m.Get("/lvm/phys_vols/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.LVM.PhysVols)
	})

	m.Get("/lvm/log_vols/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.LVM.LogVols)
	})

	m.Get("/lvm/vol_grps/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.LVM.VolGrps)
	})

	m.Get("/pci/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.PCI)
	})

	m.Get("/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw)
	})

	m.Get("/cpu/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.CPU)
	})

	m.Get("/memory/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.Memory)
	})

	m.Get("/network/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.Network)
	})

	m.Get("/network/interfaces/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.Network.Interfaces)
	})

	m.Get("/opsys/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.OpSys)
	})

	m.Get("/network/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.Network)
	})

	m.Get("/dock2box/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.Dock2Box)
	})

	m.Get("/dock2box/layers/json", func(ctx *macaron.Context) {
		ctx.JSON(200, &hw.Dock2Box.Layers)
	})
}
