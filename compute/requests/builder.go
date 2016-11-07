package requests

import "net/http"

// RequestBuilder is a delegate for building an HTTP request.
//
// The context parameter represents the current request.
//
// Returns an error to abort the request.
type RequestBuilder func(context RequestContext) (*http.Request, error)
