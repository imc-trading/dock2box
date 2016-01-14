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

func NewCreateSiteCommand() cli.Command {
	return cli.Command{
		Name:  "site",
		Usage: "Create site",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			createSiteCommandFunc(c)
		},
	}
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

	s := client.Site{
		Site:               site,
		Domain:             prompt.String("Domain", prompt.Prompt{NoDefault: true, FuncPtr: prompt.Regex, FuncInp: ""}),
		DNS:                strings.Split(prompt.String("DNS", prompt.Prompt{NoDefault: true, FuncPtr: validateIPv4List, FuncInp: ""}), ","),
		DockerRegistry:     prompt.String("Docker Registry", prompt.Prompt{NoDefault: true, FuncPtr: prompt.Regex, FuncInp: ""}),
		ArtifactRepository: prompt.String("Artifact Repository", prompt.Prompt{NoDefault: true, FuncPtr: prompt.Regex, FuncInp: ""}),
		NamingScheme:       prompt.String("Naming Scheme", prompt.Prompt{Default: "hardware-address", FuncPtr: prompt.Enum, FuncInp: "serial-number,hardware-address,external"}),
		PXETheme:           prompt.String("PXE Theme", prompt.Prompt{Default: "night", FuncPtr: prompt.Regex, FuncInp: ""}),
	}

	// Is this correct?
	fmt.Println(string(s.JSON()))
	if !prompt.Bool("Is this correct", true) {
		os.Exit(1)
	}

	// Create site
	clnt.Site.Create(&s)
}
