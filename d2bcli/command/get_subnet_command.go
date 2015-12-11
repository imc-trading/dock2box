package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetSubnetCommand() cli.Command {
	return cli.Command{
		Name:  "subnet",
		Usage: "Get subnet",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "all, a", Usage: "Get all subnets"},
		},
		Action: func(c *cli.Context) {
			getSubnetCommandFunc(c)
		},
	}
}

func getSubnetCommandFunc(c *cli.Context) {
	var subnet string
	if !c.Bool("all") {
		if len(c.Args()) == 0 {
			log.Fatal("You need to specify a subnet")
		} else {
			subnet = c.Args()[0]
		}
	}

	clnt := client.New(c.GlobalString("server"))

	if c.Bool("all") {
		s, err := clnt.Subnet.All()
		if err != nil {
			log.Fatal(err.Error())
		}
		b, _ := json.MarshalIndent(s, "", "  ")
		fmt.Printf("%v\n", string(b))
	} else {
		s, err := clnt.Subnet.Get(subnet)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%v\n", string(s.JSON()))
	}
}
