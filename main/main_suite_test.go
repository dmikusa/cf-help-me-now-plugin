package main_test

import (
	"github.com/cloudfoundry/cli/testhelpers/plugin_builder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)

	plugin_builder.BuildTestBinary(".", "help_request_plugin")

	RunSpecs(t, "Main Suite")
}
