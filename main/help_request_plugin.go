package main

import (
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"os"
)

type HelpRequestPlugin struct {
}

func (p *HelpRequestPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	ui := terminal.NewUI(os.Stdin, terminal.NewTeePrinter())
	if args[0] == "help-me-now" {
		ui.Say("Args: %v", args)
	}
}

func (p *HelpRequestPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "HelpRequestPlugin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "help-me-now",
				HelpText: "Submit a request for help now!",
				UsageDetails: plugin.Usage{
					Usage: "help-me-now\n    cf help-me-now",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(HelpRequestPlugin))
}
