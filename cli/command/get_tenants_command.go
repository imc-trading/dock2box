package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewGetTenantsCommand() cli.Command {
	return cli.Command{
		Name:  "tenants",
		Usage: "Get all tenant",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			getTenantsCommandFunc(c)
		},
	}
}

func getTenantsCommandFunc(c *cli.Context) {
	clnt := client.New(c.GlobalString("server"))

	t, err := clnt.Tenant.All()
	if err != nil {
		log.Fatal(err.Error())
	}
	b, _ := json.MarshalIndent(t, "", "  ")
	fmt.Printf("%v\n", string(b))
}
