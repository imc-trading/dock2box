package command

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewDeleteSubnetCommand() cli.Command {
	return cli.Command{
		Name:  "subnet",
		Usage: "Delete subnet",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			deleteSubnetCommandFunc(c)
		},
	}
}

func deleteSubnetCommandFunc(c *cli.Context) {
	var subnet string
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a subnet")
	} else {
		subnet = c.Args()[0]
	}

	clnt := client.New(c.GlobalString("server"))

	h, err := clnt.Subnet.Delete(subnet)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", string(h.JSON()))
}
