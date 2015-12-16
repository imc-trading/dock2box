package command

import (
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
	"github.com/imc-trading/dock2box/d2bcli/prompt"
)

func NewUpdateTenantCommand() cli.Command {
	return cli.Command{
		Name:  "tenant",
		Usage: "Create tenant",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			updateTenantCommandFunc(c)
		},
	}
}

func updateTenantCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a tenant")
	}
	tenant := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	v, err := clnt.Tenant.Get(tenant)
	if err != nil {
		log.Fatal(err.Error())
	}

	s := client.Tenant{
		ID:     v.ID,
		Tenant: prompt.String("Tenant", prompt.Prompt{Default: v.Tenant, FuncPtr: prompt.Regex, FuncInp: ""}),
	}

	// Create tenant
	clnt.Tenant.Update(tenant, &s)
	return
}
