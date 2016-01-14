package command

import (
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/cli/prompt"
	"github.com/imc-trading/dock2box/client"
)

func NewUpdateImageCommand() cli.Command {
	return cli.Command{
		Name:  "image",
		Usage: "Update image",
		Flags: []cli.Flag{},
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

	s := client.Image{
		ID:        v.ID,
		Image:     prompt.String("Image", prompt.Prompt{Default: v.Image, FuncPtr: prompt.Regex, FuncInp: ""}),
		Type:      prompt.String("Type", prompt.Prompt{Default: v.Type, FuncPtr: prompt.Enum, FuncInp: "file,docker"}),
		BootTagID: *chooseTag(clnt, v.BootTagID),
	}

	// Update image
	clnt.Image.Update(image, &s)
}
