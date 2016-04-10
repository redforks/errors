package errors_test

import (
	syserr "errors"

	"github.com/redforks/errors"

	bdd "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

var _ = bdd.Describe("errors", func() {

	assertError := func(e errors.Error, msg string, causedBy errors.CausedBy) {
		assert.Equal(t(), msg, e.Error())
		assert.Equal(t(), causedBy, e.CausedBy())
	}

	bdd.It("New", func() {
		assertError(errors.New("foo"), "foo", errors.ByBug)
	})

	bdd.It("NewBug", func() {
		assertError(errors.NewBug(syserr.New("foo")), "foo", errors.ByBug)
	})

	bdd.It("NewRuntime", func() {
		assertError(errors.NewRuntime(syserr.New("foo")), "foo", errors.ByRuntime)
	})

	bdd.It("NewExternal", func() {
		assertError(errors.NewExternal(syserr.New("foo")), "foo", errors.ByExternal)
	})

	bdd.It("NewInput", func() {
		assertError(errors.NewInput(syserr.New("foo")), "foo", errors.ByInput)
	})

	bdd.It("NewBug - nil", func() {
		assert.Nil(t(), errors.NewBug(nil))
	})

	bdd.It("NewRuntime - nil", func() {
		assert.Nil(t(), errors.NewRuntime(nil))
	})

	bdd.It("NewInput - nil", func() {
		assert.Nil(t(), errors.NewInput(nil))
	})

	bdd.It("NewExternal - nil", func() {
		assert.Nil(t(), errors.NewExternal(nil))
	})

	bdd.It("Bug", func() {
		assertError(errors.Bug("foo"), "foo", errors.ByBug)
		assertError(errors.Bug("foo %s", "bar"), "foo bar", errors.ByBug)
	})

	bdd.It("Runtime", func() {
		assertError(errors.Runtime("foo"), "foo", errors.ByRuntime)
		assertError(errors.Runtime("foo %s", "bar"), "foo bar", errors.ByRuntime)
	})

	bdd.It("External", func() {
		assertError(errors.External("foo"), "foo", errors.ByExternal)
		assertError(errors.External("foo %s", "bar"), "foo bar", errors.ByExternal)
	})

	bdd.It("Input", func() {
		assertError(errors.Input("foo"), "foo", errors.ByInput)
		assertError(errors.Input("foo %s", "bar"), "foo bar", errors.ByInput)
	})

	bdd.Context("GetCausedBy", func() {

		bdd.It("Default to ByBug", func() {
			assert.Equal(t(), errors.ByBug, errors.GetCausedBy(syserr.New("foo")))
		})

		bdd.It("Error object", func() {
			assert.Equal(t(), errors.ByExternal, errors.GetCausedBy(errors.External("foo")))
		})

	})

	bdd.Context("GetPanicCausedBy", func() {

		bdd.It("nil", func() {
			assert.Equal(t(), errors.NoError, errors.GetPanicCausedBy(nil))
		})

		bdd.It("error", func() {
			assert.Equal(t(), errors.ByInput, errors.GetPanicCausedBy(errors.Input("foo")))
		})

		bdd.It("Other value", func() {
			assert.Equal(t(), errors.ByBug, errors.GetPanicCausedBy(0))
		})

	})
})
