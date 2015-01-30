package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"io"
	"net/http"
	"os"
	"strings"
)

type HelpRequestPlugin struct {
	ui      terminal.UI
	Name    string
	Phone   string
	Email   string
	Desc    string
	ReqType string
	ReqUrl  string
}

func NewHelpRequestPlugin(reader io.Reader) *HelpRequestPlugin {
	return &HelpRequestPlugin{
		ui: terminal.NewUI(reader, terminal.NewTeePrinter()),
	}
}

func (p *HelpRequestPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "help-me-now" {
		p.Greet()
		p.Name = p.PromptFor("name")
		p.Phone = p.PromptFor("phone")
		p.Email = p.PromptFor("email")
		p.Desc = p.PromptFor("problem description")
		p.ReqType = p.PromptForRequestType()
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

func (p *HelpRequestPlugin) PromptForRequestType() string {
	resp := p.ui.Ask("How would you like to be contacted? (phone or IM)")
	if strings.ToLower(resp) == "phone" {
		return "PHONE"
	} else if strings.ToUpper(resp) == "IM" {
		return "IM"
	} else {
		p.ui.Say("Invalid response, defaulting to IM.")
		return "IM"
	}
}

func (p *HelpRequestPlugin) SubmitRequest() {
	var err error
	p.ui.Say("Submitting request...")
	p.ReqUrl, err = p.send()
	if err != nil {
		p.ui.Say("Failed: <%v>", err)
	}
	p.ui.Say("done!")
}

func (p *HelpRequestPlugin) send() (string, error) {
	//url := "https://pws-callme.cfapps.io/helprequests"
	url := "http://localhost:8080/helprequests"
	req, err := http.NewRequest("POST", url,
		bytes.NewBufferString(
			`{"fullname":"`+p.Name+
				`","phone":"`+p.Phone+
				`","email":"`+p.Email+
				`","username":"`+p.Email+
				`","description":"`+p.Desc+
				`","type":"`+p.ReqType+`"}`))
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Security-Token", "NFJKJDIJ#RIOJFLSNDFIOEWJF)#U$8238r3")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	} else {
		if resp.StatusCode == 201 {
			return strings.Join(resp.Header["Location"], ","), nil
		} else if resp.StatusCode == 403 {
			return "", errors.New("Forbidden, invalid security token.")
		} else {
			return "", errors.New(fmt.Sprintf("Status: <%s> - <%d>", resp.Status, resp.StatusCode))
		}
	}
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
