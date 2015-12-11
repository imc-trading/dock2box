package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetHostCommand() cli.Command {
	return cli.Command{
		Name:  "host",
		Usage: "Get host",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "all, a", Usage: "Get all hosts"},
		},
		Action: func(c *cli.Context) {
			getHostCommandFunc(c)
		},
	}
}

func getHostCommandFunc(c *cli.Context) {
	var hostname string
	if !c.Bool("all") {
		if len(c.Args()) == 0 {
			log.Fatal("You need to specify a hostname")
		} else {
			hostname = c.Args()[0]
		}
	}

	clnt := client.New(c.GlobalString("server"))

	if c.Bool("all") {
		h, err := clnt.Host.All()
		if err != nil {
			log.Fatal(err.Error())
		}
		b, _ := json.MarshalIndent(h, "", "  ")
		fmt.Printf("%v\n", string(b))
	} else {
		h, err := clnt.Host.Get(hostname)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%v\n", string(h.JSON()))
	}
}
