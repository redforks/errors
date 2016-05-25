package errors_test

import (
	"github.com/redforks/testing/reset"

	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	"github.com/redforks/errors"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("handler", func() {

	BeforeEach(func() {
		reset.Enable()
	})

	AfterEach(func() {
		reset.Disable()
	})

	It("Context not nil", func() {
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

	It("Context is nil", func() {
		called := 0
		errors.SetHandler(func(ctx context.Context, err interface{}) {
			called++
			assert.Equal(t(), context.Background(), ctx)
		})

		errors.Handle(nil, 2)
		assert.Equal(t(), 1, called)
	})

})
