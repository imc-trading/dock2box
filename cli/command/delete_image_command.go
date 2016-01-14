package command

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewDeleteImageCommand() cli.Command {
	return cli.Command{
		Name:  "image",
		Usage: "Delete image",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			deleteImageCommandFunc(c)
		},
	}
}

func deleteImageCommandFunc(c *cli.Context) {
	var image string
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a image")
	} else {
		image = c.Args()[0]
	}

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	h, err := clnt.Image.Delete(image)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", string(h.JSON()))
}
