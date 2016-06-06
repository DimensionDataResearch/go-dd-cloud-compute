// Package compute contains the client for Dimension Data's cloud compute API.
package compute

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

// Client is the client for Dimension Data's cloud compute API.
type Client struct {
	baseAddress string
	username    string
	password    string
	stateLock	*sync.Mutex
	httpClient  *http.Client
	account		*Account
}

// NewClient creates a new cloud compute API client.
// region is the cloud compute region identifier.
func NewClient(region string, username string, password string) *Client {
	baseAddress := fmt.Sprintf("https://api-%s.dimensiondata.com", region)

	return &Client{
		baseAddress,
		username,
		password,
		&sync.Mutex{},
		&http.Client{},
		nil,
	}
}

// Reset clears all cached data from the Client.
func (client *Client) Reset() {
	client.stateLock.Lock()
	defer client.stateLock.Unlock()

	client.account = nil
}

// GetAccount retrieves the current user's account information
func (client *Client) GetAccount() (*Account, error) {
	client.stateLock.Lock()
	defer client.stateLock.Unlock()

	if client.account != nil {
		return client.account, nil
	}

	request, err := client.newRequestV1("myaccount", http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	if statusCode == 401 {
		return nil, fmt.Errorf("Cannot connect to compute API (invalid credentials).")
	}

	account := &Account{}
	err = xml.Unmarshal(responseBody, account)
	if err != nil {
		return nil, err
	}

	client.account = account

	return account, nil
}

// ListNetworkDomains retrieves a list of all network domains.
// TODO: Support filtering and sorting.
func (client *Client) ListNetworkDomains() (domains *NetworkDomains, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/networkDomain", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		// TODO: Consider reading response as APIResponse and returning response code / message.

		return nil, fmt.Errorf("Request failed with status code %d.", statusCode)
	}

	domains = &NetworkDomains{}
	err = json.Unmarshal(responseBody, domains)
	if err != nil {
		return nil, err
	}

	return domains, nil
}

// getOrganizationID gets the current user's organisation Id.
func (client *Client) getOrganizationID() (organizationID string, err error) {
	account, err := client.GetAccount()
	if err != nil {
		return "", err
	}

	return account.OrganizationID, nil
}

// Create a basic request for the compute API (V1, XML).
func (client *Client) newRequestV1(relativeURI string, method string, body interface{}) (*http.Request, error) {
	requestURI := fmt.Sprintf("%s/oec/0.9/%s", client.baseAddress, relativeURI)

	var (
		request    *http.Request
		bodyReader io.Reader
		err        error
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

// executeRequest performs the specified request and returns the entire response body, together with the HTTP status code.
func (client *Client) executeRequest(request *http.Request) (responseBody []byte, statusCode int, err error) {
	if request.Body != nil {
		defer request.Body.Close()
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()

	statusCode = response.StatusCode

	responseBody, err = ioutil.ReadAll(response.Body)

	return
}

// Create a basic request for the compute API (V2.2, JSON).
func (client *Client) newRequestV22(relativeURI string, method string, body interface{}) (*http.Request, error) {
	requestURI := fmt.Sprintf("%s/caas/2.2/%s", client.baseAddress, relativeURI)

	var (
		request    *http.Request
		bodyReader io.Reader
		err        error
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
