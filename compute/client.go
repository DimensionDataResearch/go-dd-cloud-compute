// Package compute contains the Go client for Dimension Data's cloud compute API.
package compute

import (
	"encoding/xml"
	"fmt"
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

// GetMyAccount retrieves the current user's account information
func (client *Client) GetMyAccount() (*Account, error) {
	request, err := client.newRequestV1("myaccount", "GET")
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

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

// Create a basic request for the compute API (V1, XML).
func (client *Client) newRequestV1(relativeURI string, method string) (*http.Request, error) {
	requestURI := fmt.Sprintf("%s/oec/0.9/%s", client.baseAddress, relativeURI)

	request, err := http.NewRequest(method, requestURI, nil)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(client.username, client.password)
	request.Header.Add("Accept", "text/xml")

	return request, nil
}

// Create a basic request for the compute API (V2, JSON).
func (client *Client) newRequestV2(relativeURI string, method string) (*http.Request, error) {
	requestURI := fmt.Sprintf("%s/oec/0.9/%s", client.baseAddress, relativeURI)

	request, err := http.NewRequest(method, requestURI, nil)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(client.username, client.password)
	request.Header.Add("Accept", "application/json")

	return request, nil
}

// UseBaseAddress configures the Client to use the specified base address.
func (client *Client) UseBaseAddress(baseAddress string) error {
	if len(baseAddress) == 0 {
		return fmt.Errorf("Must supply a valid base URI.")
	}

	client.baseAddress = baseAddress

	return nil
}
