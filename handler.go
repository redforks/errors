package errors

import "log"

var (
	handler Handler = defaultHandler
)

type Handler func(err interface{})

// Handle use handler to handle non-nil err value. Use SetHandler() to switch
// handler, default handler is a plain log.Print()
func Handle(err interface{}) {
}

// SetHandler switch error handler, NOTE: no sync lock to internal handler
// variable, only call SetHandler in application initialization code, to
// prevent data race.
func SetHandler(h Handler) {
	handler = h
}

func defaultHandler(err interface{}) {
	log.Print(err)
}
