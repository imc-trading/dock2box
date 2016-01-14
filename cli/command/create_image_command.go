package command

import (
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/cli/prompt"
	"github.com/imc-trading/dock2box/client"
)

func NewCreateImageCommand() cli.Command {
	return cli.Command{
		Name:  "image",
		Usage: "Create image",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			createImageCommandFunc(c)
		},
	}
}

func createImageCommandFunc(c *cli.Context) {
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

	s := client.Image{
		Image:     image,
		Type:      prompt.String("Type", prompt.Prompt{Default: "docker", FuncPtr: prompt.Enum, FuncInp: "file,docker"}),
		BootTagID: *chooseTag(clnt, ""),
	}

	// Create image
	clnt.Image.Create(&s)
}
