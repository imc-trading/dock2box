package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetSitesCommand() cli.Command {
	return cli.Command{
		Name:  "sites",
		Usage: "Get all sites",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			getSitesCommandFunc(c)
		},
	}
}

func getSitesCommandFunc(c *cli.Context) {
	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	s, err := clnt.Site.All()
	if err != nil {
		log.Fatal(err.Error())
	}
	b, _ := json.MarshalIndent(s, "", "  ")
	fmt.Printf("%v\n", string(b))
}
