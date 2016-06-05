// Package compute contains the client for Dimension Data's cloud compute API.
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
	request, err := client.newRequestV1("myaccount", http.MethodGet, nil)
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

// Create a basic request for the compute API (V1, XML).
func (client *Client) newRequestV1(relativeURI string, method string, body interface{}) (*http.Request, error) {
	requestURI := fmt.Sprintf("%s/oec/0.9/%s", client.baseAddress, relativeURI)

	var (
		request		*http.Request
		bodyReader	io.Reader
		err			error
	)

	bodyReader, err = newReaderFromXML(body)
	if err != nil {
		return nil, err
	}

	request, err = http.NewRequest(method, requestURI, bodyReader)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(client.username, client.password)
	request.Header.Set("Accept", "text/xml")

	if bodyReader != nil {
		request.Header.Set("Content-Type", "text/xml")
	}

	return request, nil
}

// Create a basic request for the compute API (V2.2, JSON).
func (client *Client) newRequestV22(relativeURI string, method string, body interface{}) (*http.Request, error) {
	requestURI := fmt.Sprintf("%s/caas/2.2/%s", client.baseAddress, relativeURI)

	var (
		request		*http.Request
		bodyReader	io.Reader
		err			error
	)

	bodyReader, err = newReaderFromJSON(body)
	if err != nil {
		return nil, err
	}

	request, err = http.NewRequest(method, requestURI, bodyReader)
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

// newReaderFromJSON serialises the specified data as JSON and returns an io.Reader over that JSON.
func newReaderFromJSON(data interface{}) (io.Reader, error) {
	if data == nil {
		return nil, nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonData), nil
}

// newReaderFromXML serialises the specified data as XML and returns an io.Reader over that XML.
func newReaderFromXML(data interface{}) (io.Reader, error) {
	if data == nil {
		return nil, nil
	}

	xmlData, err := xml.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(xmlData), nil
}
