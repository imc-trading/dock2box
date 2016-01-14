package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetSubnetsCommand() cli.Command {
	return cli.Command{
		Name:  "subnets",
		Usage: "Get all subnets",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			getSubnetsCommandFunc(c)
		},
	}
}

func getSubnetsCommandFunc(c *cli.Context) {
	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	s, err := clnt.Subnet.All()
	if err != nil {
		log.Fatal(err.Error())
	}
	b, _ := json.MarshalIndent(s, "", "  ")
	fmt.Printf("%v\n", string(b))
}
