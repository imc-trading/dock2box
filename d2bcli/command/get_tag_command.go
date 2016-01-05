package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetTagCommand() cli.Command {
	return cli.Command{
		Name:  "tag",
		Usage: "Get tag",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "all, a", Usage: "Get all tags"},
		},
		Action: func(c *cli.Context) {
			getTagCommandFunc(c)
		},
	}
}

func getTagCommandFunc(c *cli.Context) {
	var tag string
	if !c.Bool("all") {
		if len(c.Args()) == 0 {
			log.Fatal("You need to specify a tag")
		}
		tag = c.Args()[0]
	}

	clnt := client.New(c.GlobalString("server"))

	if c.Bool("all") {
		i, err := clnt.Tag.All()
		if err != nil {
			log.Fatal(err.Error())
		}
		b, _ := json.MarshalIndent(i, "", "  ")
		fmt.Printf("%v\n", string(b))
	} else {
		i, err := clnt.Tag.Get(tag)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%v\n", string(i.JSON()))
	}
}
