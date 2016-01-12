// +build darwin

package main

import (
	"github.com/Unknwon/macaron"
	"github.com/mickep76/hwinfo"
)

func routes(m *macaron.Macaron, hw *hwinfo.HWInfo) {
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
}
