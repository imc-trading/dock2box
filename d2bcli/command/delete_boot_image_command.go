package command

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewDeleteBootImageCommand() cli.Command {
	return cli.Command{
		Name:  "boot-image",
		Usage: "Delete boot image",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			deleteBootImageCommandFunc(c)
		},
	}
}

func deleteBootImageCommandFunc(c *cli.Context) {
	var bootImage string
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a boot image")
	} else {
		bootImage = c.Args()[0]
	}

	clnt := client.New(c.GlobalString("server"))

	h, err := clnt.BootImage.Delete(bootImage)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", string(h.JSON()))
}
