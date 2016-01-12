package command

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"

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

	t, err := clnt.Tenant.Delete(tenant)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", string(t.JSON()))
}
