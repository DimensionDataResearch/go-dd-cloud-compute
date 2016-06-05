package compute

// ApiResponse represents common fields for all responses from an API call.
type ApiResponse struct {
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

// FieldMessage represents a field name together with an associated message.
type FieldMessage struct {
	// The field name.
	FieldName string `json:"name"`

	// The field message.
	Message string `json:"value"`
}
