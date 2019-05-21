package compute

import (
	"github.com/pkg/errors"
)

// APIResponse represents the response to an API call.
type APIResponse interface {
	// GetMessage gets the message associated with the API response.
	GetMessage() string

	// GetResponseCode gets the response code associated with the API response.
	GetResponseCode() string

	// GetRequestID returns the request correlation ID.
	GetRequestID() string

	// GetAPIVersion gets the version of the API that returned the response.
	GetAPIVersion() string
}

// APIError is an error representing an error response from an API.
type APIError struct {
	Message  string
	Response APIResponse
}

// Error returns the error message associated with the APIError.
func (apiError *APIError) Error() string {
	return apiError.Message
}

var _ error = &APIError{}

// IsResourceBusyError determines whether the specified error represents a RESOURCE_BUSY response (or its v1 equivalent, REASON_392) from CloudControl.
func IsResourceBusyError(err error) bool {
	return IsAPIErrorCode(err, ResponseCodeResourceBusy) || IsAPIErrorCode(err, ResultCodeResourceBusy)
}

// IsResourceNotFoundError determines whether the specified error represents a RESOURCE_NOT_FOUND response from CloudControl.
func IsResourceNotFoundError(err error) bool {
	return IsAPIErrorCode(err, ResponseCodeResourceNotFound)
}

// IsNoIPAddressAvailableError determines whether the specified error represents a NO_IP_ADDRESS_AVAILABLE response from CloudControl.
func IsNoIPAddressAvailableError(err error) bool {
	return IsAPIErrorCode(err, ResponseCodeNoIPAddressAvailable)
}

// IsExceedsLimitError determines whether the specified error represents a EXCEEDS_LIMIT response (or its v1 equivalent, REASON_751) from CloudControl.
func IsExceedsLimitError(err error) bool {
	return IsAPIErrorCode(err, ResponseCodeExceedsLimit) || IsAPIErrorCode(err, ResultCodeExceedsLimit)
}

// IsAPIErrorCode determines whether the specified error represents a CloudControl API error with the specified response code.
func IsAPIErrorCode(err error, responseCode string) bool {
	apiError, ok := errors.Cause(err).(*APIError)
	if !ok {
		return false
	}

	return apiError.Response.GetResponseCode() == responseCode
}

// IsAPIError determines whether the specified error represents a CloudControl API error.
func IsAPIError(err error, responseCode string) bool {
	_, ok := errors.Cause(err).(*APIError)

	return ok
}

// Well-known API (v1) results

const (
	// ResultSuccess is a v1 API result indicating that an operation completed successfully.
	ResultSuccess = "SUCCESS"

	// ResultCodeSuccess is a v1 API result code indicating that an operation completed successfully.
	ResultCodeSuccess = "REASON_0"

	// ResultCodeResourceBusy is a v1 API result code indicating that an operation cannot be performed on a resource because the resource is busy.
	ResultCodeResourceBusy = "REASON_392"

	// ResultCodeServerNotFound is a v1 API result code indicating that an operation cannot be performed on a server because the server could not be found.
	ResultCodeServerNotFound = "REASON_395"

	// ResultCodeBackupNotEnabledForServer is a v1 API result code indicating that an operation cannot be performed on a server because backup is not enabled for that server.
	ResultCodeBackupNotEnabledForServer = "REASON_543"

	// ResultCodeBackupEnablementInProgressForServer is a v1 API result code indicating that an operation cannot be performed on a server because backup is in the process of being enabled for that server.
	ResultCodeBackupEnablementInProgressForServer = "REASON_544"

	// ResultCodeBackupNotEnabledForSubscription is a v1 API result code indicating that an operation cannot be performed on a server because backup is not enabled for that Subscription
	ResultCodeBackupNotEnabledForSubscription = "REASON_541"

	// ResultCodeBackupClientNotFound is a v1 API result code indicating that an operation cannot be performed on a backup client because the backup client does not exist (or was never actually installed on the target server).
	ResultCodeBackupClientNotFound = "REASON_545"

	// ResultCodeBackupJobInProgress is a v1 API result code indicating that an operation cannot be performed on a server because one or more backup jobs are in progress for that server.
	ResultCodeBackupJobInProgress = "REASON_547"

	// ResultCodeServerHasBackupAgents is a v1 API result code indicating that an operation cannot be performed on a server because the server still has one or more backup agents assigned to it.
	ResultCodeServerHasBackupAgents = "REASON_548"

	// ResultCodeBackupEnabledForServer is a v1 API result code indicating that an operation cannot be performed on a server because backup is enabled for that server.
	ResultCodeBackupEnabledForServer = "REASON_550"

	// ResultCodeExceedsLimit is a v1 API result code indicating that an operation cannot be performed on a resource because a resource limit was exceeded.
	ResultCodeExceedsLimit = "REASON_751"
)

