package errors_test

import (
	"github.com/redforks/testing/reset"

	. "github.com/onsi/gomega"

	"context"

	. "github.com/onsi/ginkgo"
	"github.com/redforks/errors"
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
			Ω(err).Should(Equal(3))
			Ω(actx).Should(Equal(ctx))
		})

		errors.Handle(ctx, 3)
		Ω(called).Should(Equal(1))
	})

	It("Context is nil", func() {
		called := 0
		errors.SetHandler(func(ctx context.Context, err interface{}) {
			called++
			Ω(ctx).Should(Equal(context.Background()))
		})

		errors.Handle(nil, 2) // nolint:staticcheck
		Ω(called).Should(Equal(1))
	})

})
