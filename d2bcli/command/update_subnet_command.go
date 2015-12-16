package command

import (
	"log"
	"strings"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
	"github.com/imc-trading/dock2box/d2bcli/prompt"
)

func NewUpdateSubnetCommand() cli.Command {
	return cli.Command{
		Name:  "subnet",
		Usage: "Update subnet",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "prompt, p", Usage: "Prompt for resource input"},
			cli.StringFlag{Name: "mask, m", Usage: "Mask"},
			cli.StringFlag{Name: "gateway, g", Usage: "Gateway"},
			cli.StringFlag{Name: "site, s", Usage: "Site"},
		},
		Action: func(c *cli.Context) {
			updateSubnetCommandFunc(c)
		},
	}
}

// Get existing subnet as default for prompt
// Override with args if they are set

func updateSubnetCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a subnet")
	}
	subnet := strings.Replace(c.Args()[0], "/", "-", 1)

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	s, err := clnt.Subnet.Get(subnet)
	if err != nil {
		log.Fatal(err.Error())
	}

	if c.Bool("prompt") {
		s := client.Subnet{
			ID:     s.ID,
			Subnet: strings.Replace(subnet, "-", "/", 1),
			Mask:   prompt.String("Mask", prompt.Prompt{Default: s.Mask, FuncPtr: validateIPv4, FuncInp: ""}),
			Gw:     prompt.String("Gateway", prompt.Prompt{Default: s.Gw, FuncPtr: validateIPv4, FuncInp: ""}),
			SiteID: *chooseSite(clnt, s.SiteID),
		}

		// Create subnet
		clnt.Subnet.Update(subnet, &s)
		return
	}
}
