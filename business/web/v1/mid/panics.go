package mid

import (
	"context"
	"fmt"
	"github.com/halilylm/micro/business/web/metrics"
	"github.com/halilylm/micro/foundation/web"
	"net/http"
	"runtime/debug"
)

func Panics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			// defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))
					metrics.AddPanics(ctx)
				}
			}()
			return handler(ctx, w, r)
		}
		return h
	}
	return m
}
