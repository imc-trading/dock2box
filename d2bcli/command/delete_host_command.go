package command

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
)

func NewDeleteHostCommand() cli.Command {
	return cli.Command{
		Name:  "host",
		Usage: "Delete host",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			deleteHostCommandFunc(c)
		},
	}
}

func deleteHostCommandFunc(c *cli.Context) {
	var hostname string
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a hostname")
	} else {
		hostname = c.Args()[0]
	}
	fmt.Println("Delete Host: ", hostname)
}
