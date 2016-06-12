package cmdline

import (
	"os"
	"spork/life"

	"github.com/redforks/testing/reset"

	. "github.com/onsi/gomega"

	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	"github.com/redforks/errors"
	"github.com/redforks/hal"
)

var _ = Describe("cmdline", func() {
	var (
		exitCodes        []int
		onAbort, onError int
	)

	BeforeEach(func() {
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

	AfterEach(func() {
		reset.Disable()
		hal.Exit = os.Exit
		errors.SetHandler(nil)
	})

	It("Without error", func() {
		hit := 0
		Go(func() error {
			hit++
			return nil
		})
		Ω(hit).Should(Equal(1))
		Ω(exitCodes).Should(BeEmpty())
		Ω(onError).Should(Equal(0))
	})

	It("Exit", func() {
		Go(func() error {
			return NewExitError(1)
		})
		Ω(exitCodes).Should(Equal([]int{1}))
		Ω(onAbort).Should(Equal(1))
		Ω(onError).Should(Equal(0))
	})

	It("Report error", func() {
		Go(func() error {
			return errors.New("foo")
		})
		Ω(onError).Should(Equal(1))
	})

})
