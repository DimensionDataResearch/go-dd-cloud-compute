package compute

import "fmt"

// IsOperationCancelledError determines if an error is an OperationCancelledError.
func IsOperationCancelledError(err error) bool {
	_, isOperationCancelledError := err.(*OperationCancelledError)

	return isOperationCancelledError
}

// OperationCancelledError is the error returned when an operation cancelled.
type OperationCancelledError struct {
	OperationDescription string
}

// Get a string representation of the error.
func (err OperationCancelledError) Error() string {
	return fmt.Sprintf("%s was cancelled.",
		err.OperationDescription,
	)
}

var _ error = &OperationCancelledError{}
