package command

import (
	"github.com/codegangsta/cli"
)

// NewUpdateCommand create new resource.
func NewUpdateCommand() cli.Command {
	return cli.Command{
		Name:  "update",
		Usage: "Update resource",
		Subcommands: []cli.Command{
			NewUpdateSubnetCommand(),
			NewUpdateSiteCommand(),
			NewUpdateImageCommand(),
			NewUpdateTenantCommand(),
			NewUpdateHostCommand(),
		},
	}
}
