package main_test

import (
	"github.com/cloudfoundry/cli/plugin/fakes"
	io_helpers "github.com/cloudfoundry/cli/testhelpers/io"
	. "github.com/dmikusa-pivotal/help_request_plugin/main"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HelpRequestPlugin", func() {
	Describe(".Run", func() {
		var fakeCliConnection *fakes.FakeCliConnection
		var helpRequestPlugin *HelpRequestPlugin

		BeforeEach(func() {
			fakeCliConnection = &fakes.FakeCliConnection{}
			helpRequestPlugin = NewHelpRequestPlugin()
		})

		It("returns the arguments given to it", func() {
			output := io_helpers.CaptureOutput(func() {
				helpRequestPlugin.Run(fakeCliConnection, []string{"help-me-now", "arg1", "arg2"})
			})
			Expect(output[0]).To(Equal("Args: [help-me-now arg1 arg2]"))
		})
	})

	Describe("Gathers user information from stdin", func() {
		It("Sends the initial greeting", func() {
			plugin := NewHelpRequestPlugin()
			output := io_helpers.CaptureOutput(func() {
				plugin.Greet()
			})
			Expect(output[0]).To(Equal("Help Request System"))
		})
	})
})
