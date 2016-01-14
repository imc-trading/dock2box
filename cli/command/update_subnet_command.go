package command

import (
	"log"
	"strings"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/cli/prompt"
	"github.com/imc-trading/dock2box/client"
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

	v, err := clnt.Subnet.Get(subnet)
	if err != nil {
		log.Fatal(err.Error())
	}

	if c.Bool("prompt") {
		s := client.Subnet{
			ID:     v.ID,
			Subnet: prompt.String("Subnet", prompt.Prompt{Default: v.Subnet, FuncPtr: prompt.Regex, FuncInp: ""}),
			Mask:   prompt.String("Mask", prompt.Prompt{Default: v.Mask, FuncPtr: validateIPv4, FuncInp: ""}),
			Gw:     prompt.String("Gateway", prompt.Prompt{Default: v.Gw, FuncPtr: validateIPv4, FuncInp: ""}),
			SiteID: *chooseSite(clnt, v.SiteID),
		}

		// Create subnet
		clnt.Subnet.Update(subnet, &s)
		return
	}
}
