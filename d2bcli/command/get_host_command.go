package command

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/mickep76/dock2box/client"
)

func NewGetHostCommand() cli.Command {
	return cli.Command{
		Name:  "host",
		Usage: "Get host",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			getHostCommandFunc(c)
		},
	}
}

func getHostCommandFunc(c *cli.Context) {
	var hostname string
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a hostname")
	} else {
		hostname = c.Args()[0]
	}

	clnt := client.New(c.GlobalString("server"))
	h, err := clnt.Host.Get(hostname)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%v\n", string(h.JSON()))
}
