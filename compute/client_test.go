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

// Configure the Client to use the specified base address.
func (client *Client) setBaseAddress(baseAddress string) error {
	if len(baseAddress) == 0 {
		return fmt.Errorf("Must supply a valid base URI.")
	}

	client.baseAddress = baseAddress

	return nil
}

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
	responseBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}

	return xml.Unmarshal(responseBody, target)
}
