package command

import (
	//	"fmt"
	"log"
	"strings"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
	"github.com/imc-trading/dock2box/d2bcli/prompt"
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
			createHostCommandFunc(c)
		},
	}
}

func chooseImage(clnt *client.Client) *string {
	r, err := clnt.Image.All()
	if err != nil {
		log.Fatalf(err.Error())
	}

	images := *r
	var list []string
	for _, v := range images {
		list = append(list, v.Image+" type: "+v.Type)
	}
	return &images[prompt.Choice("Choose image", list)].ID
}

func chooseImageVersion(clnt *client.Client, id string) string {
	r, err := clnt.ImageVersion.AllByID(id)
	if err != nil {
		log.Fatalf(err.Error())
	}

	versions := *r
	var list []string
	for _, v := range versions {
		list = append(list, v.Version+", created: "+v.Created)
	}
	return versions[prompt.Choice("Choose image version", list)].Version
}

func chooseTenants(clnt *client.Client) *string {
	r, err := clnt.Tenant.All()
	if err != nil {
		log.Fatalf(err.Error())
	}

	tenants := *r
	var list []string
	for _, v := range tenants {
		list = append(list, v.Tenant)
	}
	return &tenants[prompt.Choice("Choose tenant", list)].ID
}

func chooseSite(clnt *client.Client) *string {
	r, err := clnt.Site.All()
	if err != nil {
		log.Fatalf(err.Error())
	}

	sites := *r
	var list []string
	for _, v := range sites {
		list = append(list, v.Site+", domain: "+v.Domain)
	}
	return &sites[prompt.Choice("Choose site", list)].ID
}

func chooseSubnet(clnt *client.Client, siteID string) *string {
	r, err := clnt.Subnet.All()
	if err != nil {
		log.Fatalf(err.Error())
	}

	subnets := *r
	var list []string
	for _, v := range subnets {
		// UGLY: keep until backend supports filters
		if v.SiteID == siteID {
			list = append(list, v.Subnet)
		}
	}
	return &subnets[prompt.Choice("Choose subnet", list)].ID
}

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
			Host:    hostname,
			Build:   prompt.Bool("Build", true),
			Debug:   prompt.Bool("Debug", false),
			GPT:     prompt.Bool("GPT", false),
			ImageID: *chooseImage(clnt),
		}

		// Get image version
		h.Version = chooseImageVersion(clnt, h.ImageID)

		// Get labels
		labels := prompt.String("Comma-separated list of labels", prompt.Prompt{Default: "", FuncPtr: prompt.Regex, FuncInp: "^([a-zA-Z][a-zA-Z0-9-]+,)*([a-zA-Z][a-zA-Z0-9-]+)$"})
		if labels == "" {
			h.Labels = []string{}
		} else {
			h.Labels = strings.Split(labels, ",")
		}

		h.KOpts = prompt.String("KOpts", prompt.Prompt{Default: "", FuncPtr: prompt.Regex, FuncInp: "^(|[a-zA-Z0-9- ])+$"})
		h.TenantID = *chooseTenants(clnt)
		h.SiteID = *chooseSite(clnt)

		h.Interfaces = []client.HostInterface{
			{
				Interface: prompt.String("Interface", prompt.Prompt{Default: "eth0", FuncPtr: prompt.Regex, FuncInp: "^[a-z][a-z0-9]+$"}),
				DHCP:      prompt.Bool("DHCP", false),
				HwAddr:    prompt.String("Hardware Address", prompt.Prompt{NoDefault: true, FuncPtr: prompt.Regex, FuncInp: "^([0-9a-f]{2}:){5}[0-9a-f]{2}$"}),
			},
		}

		if !h.Interfaces[0].DHCP {
			h.Interfaces[0].IPv4 = prompt.String("IP Address", prompt.Prompt{NoDefault: true, FuncPtr: prompt.Regex, FuncInp: "^([0-9]{1,3}\\.){3}[0-9]{1,3}$"})
			h.Interfaces[0].SubnetID = *chooseSubnet(clnt, h.SiteID)
		}

		// Create host
		clnt.Host.Create(&h)
		return
	}

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

	// Get image ID
	image, err := clnt.Image.Get(c.String("image"))
	if err != nil {
		log.Fatalf(err.Error())
	}
	h.ImageID = image.ID

	// Get version
	h.Version = c.String("version")

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

	// Create host
	clnt.Host.Create(&h)
}
