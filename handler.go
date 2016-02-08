package errors

import (
	"log"

	"golang.org/x/net/context"
)

var (
	handler Handler = defaultHandler
)

// Handler is a function do the actual error handling.
type Handler func(ctx context.Context, err interface{})

// Handle use handler to handle non-nil err value. Use SetHandler() to switch
// handler, default handler is a plain log.Print(), if ctx is nil, pass
// context.Background() to error handler.
func Handle(ctx context.Context, err interface{}) {
	if ctx == nil {
		ctx = context.Background()
	}
	handler(ctx, err)
}

// SetHandler switch error handler, NOTE: no sync lock to internal handler
// variable, only call SetHandler in application initialization code, to
// prevent data race.
func SetHandler(h Handler) {
	handler = h
}

func defaultHandler(ctx context.Context, err interface{}) {
	log.Print(err)
}
