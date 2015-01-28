package main

import (
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"os"
)

type HelpRequestPlugin struct {
	ui terminal.UI
}

func NewHelpRequestPlugin() *HelpRequestPlugin {
	return &HelpRequestPlugin{
		ui: terminal.NewUI(os.Stdin, terminal.NewTeePrinter()),
	}
}

func (p *HelpRequestPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "help-me-now" {
		p.ui.Say("Args: %v", args)
	}
}

func (p *HelpRequestPlugin) Greet() {
	p.ui.Say("Help Request System")
	p.ui.Say("")
	p.ui.Say("Welcome, we need to gather a bit of information to get started.")
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
	plugin.Start(NewHelpRequestPlugin())
}
