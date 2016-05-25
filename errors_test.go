package errors_test

import (
	syserr "errors"

	"github.com/redforks/errors"

	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("errors", func() {

	assertError := func(e errors.Error, msg string, causedBy errors.CausedBy) {
		assert.Equal(t(), msg, e.Error())
		assert.Equal(t(), causedBy, e.CausedBy())
	}

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

	It("NewBug - nil", func() {
		assert.Nil(t(), errors.NewBug(nil))
	})

	It("NewRuntime - nil", func() {
		assert.Nil(t(), errors.NewRuntime(nil))
	})

	It("NewInput - nil", func() {
		assert.Nil(t(), errors.NewInput(nil))
	})

	It("NewExternal - nil", func() {
		assert.Nil(t(), errors.NewExternal(nil))
	})

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

	Context("GetCausedBy", func() {

		It("Default to ByBug", func() {
			assert.Equal(t(), errors.ByBug, errors.GetCausedBy(syserr.New("foo")))
		})

		It("Error object", func() {
			assert.Equal(t(), errors.ByExternal, errors.GetCausedBy(errors.External("foo")))
		})

	})

	Context("GetPanicCausedBy", func() {

		It("nil", func() {
			assert.Equal(t(), errors.NoError, errors.GetPanicCausedBy(nil))
		})

		It("error", func() {
			assert.Equal(t(), errors.ByInput, errors.GetPanicCausedBy(errors.Input("foo")))
		})

		It("Other value", func() {
			assert.Equal(t(), errors.ByBug, errors.GetPanicCausedBy(0))
		})

	})
})
