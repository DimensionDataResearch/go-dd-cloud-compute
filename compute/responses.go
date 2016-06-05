package compute

// Response represents the standard response from an API call.
type Response struct {
	Operation    string          `json:"operation"`
	ResponseCode string          `json:"responseCode"`
	Message      string          `json:"message"`
	Info         []NameValuePair `json:"info"`
	Warning      []NameValuePair `json:"warning"`
	Error        []NameValuePair `json:"error"`
	RequestId    string          `json:"requestId"`
}

// NameValuePair represents a name together with its associated value.
type NameValuePair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
