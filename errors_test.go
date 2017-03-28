package errors_test

import (
	syserr "errors"
	"strings"

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
			Ω(e.Error()).Should(Equal("foo"))
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
		Ω(e.Error()).Should(Equal("foo 3"))
	})

	DescribeTable("NewCaused", func(causedBy errors.CausedBy) {
		e := syserr.New("foo")
		er := errors.NewCaused(causedBy, e)
		Ω(er.CausedBy).Should(Equal(causedBy))
		Ω(er.Error()).Should(Equal("foo"))
	},
		Entry("ByInput", errors.ByInput),
		Entry("ByBug", errors.ByBug),
		Entry("ByExternal", errors.ByExternal),
		Entry("ByRuntime", errors.ByRuntime),
	)

	DescribeTable("stacktrace", func(e *errors.Error) {
		stack := e.Stack()
		Ω(stack).Should(ContainSubstring("errors_test.go"))
		Ω(stack).ShouldNot(ContainSubstring("errors.go"))
	},
		Entry("New", errors.New("foo")),
		Entry("NewBug", errors.NewBug(syserr.New("foo"))),
		Entry("NewRuntime", errors.NewRuntime(syserr.New("foo"))),
		Entry("NewExternal", errors.NewExternal(syserr.New("foo"))),
		Entry("NewInput", errors.NewInput(syserr.New("foo"))),
		Entry("Bug", errors.Bug("foo")),
		Entry("Bugf", errors.Bugf("foo, %s", 1)),
		Entry("Runtime", errors.Runtime("foo")),
		Entry("Runtimef", errors.Runtimef("foo, %s", 1)),
		Entry("External", errors.External("foo")),
		Entry("Externalf", errors.Externalf("foo, %s", 1)),
		Entry("Input", errors.Input("foo")),
		Entry("Inputf", errors.Inputf("foo, %s", 1)),
		Entry("Caused", errors.Caused(errors.ByInput, "foo")),
		Entry("Causedf", errors.Causedf(errors.ByInput, "foo %d", 1)),
		Entry("NewCaused", errors.NewCaused(errors.ByInput, syserr.New("foo"))),
		Entry("Wrap", errors.Wrap(errors.ByInput, syserr.New("foo"), "bla")),
		Entry("Wrapf", errors.Wrapf(errors.ByInput, syserr.New("foo"), "bla %s", 1)),
	)

	Context("ErrorStack", func() {
		It("Include stack and msg", func() {
			e := errors.New("foo")
			msg := e.ErrorStack()
			Ω(msg).Should(HavePrefix("foo\n"))
			Ω(msg).Should(ContainSubstring("errors_test.go"))
			Ω(msg).ShouldNot(ContainSubstring("errors.go"))
		})

		It("Include inner error", func() {
			e := errors.New("foo")
			e = errors.Wrap(errors.ByBug, e, "bar")
			msg := e.ErrorStack()
			Ω(msg).Should(ContainSubstring("foo"))
			Ω(msg).Should(ContainSubstring("bar"))
			Ω(msg).Should(ContainSubstring("errors_test.go"))
			Ω(strings.Count(msg, "errors_test.go")).Should(Equal(2))
		})

		It("Include inner inner error", func() {
			e := errors.New("foo")
			e = errors.Wrap(errors.ByBug, e, "bar")
			e = errors.Wrap(errors.ByBug, e, "blah")
			msg := e.ErrorStack()
			Ω(msg).Should(ContainSubstring("foo"))
			Ω(msg).Should(ContainSubstring("bar"))
			Ω(msg).Should(ContainSubstring("blah"))
			Ω(msg).Should(ContainSubstring("errors_test.go"))
			Ω(strings.Count(msg, "errors_test.go")).Should(Equal(3))
		})
	})

	Context("ForLog", func() {

		It("*Error", func() {
			Ω(errors.ForLog(errors.Bug("foo"))).Should(HavePrefix("foo\n"))
		})

		It("error", func() {
			Ω(errors.ForLog(syserr.New("foo"))).Should(Equal("foo"))
		})

		It("other", func() {
			Ω(errors.ForLog(1)).Should(Equal("1"))
		})
	})

	Context("Wrap", func() {

		It("Wrap", func() {
			inner := syserr.New("foo")
			e := errors.Wrap(errors.ByBug, inner, "bar")
			Ω(e.Error()).Should(Equal("bar"))
			Ω(e.Err).Should(Equal(inner))
		})

		It("Wrapf", func() {
			inner := syserr.New("foo")
			e := errors.Wrapf(errors.ByBug, inner, "foo %s", "bar")
			Ω(e.Error()).Should(Equal("foo bar"))
			Ω(e.Err).Should(Equal(inner))
		})

		It("Wrap any value", func() {
			e := errors.Wrap(errors.ByBug, "foo", "bar")
			Ω(e.Error()).Should(Equal("bar"))
			Ω(e.Err).Should(Equal(syserr.New("foo")))
		})

		It("Wrapf any value", func() {
			e := errors.Wrapf(errors.ByBug, "foo", "foo %s", "bar")
			Ω(e.Error()).Should(Equal("foo bar"))
			Ω(e.Err).Should(Equal(syserr.New("foo")))
		})
	})

})
