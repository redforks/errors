package errors_test

import (
	"math/rand"

	. "github.com/redforks/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Code", func() {
	DescribeTable("Round", func(cause CausedBy) {
		low := uint32(rand.Int31n(10000))
		v := NewCode(cause, low)
		Ω(v.Caused()).Should(Equal(cause))
		Ω(uint32(v) - uint32(v.Caused())).Should(Equal(uint32(low)))
	},
		Entry("ByBug", ByBug),
		Entry("ByRuntime", ByRuntime),
		Entry("ByExternal", ByExternal),
		Entry("ByInput", ByInput),
		Entry("ByClientBug", ByClientBug),
	)
})
