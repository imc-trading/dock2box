package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetHostsCommand() cli.Command {
	return cli.Command{
		Name:  "hosts",
		Usage: "Get all hosts",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			getHostsCommandFunc(c)
		},
	}
}

func getHostsCommandFunc(c *cli.Context) {
	clnt := client.New(c.GlobalString("server"))

	h, err := clnt.Host.All()
	if err != nil {
		log.Fatal(err.Error())
	}
	b, _ := json.MarshalIndent(h, "", "  ")
	fmt.Printf("%v\n", string(b))
}
