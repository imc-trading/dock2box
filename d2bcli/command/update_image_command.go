package command

import (
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
	"github.com/imc-trading/dock2box/d2bcli/prompt"
)

func NewUpdateImageCommand() cli.Command {
	return cli.Command{
		Name:  "image",
		Usage: "Update image",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "prompt, p", Usage: "Prompt for resource input"},
			cli.StringFlag{Name: "type, t", Usage: "Type (file, docker)"},
			cli.StringFlag{Name: "boot-image, b", Usage: "Boot image"},
		},
		Action: func(c *cli.Context) {
			updateImageCommandFunc(c)
		},
	}
}

func updateImageCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a image")
	}
	image := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	v, err := clnt.Image.Get(image)
	if err != nil {
		log.Fatal(err.Error())
	}

	if c.Bool("prompt") {
		s := client.Image{
			ID:          v.ID,
			Image:       image,
			Type:        prompt.String("Type", prompt.Prompt{Default: v.Type, FuncPtr: prompt.Enum, FuncInp: "file,docker"}),
			BootImageID: *chooseBootImage(clnt, v.BootImageID),
		}

		// Update image
		clnt.Image.Update(image, &s)
		return
	}
}
