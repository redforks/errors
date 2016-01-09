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

	bdd.It("Bug", func() {
		assertError(errors.Bug("foo"), "foo", errors.ByBug)
	})

	bdd.It("Runtime", func() {
		assertError(errors.Runtime("foo"), "foo", errors.ByRuntime)
	})

	bdd.It("External", func() {
		assertError(errors.External("foo"), "foo", errors.ByExternal)
	})

	bdd.It("Input", func() {
		assertError(errors.Input("foo"), "foo", errors.ByInput)
	})

	bdd.Context("GetCausedBy", func() {

		bdd.It("Default to ByBug", func() {
			assert.Equal(t(), errors.ByBug, errors.GetCausedBy(syserr.New("foo")))
		})

		bdd.It("Error object", func() {
			assert.Equal(t(), errors.ByExternal, errors.GetCausedBy(errors.External("foo")))
		})

	})
})
