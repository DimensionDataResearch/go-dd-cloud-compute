package requests

// RequestContext represents the context for an outgoing HTTP request.
type RequestContext interface {
	// Retrieve a brief description of the current action being performed by the client
	GetAction() string

	// Retrieve the current user's organisation ID.
	GetOrganizationID() string

	// Retrieve the number of times the request has been retried.
	GetRetryCount() int

	// Retrieve the number of retries remaining for the request.
	GetRemainingRetryCount() int

	// Retrieve the last error (if any) encountered while processing the current request.
	GetLastError() error
}
