package command

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/cli/prompt"
	"github.com/imc-trading/dock2box/client"
)

func NewCreateSubnetCommand() cli.Command {
	return cli.Command{
		Name:  "subnet",
		Usage: "Create subnet",
		Flags: []cli.Flag{},
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

	s := client.Subnet{
		Subnet: subnet,
		// Calculate automatically based on subnet/prefix
		Mask: prompt.String("Mask", prompt.Prompt{NoDefault: true, FuncPtr: validateIPv4, FuncInp: ""}),
		// Default to .254 for subnet
		Gw:     prompt.String("Gateway", prompt.Prompt{NoDefault: true, FuncPtr: validateIPv4, FuncInp: ""}),
		SiteID: *chooseSite(clnt, ""),
	}

	// Is this correct?
	fmt.Println(string(s.JSON()))
	if !prompt.Bool("Is this correct", true) {
		os.Exit(1)
	}

	// Create subnet
	clnt.Subnet.Create(&s)
}
