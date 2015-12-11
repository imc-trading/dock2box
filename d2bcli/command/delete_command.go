package command

import (
	"github.com/codegangsta/cli"
)

// NewDeleteCommand delete resource.
func NewDeleteCommand() cli.Command {
	return cli.Command{
		Name:  "delete",
		Usage: "Delete resource",
		Subcommands: []cli.Command{
			NewDeleteHostCommand(),
			NewDeleteImageCommand(),
			NewDeleteSubnetCommand(),
			NewDeleteSiteCommand(),
			NewDeleteBootImageCommand(),
			NewDeleteTenantCommand(),
		},
	}
}
