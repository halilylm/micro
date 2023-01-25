package web

// Middleware is a function designed to run some code before and/or after
// another Handler. It is designed to remove boilerplate or other concerns not
// direct to any given Handler
type Middleware func(Handler) Handler

// wrapMiddleware creates a new handler by wrapping middleware around a final
// handler with the new wrapped handler. Looping backwards ensures that the
// first middleware of the slice is the first to be executed by the requests.
func wrapMiddleware(mw []Middleware, handler Handler) Handler {
	// Look backwards through the middleware invoking each one. Reaplce the
	// handle with the new wrapped handler. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
