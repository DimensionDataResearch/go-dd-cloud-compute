// Package compute contains the Go client for Dimension Data's cloud compute API.
package compute

import "net/http"

var apiV1BaseURI = "/oec/0.9/"
var apiV2BaseURI = "caas/2.1/"

// Create a basic request for the compute API (V1, XML).
func baseRequestV1(baseURI string, method string, username string, password string) (*http.Request, error) {
	request, err := http.NewRequest(method, baseURI, nil)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(username, password)
	request.Header.Add("Accept", "text/xml")

	return request, nil
}

// Create a basic request for the compute API (V2, JSON).
func baseRequestV2(baseURI string, method string, username string, password string) (*http.Request, error) {
	request, err := http.NewRequest(method, baseURI, nil)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(username, password)
	request.Header.Add("Accept", "application/json")

	return request, nil
}
