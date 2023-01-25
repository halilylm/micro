package web

import (
	"encoding/json"
	"github.com/dimfeld/httptreemux/v5"
	"net/http"
)

// Param returns the web call parameters from the request.
func Param(r *http.Request, key string) string {
	m := httptreemux.ContextParams(r.Context())
	return m[key]
}

// Decode reads the body of an HTTP requests looking for a JSON document.
// The body is decoded into the provided value
//
// If the provided value is a struct then it is checked for validation tags.
func Decode(r *http.Request, val any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}
	return nil
}
