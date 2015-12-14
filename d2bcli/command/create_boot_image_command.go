package command

import (
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
	"github.com/imc-trading/dock2box/d2bcli/prompt"
)

func NewCreateBootImageCommand() cli.Command {
	return cli.Command{
		Name:  "boot-image",
		Usage: "Create boot boot image",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "prompt, p", Usage: "Prompt for resource input"},
			cli.StringFlag{Name: "kopts, k", Usage: "Boot bootImage"},
		},
		Action: func(c *cli.Context) {
			createBootImageCommandFunc(c)
		},
	}
}

func createBootImageCommandFunc(c *cli.Context) {
	var bootImage string
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a bootImage")
	} else {
		bootImage = c.Args()[0]
	}

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	if c.Bool("prompt") {
		s := client.BootImage{
			Image: bootImage,
			KOpts: prompt.String("Kernel Options", prompt.Prompt{Default: "", FuncPtr: prompt.Regex, FuncInp: ""}),
		}

		// Create bootImage
		clnt.BootImage.Create(&s)
		return
	}
}
