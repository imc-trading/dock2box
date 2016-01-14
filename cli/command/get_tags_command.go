package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetTagsCommand() cli.Command {
	return cli.Command{
		Name:  "tags",
		Usage: "Get all tag",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			getTagsCommandFunc(c)
		},
	}
}

func getTagsCommandFunc(c *cli.Context) {
	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	i, err := clnt.Tag.All()
	if err != nil {
		log.Fatal(err.Error())
	}
	b, _ := json.MarshalIndent(i, "", "  ")
	fmt.Printf("%v\n", string(b))
}
