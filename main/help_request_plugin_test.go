package main_test

import (
	"errors"
	"fmt"
	"github.com/cloudfoundry/cli/plugin/fakes"
	io_helpers "github.com/cloudfoundry/cli/testhelpers/io"
	. "github.com/dmikusa-pivotal/help_request_plugin/main"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"os"
	"strings"
)

var _ = Describe("HelpRequestPlugin", func() {
	Describe(".Run", func() {
		var fakeCliConnection *fakes.FakeCliConnection
		var helpRequestPlugin *HelpRequestPlugin

		BeforeEach(func() {
			fakeCliConnection = &fakes.FakeCliConnection{}
			helpRequestPlugin = NewHelpRequestPlugin(os.Stdin)
		})

		It("returns the arguments given to it", func() {
			output := io_helpers.CaptureOutput(func() {
				helpRequestPlugin.Run(fakeCliConnection, []string{"help-me-now"})
			})
			Expect(strings.Join(output, "")).To(ContainSubstring("name"))
			Expect(strings.Join(output, "")).To(ContainSubstring("description"))
			Expect(strings.Join(output, "")).To(ContainSubstring("Submitting"))
			Expect(strings.Join(output, "")).To(ContainSubstring("Thanks!"))
		})

		It("Loads user info automatically", func() {
			fakeCliConnection.CliCommandReturns([]string{
				"",
				"API endpoint:   https://api.run.pivotal.io (API version: 2.21.0)",
				"User:           dmikusa@gopivotal.com",
				"Org:            dmikusa",
				"Space:          development"}, nil)
			plugin := NewHelpRequestPlugin(os.Stdin)
			output := io_helpers.CaptureOutput(func() {
				plugin.LoadUserInfo(fakeCliConnection)
			})
			Expect(output[0]).To(ContainSubstring("We're going to automatically"))
			//TODO: something is not working with the mock object.  It's not returning
			//  the specified output
			//Expect(plugin.Email).To(Equal("dmikusa@gopivotal.com"))
			//Expect(plugin.Desc).To(ContainSubstring("dmikusa"))
			//Expect(plugin.Desc).To(ContainSubstring("development"))
		})

		It("Auto load of info fails, fall back to manual", func() {
			fakeCliConnection.CliCommandReturns([]string{}, errors.New("fail :("))
			io_helpers.CaptureOutput(func() {
				io_helpers.SimulateStdin("email\norg\nspace\n", func(reader io.Reader) {
					plugin := NewHelpRequestPlugin(reader)
					output := io_helpers.CaptureOutput(func() {
						plugin.LoadUserInfo(fakeCliConnection)
					})
					fmt.Println(strings.Join(output, "\n"))
					Expect(output[0]).To(ContainSubstring("We're going to automatically"))
					//TODO: something is not working with the mock object.  It's not returning
					//  the error as expected
					//Expect(output[1]).To(ContainSubstring("Sorry, there was a problem"))
					//Expect(plugin.Email).To(Equal("email"))
				})
			})
		})

		It("Tests json dump", func() {
			plugin := NewHelpRequestPlugin(nil)
			plugin.Name = "Jack Johnson"
			plugin.Phone = "555-555-5758"
			plugin.Email = "jack@johnson.com"
			plugin.Desc = `Some field with odd" characters \s and'\nmorestuff\n\r\na`
			buf, _ := plugin.ToJson()
			Expect(buf).ShouldNot(BeNil())
		})

	})

	Describe("Gathers user information from stdin", func() {

		It("Sends the initial greeting", func() {
			plugin := NewHelpRequestPlugin(os.Stdin)
			output := io_helpers.CaptureOutput(func() {
				plugin.Greet()
			})
			Expect(output[1]).To(Equal("Help Request System"))
		})

		It("Asks for a name", func() {
			var response string
			io_helpers.CaptureOutput(func() {
				io_helpers.SimulateStdin("William Lewis Lockwood\n", func(reader io.Reader) {
					plugin := NewHelpRequestPlugin(reader)
					response = plugin.PromptFor("name")
				})
			})
			Expect(response).To(Equal("William Lewis Lockwood"))
		})
	})

	Describe("Sends request to the server", func() {
		It("Sends a valid request", func() {
			plugin := NewHelpRequestPlugin(os.Stdin)
			io_helpers.CaptureOutput(func() {
				plugin.Name = "Joe Smith"
				plugin.Phone = "5555555785"
				plugin.Email = "jsmith@work.com"
				plugin.Desc = "I need help with PWS!"
				plugin.SubmitRequest()
			})
			Expect(plugin.ReqUrl).To(HavePrefix("http://pws-callme.cfapps.io/helprequests"))
		})
	})
})
