// Package compute contains the Go client for Dimension Data's cloud compute API.
package compute

import (
	"bytes"
	"encoding/xml"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Client is the client for Dimension Data's cloud compute API.
type Client struct {
	baseAddress string
	username    string
	password    string
	httpClient  *http.Client
}

// NewClient creates a new cloud compute API client.
// region is the cloud compute region identifier.
func NewClient(region string, username string, password string) *Client {
	baseAddress := fmt.Sprintf("https://api-%s.dimensiondata.com", region)

	return &Client{
		baseAddress,
		username,
		password,
		&http.Client{},
	}
}

// GetAccount retrieves the current user's account information
func (client *Client) GetAccount() (*Account, error) {
	request, err := client.newRequestV1("myaccount")
	if err != nil {
		return nil, err
	}
	if request.Body != nil {
		defer request.Body.Close()
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == 401 {
		return nil, fmt.Errorf("Cannot connect to compute API (invalid credentials).")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	account := &Account{}
	err = xml.Unmarshal(body, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// Create a basic request for the compute API (V1, XML, only GET currently supported).
func (client *Client) newRequestV1(relativeURI string) (*http.Request, error) {
	requestURI := fmt.Sprintf("%s/oec/0.9/%s", client.baseAddress, relativeURI)

	request, err := http.NewRequest(http.MethodGet, requestURI, nil)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(client.username, client.password)
	request.Header.Set("Accept", "text/xml")

	return request, nil
}

// Create a basic request for the compute API (V2.1, JSON).
func (client *Client) newRequestV21(relativeURI string, method string, body interface{}) (*http.Request, error) {
	requestURI := fmt.Sprintf("%s/caas/2.1/%s", client.baseAddress, relativeURI)

	var bodyReader io.Reader
	if body != nil {
		bodyData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		bodyReader = bytes.NewReader(bodyData)
	}

	request, err := http.NewRequest(method, requestURI, bodyReader)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(client.username, client.password)
	request.Header.Add("Accept", "application/json")

	if bodyReader != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	return request, nil
}

// SetBaseAddress configures the Client to use the specified base address.
func (client *Client) SetBaseAddress(baseAddress string) error {
	if len(baseAddress) == 0 {
		return fmt.Errorf("Must supply a valid base URI.")
	}

	client.baseAddress = baseAddress

	return nil
}
