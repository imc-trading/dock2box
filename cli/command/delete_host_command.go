package command

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/cli/prompt"
	"github.com/imc-trading/dock2box/client"
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

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	if !prompt.Bool("Are you sure you wan't to remove "+hostname, true) {
		os.Exit(1)
	}

	h, err := clnt.Host.Delete(hostname)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", string(h.JSON()))
}
