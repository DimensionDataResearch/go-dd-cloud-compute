package compute

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
 * Test helpers
 */

// Pre-cache account details for the client.
func (client *Client) setAccount(account *Account) {
	client.stateLock.Lock()
	defer client.stateLock.Unlock()

	client.account = account
}

func readRequestBodyAsString(request *http.Request) (string, error) {
	if request.Body == nil {
		return "", nil
	}

	defer request.Body.Close()
	responseBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
}

func readRequestBodyAsJSON(request *http.Request, target interface{}) error {
	if request.Body == nil {
		return fmt.Errorf("Request body is missing.")
	}

	defer request.Body.Close()
	responseBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(responseBody, target)
}

func readRequestBodyAsXML(request *http.Request, target interface{}) error {
	if request.Body == nil {
		return fmt.Errorf("Request body is missing.")
	}

	defer request.Body.Close()
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}

	return xml.Unmarshal(requestBody, target)
}
