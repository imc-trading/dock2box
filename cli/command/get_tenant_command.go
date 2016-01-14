package command

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetTenantCommand() cli.Command {
	return cli.Command{
		Name:  "tenant",
		Usage: "Get tenant",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			getTenantCommandFunc(c)
		},
	}
}

func getTenantCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a tenantname")
	}
	tenant := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	t, err := clnt.Tenant.Get(tenant)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", string(t.JSON()))
}
