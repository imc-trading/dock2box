package command

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewDeleteTagCommand() cli.Command {
	return cli.Command{
		Name:  "tag",
		Usage: "Delete tag",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			deleteTagCommandFunc(c)
		},
	}
}

func deleteTagCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a tag")
	}
	tag := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	h, err := clnt.Image.Delete(tag)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", string(h.JSON()))
}
