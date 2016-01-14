package command

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/cli/prompt"
	"github.com/imc-trading/dock2box/client"
)

func NewCreateHostCommand() cli.Command {
	return cli.Command{
		Name:  "host",
		Usage: "Create host",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			createHostCommandFunc(c)
		},
	}
}

func createHostCommandFunc(c *cli.Context) {
	var hostname string
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a hostname")
	} else {
		hostname = c.Args()[0]
	}

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	h := client.Host{
		Host:  hostname,
		Build: prompt.Bool("Build", true),
		Debug: prompt.Bool("Debug", false),
		GPT:   prompt.Bool("GPT", false),
		TagID: *chooseTag(clnt, ""),
	}

	// Get labels
	labels := prompt.String("Comma-separated list of labels", prompt.Prompt{Default: "", FuncPtr: prompt.Regex, FuncInp: "^([a-zA-Z][a-zA-Z0-9-]+,)*([a-zA-Z][a-zA-Z0-9-]+)$"})
	if labels == "" {
		h.Labels = []string{}
	} else {
		h.Labels = strings.Split(labels, ",")
	}

	h.KOpts = prompt.String("KOpts", prompt.Prompt{Default: "", FuncPtr: prompt.Regex, FuncInp: "^(|[a-zA-Z0-9- ])+$"})
	h.TenantID = *chooseTenants(clnt, "")
	h.SiteID = *chooseSite(clnt, "")

	// Is this correct?
	fmt.Println(string(h.JSON()))
	if !prompt.Bool("Is this correct", true) {
		os.Exit(1)
	}

	// Create host
	clnt.Host.Create(&h)
}
