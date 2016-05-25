package cmdline

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"testing"
)

func TestCmdline(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Cmdline Suite")
}
