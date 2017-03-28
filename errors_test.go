package errors_test

import (
	syserr "errors"

	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/onsi/gomega"

	"github.com/redforks/errors"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("errors", func() {

	assertError := func(e *errors.Error, msg string, causedBy errors.CausedBy) {
		Ω(e.Error()).Should(Equal(msg))
		Ω(e.CausedBy).Should(Equal(causedBy))
	}

	Context("Wrap error", func() {

		It("New", func() {
			assertError(errors.New("foo"), "foo", errors.ByBug)
		})

		It("NewBug", func() {
			assertError(errors.NewBug(syserr.New("foo")), "foo", errors.ByBug)
		})

		It("NewRuntime", func() {
			assertError(errors.NewRuntime(syserr.New("foo")), "foo", errors.ByRuntime)
		})

		It("NewExternal", func() {
			assertError(errors.NewExternal(syserr.New("foo")), "foo", errors.ByExternal)
		})

		It("NewInput", func() {
			assertError(errors.NewInput(syserr.New("foo")), "foo", errors.ByInput)
		})

		Context("Not wrap nil", func() {

			It("NewBug", func() {
				Ω(errors.NewBug(nil)).Should(BeNil())
			})

			It("NewRuntime", func() {
				Ω(errors.NewRuntime(nil)).Should(BeNil())
			})

			It("NewInput", func() {
				Ω(errors.NewInput(nil)).Should(BeNil())
			})

			It("NewExternal", func() {
				Ω(errors.NewExternal(nil)).Should(BeNil())
			})

		})

		DescribeTable("Rewrap", func(cause errors.CausedBy) {
			alter := errors.ByRuntime
			if cause == errors.ByRuntime {
				alter = errors.ByBug
			}
			e := errors.Caused(alter, "foo")
			e = errors.NewCaused(cause, e)
			Ω(e).Should(MatchError("foo"))
			Ω(e.CausedBy).Should(Equal(cause))
		},
			Entry("ByBug", errors.ByBug),
			Entry("ByRuntime", errors.ByRuntime),
			Entry("ByExternal", errors.ByExternal),
			Entry("ByInput", errors.ByInput),
		)

	})

	Context("From error text", func() {

		It("Bug", func() {
			assertError(errors.Bug("foo"), "foo", errors.ByBug)
			assertError(errors.Bugf("foo %s", "bar"), "foo bar", errors.ByBug)
		})

		It("Runtime", func() {
			assertError(errors.Runtime("foo"), "foo", errors.ByRuntime)
			assertError(errors.Runtimef("foo %s", "bar"), "foo bar", errors.ByRuntime)
		})

		It("External", func() {
			assertError(errors.External("foo"), "foo", errors.ByExternal)
			assertError(errors.Externalf("foo %s", "bar"), "foo bar", errors.ByExternal)
		})

		It("Input", func() {
			assertError(errors.Input("foo"), "foo", errors.ByInput)
			assertError(errors.Inputf("foo %s", "bar"), "foo bar", errors.ByInput)
		})

	})

	Context("GetCausedBy", func() {

		It("Default to ByBug", func() {
			Ω(errors.GetCausedBy(syserr.New("foo"))).Should(Equal(errors.ByBug))
		})

		It("Error object", func() {
			Ω(errors.GetCausedBy(errors.External("foo"))).Should(Equal(errors.ByExternal))
		})

		It("nil is NoError", func() {
			Ω(errors.GetCausedBy(nil)).Should(Equal(errors.NoError))

		})

	})

	Context("GetPanicCausedBy", func() {

		It("nil", func() {
			Ω(errors.GetPanicCausedBy(nil)).Should(Equal(errors.NoError))
		})

		It("error", func() {
			Ω(errors.GetPanicCausedBy(errors.Input("foo"))).Should(Equal(errors.ByInput))
		})

		It("Other value", func() {
			Ω(errors.GetPanicCausedBy(0)).Should(Equal(errors.ByBug))
		})

	})

	DescribeTable("Caused", func(causedBy errors.CausedBy) {
		Ω(errors.Caused(causedBy, "foo").CausedBy).Should(Equal(causedBy))
	},
		Entry("ByInput", errors.ByInput),
		Entry("ByBug", errors.ByBug),
		Entry("ByExternal", errors.ByExternal),
		Entry("ByRuntime", errors.ByRuntime),
	)

	It("Causedf", func() {
		e := errors.Causedf(errors.ByInput, "foo %d", 3)
		Ω(e.CausedBy).Should(Equal(errors.ByInput))
		Ω(e.Error()).Should(HavePrefix("foo 3"))
	})

	DescribeTable("NewCaused", func(causedBy errors.CausedBy) {
		e := syserr.New("foo")
		er := errors.NewCaused(causedBy, e)
		Ω(er.CausedBy).Should(Equal(causedBy))
		Ω(er.Error()).Should(HavePrefix("foo"))
	},
		Entry("ByInput", errors.ByInput),
		Entry("ByBug", errors.ByBug),
		Entry("ByExternal", errors.ByExternal),
		Entry("ByRuntime", errors.ByRuntime),
	)
})
