package command

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/cli/prompt"
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
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	if !prompt.Bool("Are you sure you wan't to remove "+subnet, true) {
		os.Exit(1)
	}

	h, err := clnt.Subnet.Delete(subnet)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", string(h.JSON()))
}
