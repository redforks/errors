package errors_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"testing"
)

func TestErrors(t *testing.T) {
	gomega.RegisterFailHandler(Fail)
	RunSpecs(t, "Errors Suite")
}
