package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetImageCommand() cli.Command {
	return cli.Command{
		Name:  "image",
		Usage: "Get image",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "all, a", Usage: "Get all images"},
		},
		Action: func(c *cli.Context) {
			getImageCommandFunc(c)
		},
	}
}

func getImageCommandFunc(c *cli.Context) {
	var image string
	if !c.Bool("all") {
		if len(c.Args()) == 0 {
			log.Fatal("You need to specify a image")
		} else {
			image = c.Args()[0]
		}
	}

	clnt := client.New(c.GlobalString("server"))

	if c.Bool("all") {
		i, err := clnt.Image.All()
		if err != nil {
			log.Fatal(err.Error())
		}
		b, _ := json.MarshalIndent(i, "", "  ")
		fmt.Printf("%v\n", string(b))
	} else {
		i, err := clnt.Image.Get(image)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%v\n", string(i.JSON()))
	}
}
