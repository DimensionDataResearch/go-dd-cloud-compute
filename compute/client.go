// Package compute contains the client for Dimension Data's cloud compute API.
package compute

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// Client is the client for Dimension Data's cloud compute API.
type Client struct {
	baseAddress   string
	username      string
	password      string
	maxRetryCount int
	retryDelay    time.Duration
	stateLock     *sync.Mutex
	httpClient    *http.Client
	account       *Account
}

// NewClient creates a new cloud compute API client.
// region is the cloud compute region identifier.
func NewClient(region string, username string, password string) *Client {
	baseAddress := fmt.Sprintf("https://api-%s.dimensiondata.com", region)

	return &Client{
		baseAddress,
		username,
		password,
		0,
		0 * time.Second,
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

// ConfigureRetry configures the client's retry facility.
// Set maxRetryCount to 0 (the default) to disable retry.
func (client *Client) ConfigureRetry(maxRetryCount int, retryDelay time.Duration) {
	client.stateLock.Lock()
	defer client.stateLock.Unlock()

	client.maxRetryCount = maxRetryCount
	client.retryDelay = retryDelay
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

	log.Printf("Invoking '%s' request to '%s'...",
		request.Method,
		request.URL.String(),
	)

	response, err := client.httpClient.Do(request)
	if err != nil {
		log.Printf("Unexpected error while performing '%s' request to '%s': %s.",
			request.Method,
			request.URL.String(),
			err.Error(),
		)

		for retryCount := 0; retryCount < client.maxRetryCount; retryCount++ {
			log.Printf("Retrying '%s' request to '%s' (%d retries remaining)...",
				request.Method,
				request.URL.String(),
				retryCount-client.maxRetryCount,
			)

			response, err = client.httpClient.Do(request)

			if err != nil {
				log.Printf("Still failing - '%s' request to '%s': %s.",
					request.Method,
					request.URL.String(),
					err.Error(),
				)

				continue
			}

			log.Printf("'%s' request to '%s' succeeded.",
				request.Method,
				request.URL.String(),
			)

			break
		}

		if err != nil {
			err = fmt.Errorf("Unexpected error while performing '%s' request to '%s': %s",
				request.Method,
				request.URL.String(),
				err.Error(),
			)

			return
		}
	}
	defer response.Body.Close()

	statusCode = response.StatusCode

	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("Error reading response body for '%s': %s", request.URL.String(), err.Error())
	}

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

// Read an APIResponseV1 (as XML) from the response body.
func readAPIResponseV1(responseBody []byte, statusCode int) (apiResponse *APIResponseV1, err error) {
	apiResponse = &APIResponseV1{}
	err = xml.Unmarshal(responseBody, apiResponse)
	if err != nil {
		err = fmt.Errorf("Error reading API response (v1) from XML: %s", err.Error())

		return
	}

	if len(apiResponse.Result) == 0 {
		apiResponse.Result = "UNKNOWN_RESULT"
	}

	if len(apiResponse.Message) == 0 {
		apiResponse.Message = "An unexpected response was received from the compute API."
	}

	return
}

// Read an APIResponseV2 (as JSON) from the response body.
func readAPIResponseAsJSON(responseBody []byte, statusCode int) (apiResponse *APIResponseV2, err error) {
	apiResponse = &APIResponseV2{}
	err = json.Unmarshal(responseBody, apiResponse)
	if err != nil {
		err = fmt.Errorf("Error reading API response (v2) from JSON: %s", err.Error())

		return
	}

	if len(apiResponse.ResponseCode) == 0 {
		apiResponse.ResponseCode = "UNKNOWN_RESPONSE_CODE"
	}

	if len(apiResponse.Message) == 0 {
		apiResponse.Message = "An unexpected response was received from the compute API."
	}

	return
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
