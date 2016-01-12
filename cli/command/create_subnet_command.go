package command

import (
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
	"github.com/imc-trading/dock2box/cli/prompt"
)

func NewCreateSubnetCommand() cli.Command {
	return cli.Command{
		Name:  "subnet",
		Usage: "Create subnet",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "prompt, p", Usage: "Prompt for resource input"},
			cli.StringFlag{Name: "mask, m", Usage: "Mask"},
			cli.StringFlag{Name: "gateway, g", Usage: "Gateway"},
			cli.StringFlag{Name: "site, s", Usage: "Site"},
		},
		Action: func(c *cli.Context) {
			createSubnetCommandFunc(c)
		},
	}
}

func createSubnetCommandFunc(c *cli.Context) {
	var subnet string
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a subnet")
	} else {
		subnet = c.Args()[0]
	}

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	if c.Bool("prompt") {
		s := client.Subnet{
			Subnet: subnet,
			// Calculate automatically based on subnet/prefix
			Mask: prompt.String("Mask", prompt.Prompt{NoDefault: true, FuncPtr: validateIPv4, FuncInp: ""}),
			// Default to .254 for subnet
			Gw:     prompt.String("Gateway", prompt.Prompt{NoDefault: true, FuncPtr: validateIPv4, FuncInp: ""}),
			SiteID: *chooseSite(clnt, ""),
		}

		// Create subnet
		clnt.Subnet.Create(&s)
		return
	}
}
