package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetSiteCommand() cli.Command {
	return cli.Command{
		Name:  "site",
		Usage: "Get site",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "all, a", Usage: "Get all sites"},
		},
		Action: func(c *cli.Context) {
			getSiteCommandFunc(c)
		},
	}
}

func getSiteCommandFunc(c *cli.Context) {
	var site string
	if !c.Bool("all") {
		if len(c.Args()) == 0 {
			log.Fatal("You need to specify a site")
		} else {
			site = c.Args()[0]
		}
	}

	clnt := client.New(c.GlobalString("server"))

	if c.Bool("all") {
		s, err := clnt.Site.All()
		if err != nil {
			log.Fatal(err.Error())
		}
		b, _ := json.MarshalIndent(s, "", "  ")
		fmt.Printf("%v\n", string(b))
	} else {
		s, err := clnt.Site.Get(site)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%v\n", string(s.JSON()))
	}
}
