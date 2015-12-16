package command

import (
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
	"github.com/imc-trading/dock2box/d2bcli/prompt"
)

func NewCreateImageCommand() cli.Command {
	return cli.Command{
		Name:  "image",
		Usage: "Create image",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "prompt, p", Usage: "Prompt for resource input"},
			cli.StringFlag{Name: "type, t", Usage: "Type (file, docker)"},
			cli.StringFlag{Name: "boot-image, b", Usage: "Boot image"},
		},
		Action: func(c *cli.Context) {
			createImageCommandFunc(c)
		},
	}
}

func chooseBootImage(clnt *client.Client) *string {
	r, err := clnt.BootImage.All()
	if err != nil {
		log.Fatalf(err.Error())
	}

	images := *r
	var list []string
	for _, v := range images {
		list = append(list, v.Image)
	}
	return &images[prompt.Choice("Choose image", -1, list)].ID
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

	if c.Bool("prompt") {
		s := client.Image{
			Image:       image,
			Type:        prompt.String("Type", prompt.Prompt{Default: "docker", FuncPtr: prompt.Enum, FuncInp: "file,docker"}),
			BootImageID: *chooseBootImage(clnt),
		}

		// Create image
		clnt.Image.Create(&s)
		return
	}
}
