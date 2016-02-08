package errors_test

import (
	"spork/testing/reset"

	"golang.org/x/net/context"

	bdd "github.com/onsi/ginkgo"
	"github.com/redforks/errors"
	"github.com/stretchr/testify/assert"
)

var _ = bdd.Describe("handler", func() {

	bdd.BeforeEach(func() {
		reset.Enable()
	})

	bdd.AfterEach(func() {
		reset.Disable()
	})

	bdd.It("Context not nil", func() {
		called := 0
		ctx := context.WithValue(context.Background(), "foo", 1)

		errors.SetHandler(func(actx context.Context, err interface{}) {
			called++
			assert.Equal(t(), 3, err)
			assert.Equal(t(), ctx, actx)
		})

		errors.Handle(ctx, 3)
		assert.Equal(t(), 1, called)
	})

	bdd.It("Context is nil", func() {
		called := 0
		errors.SetHandler(func(ctx context.Context, err interface{}) {
			called++
			assert.Equal(t(), context.Background(), ctx)
		})

		errors.Handle(nil, 2)
		assert.Equal(t(), 1, called)
	})

})
