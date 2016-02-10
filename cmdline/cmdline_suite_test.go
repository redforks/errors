package cmdline

import (
	"github.com/onsi/ginkgo"

	"testing"
)

var t = ginkgo.GinkgoT

func TestCmdline(t *testing.T) {
	ginkgo.RunSpecs(t, "Cmdline Suite")
}
