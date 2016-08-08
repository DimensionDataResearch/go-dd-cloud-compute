package compute

import "fmt"

// APIResponseV2 represents the basic response most commonly received when making v2 API calls.
type APIResponseV2 struct {
	// The operation that was performed.
	Operation string `json:"operation"`

	// The API response code.
	ResponseCode string `json:"responseCode"`

	// The API status message (if any).
	Message string `json:"message"`

	// Informational messages (if any) relating to request fields.
	FieldMessages []FieldMessage `json:"info"`

	// Warning messages (if any) relating to request fields.
	FieldWarnings []FieldMessage `json:"warning"`

	// Error messages (if any) relating to request fields.
	FieldErrors []FieldMessage `json:"error"`

	// The request ID (correlation identifier).
	RequestID string `json:"requestId"`
}

// GetMessage gets the message associated with the API response.
func (response *APIResponseV2) GetMessage() string {
	return response.Message
}

// GetResponseCode gets the response code associated with the API response.
func (response *APIResponseV2) GetResponseCode() string {
	return response.ResponseCode
}

// GetRequestID gets the request correlation ID.
func (response *APIResponseV2) GetRequestID() string {
	return response.RequestID
}

// GetAPIVersion gets the response code associated with the API response.
func (response *APIResponseV2) GetAPIVersion() string {
	return "v2"
}

// GetFieldError retrieves the value of the specified field error message (if any).
// Returns nil if the no field error message with the specified name is present in the API response.
func (response *APIResponseV2) GetFieldError(fieldName string) *string {
	for index := range response.FieldErrors {
		errorMessage := response.FieldErrors[index]
		if errorMessage.FieldName == fieldName {
			return &errorMessage.Message
		}
	}

	return nil
}

// GetFieldWarning retrieves the value of the specified field warning message (if any).
// Returns nil if the no field warning message with the specified name is present in the API response.
func (response *APIResponseV2) GetFieldWarning(fieldName string) *string {
	for index := range response.FieldWarnings {
		warningMessage := response.FieldWarnings[index]
		if warningMessage.FieldName == fieldName {
			return &warningMessage.Message
		}
	}

	return nil
}

// GetFieldMessage retrieves the value of the specified field message (if any).
// Returns nil if the no field message with the specified name is present in the API response.
func (response *APIResponseV2) GetFieldMessage(fieldName string) *string {
	for index := range response.FieldMessages {
		fieldMessage := response.FieldMessages[index]
		if fieldMessage.FieldName == fieldName {
			return &fieldMessage.Message
		}
	}

	return nil
}

var _ APIResponse = &APIResponseV2{}

// ToError creates an error representing the API response.
func (response *APIResponseV2) ToError(errorMessageOrFormat string, formatArgs ...interface{}) error {
	return &APIError{
		Message:  fmt.Sprintf(errorMessageOrFormat, formatArgs...),
		Response: response,
	}
}

// FieldMessage represents a field name together with an associated message.
type FieldMessage struct {
	// The field name.
	FieldName string `json:"name"`

	// The field message.
	Message string `json:"value"`
}
