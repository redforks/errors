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
		Ω(e.Code().Caused()).Should(Equal(causedBy))
	}

	Context("Wrap error", func() {
		cases := []TableEntry{
			Entry("Bug", errors.ByBug, errors.NewBug),
			Entry("Runtime", errors.ByRuntime, errors.NewRuntime),
			Entry("External", errors.ByExternal, errors.NewExternal),
			Entry("Input", errors.ByInput, errors.NewInput),
			Entry("ClientBug", errors.ByClientBug, errors.NewClientBug),
		}

		DescribeTable("not nil", func(causedBy errors.CausedBy, fn func(err error) *errors.Error) {
			assertError(fn(syserr.New("foo")), "foo", causedBy)
		}, cases...)

		DescribeTable("nil", func(causedBy errors.CausedBy, fn func(err error) *errors.Error) {
			Ω(fn(nil)).Should(BeNil())
		}, cases...)

		DescribeTable("Rewrap", func(cause errors.CausedBy, fn func(err error) *errors.Error) {
			alter := errors.ByRuntime
			if cause == errors.ByRuntime {
				alter = errors.ByBug
			}
			e := errors.Caused(alter, "foo")
			e = fn(e)
			Ω(e.Error()).Should(Equal("foo"))
			Ω(e.Code().Caused()).Should(Equal(cause))
		}, cases...)

	})

	Context("From error text", func() {

		DescribeTable("without format", func(causedBy errors.CausedBy, fn func(msg string) *errors.Error) {
			assertError(fn("foo"), "foo", causedBy)
		},
			Entry("New", errors.ByBug, errors.New),
			Entry("Bug", errors.ByBug, errors.Bug),
			Entry("Runtime", errors.ByRuntime, errors.Runtime),
			Entry("External", errors.ByExternal, errors.External),
			Entry("Input", errors.ByInput, errors.Input),
			Entry("ClientBug", errors.ByClientBug, errors.ClientBug),
		)

		DescribeTable("with format", func(causedBy errors.CausedBy, fn func(text string, a ...interface{}) *errors.Error) {
			assertError(fn("foo %s", "bar"), "foo bar", causedBy)
		},
			Entry("Bug", errors.ByBug, errors.Bugf),
			Entry("Runtime", errors.ByRuntime, errors.Runtimef),
			Entry("External", errors.ByExternal, errors.Externalf),
			Entry("Input", errors.ByInput, errors.Inputf),
			Entry("ClientBug", errors.ByClientBug, errors.ClientBugf),
		)
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
		Ω(errors.Caused(causedBy, "foo").Code().Caused()).Should(Equal(causedBy))
	},
		Entry("ByInput", errors.ByInput),
		Entry("ByBug", errors.ByBug),
		Entry("ByExternal", errors.ByExternal),
		Entry("ByRuntime", errors.ByRuntime),
	)

	It("Causedf", func() {
		e := errors.Causedf(errors.ByInput, "foo %d", 3)
		Ω(e.Code().Caused()).Should(Equal(errors.ByInput))
		Ω(e.Error()).Should(Equal("foo 3"))
	})

	DescribeTable("NewCaused", func(causedBy errors.CausedBy) {
		e := syserr.New("foo")
		er := errors.NewCaused(causedBy, e)
		Ω(er.Code().Caused()).Should(Equal(causedBy))
		Ω(er.Error()).Should(Equal("foo"))
	},
		Entry("ByInput", errors.ByInput),
		Entry("ByBug", errors.ByBug),
		Entry("ByExternal", errors.ByExternal),
		Entry("ByRuntime", errors.ByRuntime),
		Entry("ByClientBug", errors.ByClientBug),
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
		Entry("Bugf", errors.Bugf("foo, %d", 1)),
		Entry("Runtime", errors.Runtime("foo")),
		Entry("Runtimef", errors.Runtimef("foo, %d", 1)),
		Entry("External", errors.External("foo")),
		Entry("Externalf", errors.Externalf("foo, %d", 1)),
		Entry("Input", errors.Input("foo")),
		Entry("Inputf", errors.Inputf("foo, %d", 1)),
		Entry("ClientBug", errors.ClientBugf("foo, %d", 1)),
		Entry("Caused", errors.Caused(errors.ByInput, "foo")),
		Entry("Causedf", errors.Causedf(errors.ByInput, "foo %d", 1)),
		Entry("NewCaused", errors.NewCaused(errors.ByInput, syserr.New("foo"))),
		Entry("Wrap", errors.Wrap(errors.ByInput, syserr.New("foo"), "bla")),
		Entry("Wrapf", errors.Wrapf(errors.ByInput, syserr.New("foo"), "bla %d", 1)),
	)

	Context("ErrorStack", func() {
		It("Include stack and msg", func() {
			e := errors.New("foo")
			msg := e.ErrorStack()
			Ω(msg).Should(HavePrefix("foo\n"))
			Ω(msg).Should(ContainSubstring("errors_test.go"))
			Ω(msg).ShouldNot(ContainSubstring("errors.go"))
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
			Ω(errors.ForLog(1)).Should(HavePrefix("1\n"))
		})

		It("inner error", func() {
			e := errors.Wrap(errors.ByInput, errors.Bug("foo"), "bar")
			s := errors.ForLog(e)
			idxBar := strings.Index(s, "bar")
			idxFoo := strings.Index(s, "Inner error:\nfoo")
			Ω(idxFoo > idxBar).Should(BeTrue())
			Ω(idxBar).Should(BeNumerically(">", -1))
			Ω(idxFoo).Should(BeNumerically(">", -1))
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
