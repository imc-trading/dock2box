package command

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetSiteCommand() cli.Command {
	return cli.Command{
		Name:  "site",
		Usage: "Get site",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			getSiteCommandFunc(c)
		},
	}
}

func getSiteCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a site")
	}
	site := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))

	s, err := clnt.Site.Get(site)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", string(s.JSON()))
}