// Well-known API (v2) response codes

const (
	// ResponseCodeOK is a v2 API response code indicating that an operation completed successfully.
	ResponseCodeOK = "OK"

	// ResponseCodeInProgress is a v2 API response code indicating that an operation is in progress.
	ResponseCodeInProgress = "IN_PROGRESS"

	// ResponseCodeResourceNotFound is a v2 API response code indicating that an operation failed because a target resource was not found.
	ResponseCodeResourceNotFound = "RESOURCE_NOT_FOUND"

	// ResponseCodeAuthorizationFailure is a v2 API response code indicating that an operation failed because the caller was not authorised to perform that operation (e.g. target resource belongs to another organisation).
	ResponseCodeAuthorizationFailure = "AUTHORIZATION_FAILURE"

	// ResponseCodeInvalidInputData is a v2 API response code indicating that an operation failed due to invalid input data.
	ResponseCodeInvalidInputData = "INVALID_INPUT_DATA"

	// ResponseCodeResourceNameNotUnique is a v2 API response code indicating that an operation failed due to the use of a name that duplicates an existing name.
	ResponseCodeResourceNameNotUnique = "NAME_NOT_UNIQUE"

	// ResponseCodeIPAddressNotUnique is a v2 API response code indicating that an operation failed due to the use of an IP address that duplicates an existing IP address.
	ResponseCodeIPAddressNotUnique = "IP_ADDRESS_NOT_UNIQUE"

	// ResponseCodeIPAddressOutOfRange is a v2 API response code indicating that an operation failed due to the use of an IP address lies outside the supported range (e.g. outside of the target subnet).
	ResponseCodeIPAddressOutOfRange = "IP_ADDRESS_OUT_OF_RANGE"

	// ResponseCodeNoIPAddressAvailable is a v2 API response code indicating that there are no remaining unreserved IPv4 addresses in the target subnet.
	ResponseCodeNoIPAddressAvailable = "NO_IP_ADDRESS_AVAILABLE"

	// ResponseCodeResourceHasDependency is a v2 API response code indicating that an operation cannot be performed on a resource because of a resource that depends on it.
	ResponseCodeResourceHasDependency = "HAS_DEPENDENCY"

	// ResponseCodeResourceBusy is a v2 API response code indicating that an operation cannot be performed on a resource because the resource is busy.
	ResponseCodeResourceBusy = "RESOURCE_BUSY"

	// ResponseCodeResourceLocked is a v2 API response code indicating that an operation cannot be performed on a resource because the resource is locked.
	ResponseCodeResourceLocked = "RESOURCE_LOCKED"

	// ResponseCodeExceedsLimit is a v2 API response code indicating that an operation failed because a resource limit was exceeded.
	ResponseCodeExceedsLimit = "EXCEEDS_LIMIT"

	// ResponseCodeOutOfResources is a v2 API response code indicating that an operation failed because some type of resource (e.g. free IPv4 addresses) has been exhausted.
	ResponseCodeOutOfResources = "OUT_OF_RESOURCES"

	// ResponseCodeOperationNotSupported is a v2 API response code indicating that an operation is not supported.
	ResponseCodeOperationNotSupported = "OPERATION_NOT_SUPPORTED"

	// ResponseCodeInfrastructureInMaintenance is a v2 API response code indicating that an operation failed due to maintenance being performed on the supporting infrastructure.
	ResponseCodeInfrastructureInMaintenance = "INFRASTRUCTURE_IN_MAINTENANCE"

	// ResponseCodeUnexpectedError is a v2 API response code indicating that the CloudControl API encountered an unexpected error.
	ResponseCodeUnexpectedError = "UNEXPECTED_ERROR"
)
