package command

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/cli/prompt"
	"github.com/imc-trading/dock2box/client"
)

func NewCreateHostCommand() cli.Command {
	return cli.Command{
		Name:  "host",
		Usage: "Create host",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "prompt, p", Usage: "Prompt for resource input"},
			cli.BoolFlag{Name: "disable-build", Usage: "Disable PXE build, this prevents a host from being provisioned (enabled by default)"},
			cli.BoolFlag{Name: "debug", Usage: "Enable debug during host provisioning (disabled by default)"},
			cli.BoolFlag{Name: "gpt", Usage: "Enable use of GUID Partition Table (disabled by default)"},
			cli.StringFlag{Name: "tag, t", Value: "", Usage: "Tag to use for provisioning"},
			cli.StringFlag{Name: "version, v", Value: "latest", Usage: "Image version to use for provisioning"},
			cli.StringFlag{Name: "kopts, k", Usage: "Kernel options"},
			cli.StringFlag{Name: "tenant, T", Usage: "Tenant"},
			cli.StringFlag{Name: "labels, l", Usage: "Comma-separated list of labels"},
			cli.StringFlag{Name: "site, s", Usage: "Site"},
			cli.StringFlag{Name: "interface, I", Value: "eth0", Usage: "Interface"},
			cli.BoolFlag{Name: "dhcp, D", Usage: "DHCP"},
			cli.StringFlag{Name: "hwaddr, H", Usage: "Hardware address"},
			cli.StringFlag{Name: "ipv4, P", Usage: "IPv4 address"},
			cli.StringFlag{Name: "subnet, S", Usage: "Subnet address using prefix ex. 192.168.0.1/24"},
		},
		Action: func(c *cli.Context) {
			createHostCommandFunc(c)
		},
	}
}

/*
func addInterface(clnt *client.Client, siteID string) client.Interface {
	ifs := client.Interface{
		Interface: prompt.String("Interface", prompt.Prompt{Default: "eth0", FuncPtr: prompt.Regex, FuncInp: "^[a-z][a-z0-9]+$"}),
		DHCP:      prompt.Bool("DHCP", false),
		HwAddr:    prompt.String("Hardware Address", prompt.Prompt{NoDefault: true, FuncPtr: validateHwAddr}),
	}

	if !ifs.DHCP {
		ifs.IPv4 = prompt.String("IP Address", prompt.Prompt{NoDefault: true, FuncPtr: validateIPv4})
		ifs.SubnetID = *chooseSubnet(clnt, siteID, "")
		// Check subnet match IPv4
	}

	return ifs
}
*/

func createHostCommandFunc(c *cli.Context) {
	var hostname string
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a hostname")
	} else {
		hostname = c.Args()[0]
	}

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	if c.Bool("prompt") {
		h := client.Host{
			Host:  hostname,
			Build: prompt.Bool("Build", true),
			Debug: prompt.Bool("Debug", false),
			GPT:   prompt.Bool("GPT", false),
			TagID: *chooseTag(clnt, ""),
		}

		// Get labels
		labels := prompt.String("Comma-separated list of labels", prompt.Prompt{Default: "", FuncPtr: prompt.Regex, FuncInp: "^([a-zA-Z][a-zA-Z0-9-]+,)*([a-zA-Z][a-zA-Z0-9-]+)$"})
		if labels == "" {
			h.Labels = []string{}
		} else {
			h.Labels = strings.Split(labels, ",")
		}

		h.KOpts = prompt.String("KOpts", prompt.Prompt{Default: "", FuncPtr: prompt.Regex, FuncInp: "^(|[a-zA-Z0-9- ])+$"})
		h.TenantID = *chooseTenants(clnt, "")
		h.SiteID = *chooseSite(clnt, "")

		/*
			h.Interfaces = []client.Interface{addInterface(clnt, h.SiteID)}
			if prompt.Bool("Do you want to add another network interface", false) {
				h.Interfaces = append(h.Interfaces, addInterface(clnt, h.SiteID))
			}
		*/

		// Is this correct?
		fmt.Println(string(h.JSON()))
		if !prompt.Bool("Is this correct", true) {
			os.Exit(1)
		}

		// Create host
		clnt.Host.Create(&h)
		return
	}

	/*
		h := client.Host{
			Host:  hostname,
			Debug: c.Bool("debug"),
			GPT:   c.Bool("gpt"),
			KOpts: c.String("kopts"),
		}

		// Get build
		if c.Bool("disable-build") {
			h.Build = false
		} else {
			h.Build = true
		}

		// Check arguments
		if !c.IsSet("image") {
			log.Fatalf("You need to specify image")
		}

		if !c.IsSet("tenant") {
			log.Fatalf("You need to specify tenant")
		}

		if !c.IsSet("site") {
			log.Fatalf("You need to specify site")
		}

		// Get tag ID
		tag, err := clnt.Tag.Get(c.String("tag"))
		if err != nil {
			log.Fatalf(err.Error())
		}
		h.TagID = tag.ID

		// Get tenant
		tenant, err := clnt.Tenant.Get(c.String("tenant"))
		if err != nil {
			log.Fatalf(err.Error())
		}
		h.TenantID = tenant.ID

		// Get labels
		if !c.IsSet("labels") {
			h.Labels = []string{}
		} else {
			h.Labels = strings.Split(c.String("labels"), ",")
		}

		// Get site
		site, err := clnt.Site.Get(c.String("site"))
		if err != nil {
			log.Fatalf(err.Error())
		}
		h.SiteID = site.ID

		// Host Interface
		ifs := client.Interface{
			Interface: c.String("interface"),
			DHCP:      c.Bool("dhcp"),
			HwAddr:    c.String("hwaddr"),
		}

		if !ifs.DHCP {
			ifs.IPv4 = c.String("ipv4")

			// Get subnet
			subnet, err := clnt.Subnet.Get(c.String("subnet"))
			if err != nil {
				log.Fatalf(err.Error())
			}
			ifs.SubnetID = subnet.ID

			h.Interfaces = []client.Interface{ifs}
		}

		// Create host
		clnt.Host.Create(&h)
	*/
}
