package compute

import (
	"encoding/xml"
	"fmt"
)

// APIResponseV1 represents a response from the CloudControl v1 API for an asynchronous operation.
type APIResponseV1 struct {
	// The XML name for the "APIResponseV1" data contract
	XMLName xml.Name `xml:"Status"`

	// The operation for which status is being reported.
	Operation string `xml:"operation"`

	// The operation result.
	Result string `xml:"result"`

	// A brief message describing the operation result.
	Message string `xml:"resultDetail"`

	// The operation result code
	ResultCode string `xml:"resultCode"`

	// Additional information (if any).
	AdditionalInformation []APIResponseAdditionalInformationV1 `xml:"additionalInformation"`
}

// APIResponseAdditionalInformationV1 represents additional information in a V1 API response (in the form of a name / value pair).
type APIResponseAdditionalInformationV1 struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value"`
}

// GetMessage gets the message associated with the API response.
func (response *APIResponseV1) GetMessage() string {
	return response.Message
}

// GetResponseCode gets the response code associated with the API response.
func (response *APIResponseV1) GetResponseCode() string {
	return response.Result
}

// GetRequestID gets the request correlation ID.
func (response *APIResponseV1) GetRequestID() string {
	return "NOT-PRESENT-IN-V1-API"
}

// GetAPIVersion gets the response code associated with the API response.
func (response *APIResponseV1) GetAPIVersion() string {
	return "v1"
}

var _ APIResponse = &APIResponseV1{}

// GetAdditionalInformation retrieves additional information (if available) by name from the API response.
//
// Returns nil if no matching additional information is found with the specified name.
func (response *APIResponseV1) GetAdditionalInformation(name string) *string {
	for _, additionalInformation := range response.AdditionalInformation {
		if additionalInformation.Name == name {
			return &additionalInformation.Value
		}
	}

	return nil
}

// ToError creates an error representing the API response.
func (response *APIResponseV1) ToError(errorMessageOrFormat string, formatArgs ...interface{}) error {
	return &APIError{
		Message:  fmt.Sprintf(errorMessageOrFormat, formatArgs...),
		Response: response,
	}
}
