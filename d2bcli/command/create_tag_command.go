package command

import (
	"log"
	"time"

	"github.com/codegangsta/cli"
	"github.com/jehiah/go-strftime"

	"github.com/imc-trading/dock2box/client"
	"github.com/imc-trading/dock2box/d2bcli/prompt"
)

func NewCreateTagCommand() cli.Command {
	return cli.Command{
		Name:  "tag",
		Usage: "Create tag",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "prompt, p", Usage: "Prompt for resource input"},
			cli.StringFlag{Name: "created, c", Usage: "Created"},
			cli.StringFlag{Name: "sha256, s", Usage: "SHA256"},
			cli.StringFlag{Name: "image, i", Usage: "Image"},
		},
		Action: func(c *cli.Context) {
			createTagCommandFunc(c)
		},
	}
}

func createTagCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a tag")
	}
	tag := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	if c.Bool("prompt") {
		s := client.Tag{
			Tag:     tag,
			Created: prompt.String("Created", prompt.Prompt{Default: strftime.Format("%Y-%m-%dT%H:%M:%SZ", time.Now()), FuncPtr: prompt.Regex, FuncInp: ""}),
			SHA256:  prompt.String("SHA256", prompt.Prompt{NoDefault: true, FuncPtr: prompt.Regex, FuncInp: "^[0-9a-f]+$"}),
			ImageID: *chooseImage(clnt, ""),
		}

		// Create image
		clnt.Tag.Create(&s)
		return
	}
}
