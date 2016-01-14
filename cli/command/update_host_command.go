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

func NewUpdateHostCommand() cli.Command {
	return cli.Command{
		Name:  "host",
		Usage: "Update host",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			updateHostCommandFunc(c)
		},
	}
}

func updateHostCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a hostname")
	}
	hostname := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	v, err := clnt.Host.Get(hostname)
	if err != nil {
		log.Fatal(err.Error())
	}

	h := client.Host{
		ID:    v.ID,
		Host:  prompt.String("Host", prompt.Prompt{Default: v.Host, FuncPtr: prompt.Regex, FuncInp: ""}),
		Build: prompt.Bool("Build", v.Build),
		Debug: prompt.Bool("Debug", v.Debug),
		GPT:   prompt.Bool("GPT", v.GPT),
		TagID: *chooseTag(clnt, v.TagID),
	}

	// Get labels
	labels := prompt.String("Comma-separated list of labels", prompt.Prompt{Default: strings.Join(v.Labels, ","), FuncPtr: prompt.Regex, FuncInp: "^([a-zA-Z][a-zA-Z0-9-]+,)*([a-zA-Z][a-zA-Z0-9-]+)$"})
	if labels == "" {
		h.Labels = []string{}
	} else {
		h.Labels = strings.Split(labels, ",")
	}

	h.KOpts = prompt.String("KOpts", prompt.Prompt{Default: v.KOpts, FuncPtr: prompt.Regex, FuncInp: "^(|[a-zA-Z0-9- ])+$"})
	h.TenantID = *chooseTenants(clnt, v.TenantID)
	h.SiteID = *chooseSite(clnt, v.SiteID)

	// Is this correct?
	fmt.Println(string(h.JSON()))
	if !prompt.Bool("Is this correct", true) {
		os.Exit(1)
	}

	// Update host
	clnt.Host.Update(hostname, &h)
}
