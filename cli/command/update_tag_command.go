package command

import (
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
	"github.com/imc-trading/dock2box/cli/prompt"
)

func NewUpdateTagCommand() cli.Command {
	return cli.Command{
		Name:  "tag",
		Usage: "Update tag",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "prompt, p", Usage: "Prompt for resource input"},
			cli.StringFlag{Name: "created, c", Usage: "Created"},
			cli.StringFlag{Name: "sha256, s", Usage: "SHA256"},
			cli.StringFlag{Name: "image, i", Usage: "Image"},
		},
		Action: func(c *cli.Context) {
			updateTagCommandFunc(c)
		},
	}
}

func updateTagCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a tag")
	}
	tag := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	v, err := clnt.Tag.Get(tag)
	if err != nil {
		log.Fatal(err.Error())
	}

	if c.Bool("prompt") {
		s := client.Tag{
			ID:      v.ID,
			Tag:     prompt.String("Image", prompt.Prompt{Default: v.Tag, FuncPtr: prompt.Regex, FuncInp: ""}),
			Created: prompt.String("Created", prompt.Prompt{Default: v.Created, FuncPtr: prompt.Regex, FuncInp: ""}),
			SHA256:  prompt.String("SHA256", prompt.Prompt{Default: v.SHA256, FuncPtr: prompt.Regex, FuncInp: "^[0-9a-f]+$"}),
			ImageID: *chooseImage(clnt, v.ImageID),
		}

		// Update tag
		clnt.Tag.Update(tag, &s)
		return
	}
}
