package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetBootImageCommand() cli.Command {
	return cli.Command{
		Name:  "boot-image",
		Usage: "Get boot image",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "all, a", Usage: "Get all boot images"},
		},
		Action: func(c *cli.Context) {
			getBootImageCommandFunc(c)
		},
	}
}

func getBootImageCommandFunc(c *cli.Context) {
	var bootImage string
	if !c.Bool("all") {
		if len(c.Args()) == 0 {
			log.Fatal("You need to specify a boot image")
		} else {
			bootImage = c.Args()[0]
		}
	}

	clnt := client.New(c.GlobalString("server"))

	if c.Bool("all") {
		i, err := clnt.BootImage.All()
		if err != nil {
			log.Fatal(err.Error())
		}
		b, _ := json.MarshalIndent(i, "", "  ")
		fmt.Printf("%v\n", string(b))
	} else {
		i, err := clnt.BootImage.Get(bootImage)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%v\n", string(i.JSON()))
	}
}
