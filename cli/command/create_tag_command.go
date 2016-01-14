package command

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/jehiah/go-strftime"

	"github.com/imc-trading/dock2box/cli/prompt"
	"github.com/imc-trading/dock2box/client"
)

func NewCreateTagCommand() cli.Command {
	return cli.Command{
		Name:  "tag",
		Usage: "Create tag",
		Flags: []cli.Flag{},
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

	s := client.Tag{
		Tag:     tag,
		Created: prompt.String("Created", prompt.Prompt{Default: strftime.Format("%Y-%m-%dT%H:%M:%SZ", time.Now()), FuncPtr: prompt.Regex, FuncInp: ""}),
		SHA256:  prompt.String("SHA256", prompt.Prompt{NoDefault: true, FuncPtr: prompt.Regex, FuncInp: "^[0-9a-f]+$"}),
		ImageID: *chooseImage(clnt, ""),
	}

	// Is this correct?
	fmt.Println(string(s.JSON()))
	if !prompt.Bool("Is this correct", true) {
		os.Exit(1)
	}

	// Create image
	clnt.Tag.Create(&s)
}
