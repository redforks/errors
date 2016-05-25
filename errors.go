//go:generate stringer -type=CausedBy

// Enhanced errors package, add CausedBy information. Can replace standard
// errors package.
//
// Error classified by who caused the error:
//
//  1. Bug: we programmers are the one responsible for the error.
//  2. Runtime: OS, hardware, golang, code libraries, normally can be fixed by
//  DevOp.
//  3. External: caused by external service, or network connection.
//  4. Input: caused by bad input.
//
// The first two kinds can be fixed by ourself, and the last two we can not
// fixed must report back.
package errors

import (
	syserr "errors"
	"fmt"
)

// Describe who caused this error.
type CausedBy int

// Extended error interface that describe error source (CausedBy)
type Error interface {
	error

	CausedBy() CausedBy
}

const (
	// Error caused by a bug. This kind of error normally logged and/or report to
	// error report service, notify to developer. Did not show detail error to
	// end-user, it is not their fault, do not blame them, and apologize for this
	// internal error.
	ByBug CausedBy = iota

	// Error caused by golang runtime, such as run out of memory, file read/write
	// error. Any error caused by OS, or hardware. Network interface error is
	// caused by Runtime, error caused by network cable is not. Report this kind
	// of error to health monitor service, notify maintains team as fast as
	// possible. It not need to report to error service.
	ByRuntime

	// Error caused by depended external service, such as a Database or other app
	// services. Or network environment, such as lost network connection. This
	// kind of error is can not fixed by patching code, patching OS, upgrade or
	// replace hardware, they are not our error, nothing we can do.
	// Report to health monitor service, notify maintains team to contact people
	// who can fix this.
	// Give end-user a brief message that who caused this error that they can
	// understand, such as payment service, and user know this error is not
	// caused by us.
	ByExternal

	// Error caused by bad input. A program always dealing with input, such as
	// user input, or a request for a service daemon. When the input is not
	// expected, return error with detail and precise reason. Of course, do not
	// need report to error report service or health monitor service.
	ByInput

	// A special value returned by GetPanicCausedBy() to indicate no error
	// happened
	NoError
)

type errorWrap struct {
	error
}

type byBug errorWrap

func (byBug) CausedBy() CausedBy {
	return ByBug
}

type byRuntime errorWrap

func (byRuntime) CausedBy() CausedBy {
	return ByRuntime
}

type byExternal errorWrap

func (byExternal) CausedBy() CausedBy {
	return ByExternal
}

type byInput errorWrap

func (byInput) CausedBy() CausedBy {
	return ByInput
}

// Replacement of standard errors.New(), create a ByBug error.
func New(text string) Error {
	return byBug{
		syserr.New(text),
	}
}

// Create a ByBug error from exist error. If e is nil, return nil. It is safe
// to write:
//
//  return errors.NewBug(exitFunc())
func NewBug(e error) Error {
	err, need := checkWrapped(e)
	if need {
		return byBug{e}
	}
	return err
}

// Create a ByRuntime error from exist error. If e is nil, return nil. It is
// safe to write:
//
//  return errors.NewRuntime(exitFunc())
func NewRuntime(e error) Error {
	err, need := checkWrapped(e)
	if need {
		return byRuntime{e}
	}
	return err
}

// Create a ByExternal error from exist error. If e is nil, return nil. It is
// safe to write:
//
//  return errors.NewExternal(exitFunc())
func NewExternal(e error) Error {
	err, need := checkWrapped(e)
	if need {
		return byExternal{e}
	}
	return err
}

// Create a ByInput error from exist error. If e is nil, return nil. It is safe
// to write:
//
//  return errors.NewInput(exitFunc())
func NewInput(e error) Error {
	err, need := checkWrapped(e)
	if need {
		return byInput{e}
	}
	return err
}

// check error dose need wrap, if not need, NexXXX() funcs use err as return value,
// if need wrap, needWrap returns true
func checkWrapped(e error) (err Error, needWrap bool) {
	if e == nil {
		return nil, false
	}

	err, ok := e.(Error)
	needWrap = !ok
	return
}

// Create a text ByBug error, use fmt.Sprintf() if contains extra arguments.
func Bug(text string) Error {
	return byBug{syserr.New(text)}
}

// Bugf printf version of Bug().
func Bugf(text string, a ...interface{}) Error {
	return Bug(fmt.Sprintf(text, a...))
}

// Create a text ByRuntime error, use fmt.Sprintf() if contains extra argumets.
func Runtime(text string) Error {
	return byRuntime{syserr.New(text)}
}

// Runtimef printf version of Runtime
func Runtimef(text string, a ...interface{}) Error {
	return Runtime(fmt.Sprintf(text, a...))
}

// Create a text ByExternal error, use fmt.Sprintf() if contains extra arguments.
func External(text string) Error {
	return byExternal{syserr.New(text)}
}

// Externalf printf version of External
func Externalf(text string, a ...interface{}) Error {
	return External(fmt.Sprintf(text, a...))
}

// Create a text ByInput error, use fmt.Sprintf() if contains extra arguments.
func Input(text string) Error {
	return byInput{syserr.New(text)}
}

// Inputf printf version of Input.
func Inputf(text string, a ...interface{}) Error {
	return Input(fmt.Sprintf(text, a...))
}

// GetCausedBy from any error. If the error is Error interface, call its
// CausedBy() method. Then all considered as ByBug.
//
// If the error is not a bug, wrap it use NewXXX() function before return:
//
//  if err := os.Open("file"); err != nil {
//		return errors.NewRuntime(err)
//  }
//
func GetCausedBy(e error) CausedBy {
	if err, ok := e.(Error); ok {
		return err.CausedBy()
	}
	return ByBug
}

// Resolve caused for recover() return value, if the value is nil, return
// NoError. For error value use GetCausedBy() to resolve, other value return
// ByBug.
func GetPanicCausedBy(v interface{}) CausedBy {
	if v == nil {
		return NoError
	}

	if val, ok := v.(error); ok {
		return GetCausedBy(val)
	}

	return ByBug
}
