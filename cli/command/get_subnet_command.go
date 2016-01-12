package command

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetSubnetCommand() cli.Command {
	return cli.Command{
		Name:  "subnet",
		Usage: "Get subnet",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			getSubnetCommandFunc(c)
		},
	}
}

func getSubnetCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a subnet")
	}
	subnet := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))

	s, err := clnt.Subnet.Get(subnet)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", string(s.JSON()))
}
