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
	})

})
