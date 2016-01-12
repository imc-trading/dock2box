package command

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetImageCommand() cli.Command {
	return cli.Command{
		Name:  "image",
		Usage: "Get image",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			getImageCommandFunc(c)
		},
	}
}

func getImageCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a image")
	}
	image := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))

	i, err := clnt.Image.Get(image)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", string(i.JSON()))
}
