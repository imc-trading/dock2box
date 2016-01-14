package command

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/cli/prompt"
	"github.com/imc-trading/dock2box/client"
)

func NewDeleteTenantCommand() cli.Command {
	return cli.Command{
		Name:  "tenant",
		Usage: "Delete tenant",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			deleteTenantCommandFunc(c)
		},
	}
}

func deleteTenantCommandFunc(c *cli.Context) {
	var tenant string
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a tenant")
	} else {
		tenant = c.Args()[0]
	}

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	if !prompt.Bool("Are you sure you wan't to remove "+tenant, true) {
		os.Exit(1)
	}

	t, err := clnt.Tenant.Delete(tenant)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", string(t.JSON()))
}
