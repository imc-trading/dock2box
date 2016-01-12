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
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			getHostCommandFunc(c)
		},
	}
}

func getHostCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a hostname")
	}
	hostname := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	h, err := clnt.Host.Query(map[string]string{"host": hostname})
	if err != nil {
		log.Fatal(err.Error())
	}

	b, _ := json.MarshalIndent(h, "", "  ")
	fmt.Printf("%v\n", string(b))
}
