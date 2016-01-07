package command

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
	"github.com/imc-trading/dock2box/d2bcli/prompt"
)

func NewUpdateHostCommand() cli.Command {
	return cli.Command{
		Name:  "host",
		Usage: "Update host",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "prompt, p", Usage: "Prompt for resource input"},
			cli.BoolFlag{Name: "disable-build", Usage: "Disable PXE build, this prevents a host from being provisioned (enabled by default)"},
			cli.BoolFlag{Name: "debug", Usage: "Enable debug during host provisioning (disabled by default)"},
			cli.BoolFlag{Name: "gpt", Usage: "Enable use of GUID Partition Table (disabled by default)"},
			cli.StringFlag{Name: "image, i", Value: "", Usage: "Image to use for provisioning"},
			cli.StringFlag{Name: "version, v", Value: "latest", Usage: "Image version to use for provisioning"},
			cli.StringFlag{Name: "kopts, k", Usage: "Kernel options"},
			cli.StringFlag{Name: "tenant, t", Usage: "Tenant"},
			cli.StringFlag{Name: "labels, l", Usage: "Comma-separated list of labels"},
			cli.StringFlag{Name: "site, s", Usage: "Site"},
			cli.StringFlag{Name: "interface, I", Value: "eth0", Usage: "Interface"},
			cli.BoolFlag{Name: "dhcp, D", Usage: "DHCP"},
			cli.StringFlag{Name: "hwaddr, H", Usage: "Hardware address"},
			cli.StringFlag{Name: "ipv4, P", Usage: "IPv4 address"},
			cli.StringFlag{Name: "subnet, S", Usage: "Subnet address using prefix ex. 192.168.0.1/24"},
		},
		Action: func(c *cli.Context) {
			updateHostCommandFunc(c)
		},
	}
}

/*
func updateHostInterface(clnt *client.Client, siteID string, v client.HostInterface) client.HostInterface {
	ifs := client.HostInterface{
		Interface: prompt.String("Interface", prompt.Prompt{Default: v.Interface, FuncPtr: prompt.Regex, FuncInp: "^[a-z][a-z0-9]+$"}),
		DHCP:      prompt.Bool("DHCP", v.DHCP),
		HwAddr:    prompt.String("Hardware Address", prompt.Prompt{Default: v.HwAddr, FuncPtr: validateHwAddr}),
	}

	if !ifs.DHCP {
		ifs.IPv4 = prompt.String("IP Address", prompt.Prompt{Default: v.IPv4, FuncPtr: validateIPv4})
		ifs.SubnetID = *chooseSubnet(clnt, siteID, v.SubnetID)
	}

	return ifs
}
*/

func updateHostCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a hostname")
	}
	hostname := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	v, err := clnt.Host.Get(hostname)
	if err != nil {
		log.Fatal(err.Error())
	}

	if c.Bool("prompt") {
		h := client.Host{
			ID:    v.ID,
			Host:  prompt.String("Host", prompt.Prompt{Default: v.Host, FuncPtr: prompt.Regex, FuncInp: ""}),
			Build: prompt.Bool("Build", v.Build),
			Debug: prompt.Bool("Debug", v.Debug),
			GPT:   prompt.Bool("GPT", v.GPT),
			TagID: *chooseTag(clnt, v.TagID),
		}

		// Get labels
		labels := prompt.String("Comma-separated list of labels", prompt.Prompt{Default: strings.Join(v.Labels, ","), FuncPtr: prompt.Regex, FuncInp: "^([a-zA-Z][a-zA-Z0-9-]+,)*([a-zA-Z][a-zA-Z0-9-]+)$"})
		if labels == "" {
			h.Labels = []string{}
		} else {
			h.Labels = strings.Split(labels, ",")
		}

		h.KOpts = prompt.String("KOpts", prompt.Prompt{Default: v.KOpts, FuncPtr: prompt.Regex, FuncInp: "^(|[a-zA-Z0-9- ])+$"})
		h.TenantID = *chooseTenants(clnt, v.TenantID)
		h.SiteID = *chooseSite(clnt, v.SiteID)

		/*
			for _, v := range v.Interfaces {
				h.Interfaces = append(h.Interfaces, updateHostInterface(clnt, h.SiteID, v))
			}
		*/

		// Is this correct?
		fmt.Println(string(h.JSON()))
		if !prompt.Bool("Is this correct", true) {
			os.Exit(1)
		}

		// Update host
		clnt.Host.Update(hostname, &h)
		return
	}
}
