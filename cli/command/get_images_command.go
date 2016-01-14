package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetImagesCommand() cli.Command {
	return cli.Command{
		Name:  "images",
		Usage: "Get all images",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			getImagesCommandFunc(c)
		},
	}
}

func getImagesCommandFunc(c *cli.Context) {
	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	i, err := clnt.Image.All()
	if err != nil {
		log.Fatal(err.Error())
	}
	b, _ := json.MarshalIndent(i, "", "  ")
	fmt.Printf("%v\n", string(b))
}
