//go:generate stringer -type=CausedBy

// Package errors Enhanced errors package, add CausedBy information. Can
// replace standard errors package.
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

// CausedBy describe who caused this error.
type CausedBy int

// Error is a special error interface contains CausedBy information.
type Error interface {
	error

	// CausedBy returns the source of error.
	CausedBy() CausedBy
}

const (
	// ByBug error caused by a bug. This kind of error normally logged and/or
	// report to error report service, notify to developer. Did not show detail
	// error to end-user, it is not their fault, do not blame them, and apologize
	// for this internal error.
	ByBug CausedBy = iota

	// ByRuntime error caused by golang runtime, such as run out of memory, file
	// read/write error. Any error caused by OS, or hardware. Network interface
	// error is caused by Runtime, error caused by network cable is not. Report
	// this kind of error to health monitor service, notify maintains team as
	// fast as possible. It not need to report to error service.
	ByRuntime

	// ByExternal error caused by depended external service, such as a Database
	// or other app services. Or network environment, such as lost network
	// connection. This kind of error is can not fixed by patching code, patching
	// OS, upgrade or replace hardware, they are not our error, nothing we can
	// do.  Report to health monitor service, notify maintains team to contact
	// people who can fix this.
	// Give end-user a brief message that who caused this error that they can
	// understand, such as payment service, and user know this error is not
	// caused by us.
	ByExternal

	// ByInput error caused by bad input. A program always dealing with input,
	// such as user input, or a request for a service daemon. When the input is
	// not expected, return error with detail and precise reason. Of course, do
	// not need report to error report service or health monitor service.
	ByInput

	// NoError is a special value returned by GetPanicCausedBy() to indicate no
	// error happened
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

// New function replace of standard errors.New(), create a ByBug error.
func New(text string) Error {
	return byBug{
		syserr.New(text),
	}
}

// NewBug wrap an exist error to ByBug. If e is nil, return nil. If e is
// already an Error, abort the wrap.
func NewBug(e error) Error {
	err, need := checkWrapped(e, ByBug)
	if need {
		return byBug{e}
	}
	return err
}

// NewRuntime wrap an exist error to ByRuntime. If e is nil, return nil. If e is
// already an Error, abort the wrap.
func NewRuntime(e error) Error {
	err, need := checkWrapped(e, ByRuntime)
	if need {
		return byRuntime{e}
	}
	return err
}

// NewExternal wrap an exist error to ByRuntime. If e is nil, return nil. If e is
// already an Error, abort the wrap.
func NewExternal(e error) Error {
	err, need := checkWrapped(e, ByExternal)
	if need {
		return byExternal{e}
	}
	return err
}

// NewInput wrap an exist error to ByRuntime. If e is nil, return nil. If e is
// already an Error, abort the wrap.
func NewInput(e error) Error {
	err, need := checkWrapped(e, ByInput)
	if need {
		return byInput{e}
	}
	return err
}

// check error dose need wrap, if not need, NexXXX() funcs use err as return value,
// if need wrap, needWrap returns true
func checkWrapped(e error, exp CausedBy) (err Error, needWrap bool) {
	if e == nil {
		return nil, false
	}

	err, ok := e.(Error)
	if !ok {
		return nil, true
	}

	return err, err.CausedBy() != exp
}

// Bug creates an Error from string.
func Bug(text string) Error {
	return byBug{syserr.New(text)}
}

// Bugf sprintf version of Bug().
func Bugf(text string, a ...interface{}) Error {
	return Bug(fmt.Sprintf(text, a...))
}

// Runtime creates an Error from string.
func Runtime(text string) Error {
	return byRuntime{syserr.New(text)}
}

// Runtimef sprintf version of Runtime().
func Runtimef(text string, a ...interface{}) Error {
	return Runtime(fmt.Sprintf(text, a...))
}

// External creates an Error from string.
func External(text string) Error {
	return byExternal{syserr.New(text)}
}

// Externalf sprintf version of Runtime().
func Externalf(text string, a ...interface{}) Error {
	return External(fmt.Sprintf(text, a...))
}

// Input creates an Error from string.
func Input(text string) Error {
	return byInput{syserr.New(text)}
}

// Inputf sprintf version of Input.
func Inputf(text string, a ...interface{}) Error {
	return Input(fmt.Sprintf(text, a...))
}

// Caused create error causedBy set by argument
func Caused(causedBy CausedBy, text string) Error {
	switch causedBy {
	case ByBug:
		return Bug(text)
	case ByInput:
		return Input(text)
	case ByExternal:
		return External(text)
	case ByRuntime:
		return Runtime(text)
	default:
		panic("Unknown causedBy")
	}
}

// Causedf sprintf version of Caused
func Causedf(causedBy CausedBy, text string, a ...interface{}) Error {
	return Caused(causedBy, fmt.Sprintf(text, a...))
}

// NewCaused wraps an exist error to specified causedBy Error,
// If e is nil, return nil.
// If already an Error, returned directly if causedBy matches, re-wrap
// with specific causedBy if not matched.
func NewCaused(causedBy CausedBy, err error) Error {
	switch causedBy {
	case ByBug:
		return NewBug(err)
	case ByInput:
		return NewInput(err)
	case ByExternal:
		return NewExternal(err)
	case ByRuntime:
		return NewRuntime(err)
	default:
		panic("Unknown causedBy")
	}
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
	switch err := e.(type) {
	case nil:
		return NoError
	case Error:
		return err.CausedBy()
	default:
		return ByBug
	}
}

// GetPanicCausedBy resolve caused for recover() return value, if the value is
// nil, return NoError. For error value use GetCausedBy() to resolve, other
// value return ByBug.
func GetPanicCausedBy(v interface{}) CausedBy {
	switch err := v.(type) {
	case nil:
		return NoError
	case Error:
		return err.CausedBy()
	default:
		return ByBug
	}
}
