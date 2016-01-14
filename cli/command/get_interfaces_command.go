package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetInterfacesCommand() cli.Command {
	return cli.Command{
		Name:  "interfaces",
		Usage: "Get all interfaces",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			getInterfacesCommandFunc(c)
		},
	}
}

func getInterfacesCommandFunc(c *cli.Context) {
	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	i, err := clnt.Interface.All()
	if err != nil {
		log.Fatal(err.Error())
	}
	b, _ := json.MarshalIndent(i, "", "  ")
	fmt.Printf("%v\n", string(b))
}
