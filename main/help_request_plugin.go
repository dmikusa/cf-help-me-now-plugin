package main

import (
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"io"
	"os"
)

type HelpRequestPlugin struct {
	ui    terminal.UI
	name  string
	phone string
	email string
}

func NewHelpRequestPlugin(reader io.Reader) *HelpRequestPlugin {
	return &HelpRequestPlugin{
		ui: terminal.NewUI(reader, terminal.NewTeePrinter()),
	}
}

func (p *HelpRequestPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "help-me-now" {
		p.Greet()
		p.name = p.PromptFor("name")
		p.phone = p.PromptFor("phone")
		p.email = p.PromptFor("email")
		p.SubmitRequest()
		p.Finish()
	}
}

func (p *HelpRequestPlugin) Greet() {
	p.ui.Say("")
	p.ui.Say("Help Request System")
	p.ui.Say("")
	p.ui.Say("Welcome, we need to gather a bit of information to get started.")
}

func (p *HelpRequestPlugin) PromptFor(item string) string {
	return p.ui.Ask("What's your %s?", item)
}

func (p *HelpRequestPlugin) SubmitRequest() {
	p.ui.Say("Submitting request...")
	p.ui.Say("...<actually submit request here [name=%s, phone=%s, email=%s]", p.name, p.phone, p.email)
	p.ui.Say("done!")
}

func (p *HelpRequestPlugin) Finish() {
	p.ui.Say("Thanks!  I've submitted your support request.")
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
	plugin.Start(NewHelpRequestPlugin(os.Stdin))
}
