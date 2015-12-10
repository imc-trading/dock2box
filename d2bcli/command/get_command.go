package command

import (
	"github.com/codegangsta/cli"
)

// NewGetCommand create new resource.
func NewGetCommand() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "get resource",
		Subcommands: []cli.Command{
			NewGetHostCommand(),
		},
	}
}
