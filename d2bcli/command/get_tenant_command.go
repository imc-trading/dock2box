package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetTenantCommand() cli.Command {
	return cli.Command{
		Name:  "tenant",
		Usage: "Get tenant",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "all, a", Usage: "Get all tenants"},
		},
		Action: func(c *cli.Context) {
			getTenantCommandFunc(c)
		},
	}
}

func getTenantCommandFunc(c *cli.Context) {
	var tenant string
	if !c.Bool("all") {
		if len(c.Args()) == 0 {
			log.Fatal("You need to specify a tenantname")
		} else {
			tenant = c.Args()[0]
		}
	}

	clnt := client.New(c.GlobalString("server"))

	if c.Bool("all") {
		t, err := clnt.Tenant.All()
		if err != nil {
			log.Fatal(err.Error())
		}
		b, _ := json.MarshalIndent(t, "", "  ")
		fmt.Printf("%v\n", string(b))
	} else {
		t, err := clnt.Tenant.Get(tenant)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%v\n", string(t.JSON()))
	}
}
