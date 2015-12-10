package command

import (
	"github.com/codegangsta/cli"
)

// NewDeleteCommand delete resource.
func NewDeleteCommand() cli.Command {
	return cli.Command{
		Name:  "delete",
		Usage: "delete resource",
		Subcommands: []cli.Command{
			NewDeleteHostCommand(),
		},
	}
}
