// Contains common structure handles errors in CausedBy way in command line
// utility applications.
//
// Service/daemon can use cmdline package to handle main goroutine's
// panic/errors, but can not handle panic/errors inside other goroutines, use
// redforks/life package to cover service goroutines.
package cmdline

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/redforks/errors"
)

type exitError int

func (code exitError) Error() string {
	return fmt.Sprintf("Exit error %d", int(code))
}

// Create a new exit error. Panic with exit error or return it in MainFunc,
// Go() detect it and call os.Exit() with specific exit code.
// os.Exit() exit the application immediately without calling deferred code
// block, by using exit error we can *fix* this.
func NewExitError(code int) error {
	return exitError(code)
}

// Your application main function type.
type MainFunc func() error

// Call your application main function, handles any error returned or paniced,
// handle error by errors.CausedBy rule.
func Go(main MainFunc) {
	defer func() {
		handleError(recover())
	}()

	handleError(main())
}

func handleError(v interface{}) {
	if err, ok := v.(exitError); ok {
		os.Exit(int(err))
		return
	}

	switch errors.GetPanicCausedBy(v) {
	case errors.NoError:
	case errors.ByBug, errors.ByRuntime:
		fmt.Fprintln(os.Stderr, v)
		debug.PrintStack()
		os.Exit(-1)
	case errors.ByInput, errors.ByExternal:
		fmt.Println(v)
		os.Exit(1)
	default:
		panic("Unknown CausedBy")
	}
}
