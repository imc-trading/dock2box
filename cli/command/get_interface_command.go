package command

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetInterfaceCommand() cli.Command {
	return cli.Command{
		Name:  "interface",
		Usage: "Get interface",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			getInterfaceCommandFunc(c)
		},
	}
}

func getInterfaceCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify an interface")
	}
	intf := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))

	i, err := clnt.Interface.Get(intf)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", string(i.JSON()))
}
