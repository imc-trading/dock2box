package command

import (
	"log"
	"strings"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
	"github.com/imc-trading/dock2box/d2bcli/prompt"
)

func NewUpdateSiteCommand() cli.Command {
	return cli.Command{
		Name:  "site",
		Usage: "Update site",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "prompt, p", Usage: "Prompt for resource input"},
			cli.StringFlag{Name: "domain, d", Usage: "Domain"},
			cli.StringFlag{Name: "dns, D", Usage: "Comma-separated list of dns servers"},
			cli.StringFlag{Name: "docker-registry, r", Value: "registry", Usage: "Docker Registry for site"},
			cli.StringFlag{Name: "artifact-repository, a", Value: "repository", Usage: "Artifact repository for site"},
			cli.StringFlag{Name: "naming-scheme, n", Value: "hardware-address", Usage: "Naming scheme (serial-number, hardware-address, external)"},
			cli.StringFlag{Name: "pxe-theme, t", Value: "night", Usage: "PXE Theme (night, matrix)"},
		},
		Action: func(c *cli.Context) {
			updateSiteCommandFunc(c)
		},
	}
}

func updateSiteCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a site")
	}
	site := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	v, err := clnt.Site.Get(site)
	if err != nil {
		log.Fatal(err.Error())
	}

	if c.Bool("prompt") {
		s := client.Site{
			ID:                 v.ID,
			Site:               prompt.String("Site", prompt.Prompt{Default: v.Site, FuncPtr: prompt.Regex, FuncInp: ""}),
			Domain:             prompt.String("Domain", prompt.Prompt{Default: v.Domain, FuncPtr: prompt.Regex, FuncInp: ""}),
			DNS:                strings.Split(prompt.String("DNS", prompt.Prompt{Default: strings.Join(v.DNS, ","), FuncPtr: validateIPv4List, FuncInp: ""}), ","),
			DockerRegistry:     prompt.String("Docker Registry", prompt.Prompt{Default: v.DockerRegistry, FuncPtr: prompt.Regex, FuncInp: ""}),
			ArtifactRepository: prompt.String("Artifact Repository", prompt.Prompt{Default: v.ArtifactRepository, FuncPtr: prompt.Regex, FuncInp: ""}),
			NamingScheme:       prompt.String("Naming Scheme", prompt.Prompt{Default: v.NamingScheme, FuncPtr: prompt.Enum, FuncInp: "serial-number,hardware-address,external"}),
			PXETheme:           prompt.String("PXE Theme", prompt.Prompt{Default: v.PXETheme, FuncPtr: prompt.Enum, FuncInp: "night,matrix"}),
		}

		// Create site
		clnt.Site.Update(site, &s)
		return
	}
}
