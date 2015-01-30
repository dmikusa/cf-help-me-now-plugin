package main_test

import (
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
			Expect(strings.Join(output, "")).To(ContainSubstring("email"))
			Expect(strings.Join(output, "")).To(ContainSubstring("Submitting"))
			Expect(strings.Join(output, "")).To(ContainSubstring("Thanks!"))
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

		It("Asks for a type, get PHONE", func() {
			var response string
			io_helpers.CaptureOutput(func() {
				io_helpers.SimulateStdin("phone\n", func(reader io.Reader) {
					plugin := NewHelpRequestPlugin(reader)
					response = plugin.PromptForRequestType()
				})
			})
			Expect(response).To(Equal("PHONE"))
		})

		It("Asks for a type, get IM", func() {
			var response string
			io_helpers.CaptureOutput(func() {
				io_helpers.SimulateStdin("IM\n", func(reader io.Reader) {
					plugin := NewHelpRequestPlugin(reader)
					response = plugin.PromptForRequestType()
				})
			})
			Expect(response).To(Equal("IM"))
		})

		It("Asks for a type, get DEFAULT", func() {
			var response string
			io_helpers.CaptureOutput(func() {
				io_helpers.SimulateStdin("dkjfdk\n", func(reader io.Reader) {
					plugin := NewHelpRequestPlugin(reader)
					response = plugin.PromptForRequestType()
				})
			})
			Expect(response).To(Equal("IM"))
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
				plugin.ReqType = "PHONE"
				plugin.SubmitRequest()
			})
			Expect(plugin.ReqUrl).To(HavePrefix("http://localhost:8080/helprequests/"))
		})
	})
})
