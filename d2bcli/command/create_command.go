package command

import (
	"github.com/codegangsta/cli"
)

// NewCreateCommand create new resource.
func NewCreateCommand() cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "create resource",
		Subcommands: []cli.Command{
			NewCreateHostCommand(),
		},
	}
}
