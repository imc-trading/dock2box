package main

import (
	"os"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/d2bcli/command"
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "d2b-cli"
	app.Version = Version
	app.Usage = "Command line tool for interacting with Dock2Box resources."
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "server, s", Value: "http://localhost:8080/v1", EnvVar: "D2B_SERVER", Usage: "URL for Dock2Box API"},
		cli.BoolFlag{Name: "debug, d", Usage: "Debug"},
	}
	app.Commands = []cli.Command{
		command.NewCreateCommand(),
		command.NewUpdateCommand(),
		command.NewDeleteCommand(),
		command.NewGetCommand(),
	}

	app.Run(os.Args)
}
