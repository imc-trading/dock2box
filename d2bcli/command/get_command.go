package command

import (
	"github.com/codegangsta/cli"
)

// NewGetCommand create new resource.
func NewGetCommand() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "Get resource",
		Subcommands: []cli.Command{
			NewGetHostCommand(),
			NewGetHostsCommand(),
			NewGetImageCommand(),
			NewGetImagesCommand(),
			NewGetTagCommand(),
			NewGetTagsCommand(),
			NewGetSubnetCommand(),
			NewGetSubnetsCommand(),
			NewGetSiteCommand(),
			NewGetSitesCommand(),
			NewGetTenantCommand(),
			NewGetTenantsCommand(),
		},
	}
}
