package command

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
	"github.com/imc-trading/dock2box/d2bcli/prompt"
)

func NewCreateSiteCommand() cli.Command {
	return cli.Command{
		Name:  "site",
		Usage: "Create site",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "prompt, p", Usage: "Prompt for resource input"},
			cli.StringFlag{Name: "domain, d", Usage: "Domain"},
			cli.StringFlag{Name: "dns, D", Usage: "Comma-separated list of dns servers"},
			cli.StringFlag{Name: "docker-registry, r", Value: "registry", Usage: "Docker Registry for site"},
			cli.StringFlag{Name: "artifact-repository, a", Value: "repository", Usage: "Artifact repository for site"},
			cli.StringFlag{Name: "naming-scheme, n", Value: "repository", Usage: "Naming scheme (serial-number, hardware-address, external)"},
		},
		Action: func(c *cli.Context) {
			createSiteCommandFunc(c)
		},
	}
}

func validateIPv4List(inp string, list string) bool {
	for _, v := range strings.Split(inp, ",") {
		if net.ParseIP(v) == nil {
			fmt.Println("Invalid IPv4 address: %s", v)
			return false
		}
	}
	return true
}

func createSiteCommandFunc(c *cli.Context) {
	var site string
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a site")
	} else {
		site = c.Args()[0]
	}

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	if c.Bool("prompt") {
		s := client.Site{
			Site:               site,
			Domain:             prompt.String("Domain", prompt.Prompt{NoDefault: true, FuncPtr: prompt.Regex, FuncInp: ""}),
			DNS:                strings.Split(prompt.String("DNS", prompt.Prompt{NoDefault: true, FuncPtr: validateIPv4List, FuncInp: ""}), ","),
			DockerRegistry:     prompt.String("Docker Registry", prompt.Prompt{NoDefault: true, FuncPtr: prompt.Regex, FuncInp: ""}),
			ArtifactRepository: prompt.String("Artifact Repository", prompt.Prompt{NoDefault: true, FuncPtr: prompt.Regex, FuncInp: ""}),
			// Def hwaddr
			NamingScheme: prompt.String("Naming Scheme", prompt.Prompt{NoDefault: true, FuncPtr: prompt.Enum, FuncInp: "serial-number,hardware-address,external"}),
		}

		// Create site
		clnt.Site.Create(&s)
		return
	}
}
