package cmdline

import (
	"os"
	"spork/life"
	"spork/testing/reset"

	"golang.org/x/net/context"

	bdd "github.com/onsi/ginkgo"
	"github.com/redforks/errors"
	"github.com/redforks/hal"
	"github.com/stretchr/testify/assert"
)

var _ = bdd.Describe("cmdline", func() {
	var (
		exitCodes        []int
		onAbort, onError int
	)

	bdd.BeforeEach(func() {
		exitCodes = nil
		hal.Exit = func(n int) {
			exitCodes = append(exitCodes, n)
		}
		reset.Enable()

		onAbort, onError = 0, 0
		life.RegisterHook("log", 0, life.OnAbort, func() {
			onAbort++
		})

		errors.SetHandler(func(_ context.Context, err interface{}) {
			onError++
		})
	})

	bdd.AfterEach(func() {
		reset.Disable()
		hal.Exit = os.Exit
		errors.SetHandler(nil)
	})

	bdd.It("Without error", func() {
		hit := 0
		Go(func() error {
			hit++
			return nil
		})
		assert.Equal(t(), 1, hit)
		assert.Empty(t(), exitCodes)
		assert.Equal(t(), 0, onError)
	})

	bdd.It("Exit", func() {
		Go(func() error {
			return NewExitError(1)
		})
		assert.Equal(t(), []int{1}, exitCodes)
		assert.Equal(t(), 1, onAbort)
		assert.Equal(t(), 0, onError)
	})

	bdd.It("Report error", func() {
		Go(func() error {
			return errors.New("foo")
		})
		assert.Equal(t(), 1, onError)
	})

})
