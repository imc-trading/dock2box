package command

import (
	"github.com/codegangsta/cli"
)

// NewQueryCommand query resources.
func NewQueryCommand() cli.Command {
	return cli.Command{
		Name:  "query",
		Usage: "Query for resources",
		Subcommands: []cli.Command{
			NewQueryHostCommand(),
			NewQueryInterfaceCommand(),
			NewQueryImageCommand(),
			NewQueryTagCommand(),
			NewQuerySiteCommand(),
			NewQuerySubnetCommand(),
			NewQueryTenantCommand(),
		},
	}
}
