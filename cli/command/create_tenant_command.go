package command

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/cli/prompt"
	"github.com/imc-trading/dock2box/client"
)

func NewCreateTenantCommand() cli.Command {
	return cli.Command{
		Name:  "tenant",
		Usage: "Create tenant",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			createTenantCommandFunc(c)
		},
	}
}

func createTenantCommandFunc(c *cli.Context) {
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

	s := client.Tenant{
		Tenant: tenant,
	}

	// Is this correct?
	fmt.Println(string(s.JSON()))
	if !prompt.Bool("Is this correct", true) {
		os.Exit(1)
	}

	// Create tenant
	clnt.Tenant.Create(&s)
	return
}
