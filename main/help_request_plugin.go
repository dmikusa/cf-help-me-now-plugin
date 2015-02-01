package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func panicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

type HelpRequestPlugin struct {
	ui     terminal.UI
	Name   string
	Phone  string
	Email  string
	Desc   string
	ReqUrl string
}

func NewHelpRequestPlugin(reader io.Reader) *HelpRequestPlugin {
	return &HelpRequestPlugin{
		ui: terminal.NewUI(reader, terminal.NewTeePrinter()),
	}
}

func (p *HelpRequestPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "help-me-now" {
		if len(args) == 1 {
			p.Clear()
			p.Greet()
			p.Name = p.PromptFor("name")
			p.Phone = p.PromptFor("phone")
			p.Desc = p.PromptFor("problem description")
			p.LoadUserInfo(cliConnection)
			p.SubmitRequest()
			p.Finish()
		} else if len(args) == 2 && args[1] == "--status" {
			path := filepath.Join(p.pluginDataDir(), "request.txt")
			if _, err := os.Stat(path); os.IsExist(err) {
				p.Load()
				p.ui.Say("Request Status: %s", p.Status())
			} else {
				p.ui.Say("Sorry, could not find an existing request.  Please submit again.")
			}
		}
	}
}

func (p *HelpRequestPlugin) appendToDesc(val string) {
	p.Desc += "\n" + val + "\n"
}

func (p *HelpRequestPlugin) pluginDataDir() string {
	u, err := user.Current()
	panicOnError(err)
	return filepath.Join(u.HomeDir, ".cf", "help-me-now")
}

func (p *HelpRequestPlugin) LoadUserInfo(cliConnection plugin.CliConnection) {
	p.ui.Say("We're going to automatically gather your account information.  One minute...")
	output, err := cliConnection.CliCommandWithoutTerminalOutput("target")
	if err != nil {
		p.ui.Say("Sorry, there was a problem loading your information.  We'll need" +
			" to collect a few things to continue.")
		p.Email = p.PromptFor("email")
		p.appendToDesc(p.PromptFor("organization name"))
		p.appendToDesc(p.PromptFor("space name"))
	} else {
		for _, line := range output {
			sections := strings.SplitN(line, ":", 2)
			if len(sections) == 2 {
				key := strings.TrimSpace(sections[0])
				val := strings.TrimSpace(sections[1])
				if key == "User" {
					p.Email = val
				} else {
					p.appendToDesc(line)
				}
			}
		}
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
	var err error
	p.ui.Say("Submitting request...")
	p.Save()
	p.ReqUrl, err = p.send()
	if err != nil {
		p.ui.Say("Sorry, we're unable to submit your request at this time.")
		p.ui.Say("Error: %q", err)
	} else {
		p.ui.Say("Your request was accepted!  You should receive a call shortly.")
		p.ui.Say("To check the status of your request, run `cf help-me-now --status`.")
	}
	p.Save()
}

func (p *HelpRequestPlugin) Clear() {
	path := filepath.Join(p.pluginDataDir(), "request.txt")
	if _, err := os.Stat(path); os.IsExist(err) {
		err = os.Remove(path)
		panicOnError(err)
	}
}

func (p *HelpRequestPlugin) Load() {
	path := filepath.Join(p.pluginDataDir(), "request.txt")
	data, err := ioutil.ReadFile(path)
	panicOnError(err)
	err = json.Unmarshal(data, &p)
	panicOnError(err)
}

func (p *HelpRequestPlugin) Save() {
	path := p.pluginDataDir()
	os.MkdirAll(path, 0755)
	path = filepath.Join(path, "request.txt")
	b, err := json.Marshal(p)
	panicOnError(err)
	err = ioutil.WriteFile(path, b, 0644)
	panicOnError(err)
}

func (p *HelpRequestPlugin) FromJson(blob []byte) error {
	m := make(map[string]string)
	err := json.Unmarshal(blob, &m)
	if err != nil {
		return err
	}
	p.Name = m["fullname"]
	p.Phone = m["phone"]
	p.Email = m["email"]
	p.Desc = m["description"]
	return nil
}

func (p *HelpRequestPlugin) ToJson() ([]byte, error) {
	m := make(map[string]string)
	m["fullname"] = p.Name
	m["phone"] = p.Phone
	m["email"] = p.Email
	m["username"] = p.Email
	m["description"] = p.Desc
	m["type"] = "PHONE"
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (p *HelpRequestPlugin) send() (string, error) {
	url := "http://pws-callme.cfapps.io/helprequests"
	b, err := p.ToJson()
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
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

func (p *HelpRequestPlugin) Status() (interface{}, error) {
	resp, err := http.Get(p.ReqUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		return "", err
	}
	return m["status"], nil
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
					Usage: "help-me-now [--status]\n    cf help-me-now",
				},
			},
		},
	}
}

func main() {
	plugin.Start(NewHelpRequestPlugin(os.Stdin))
}
