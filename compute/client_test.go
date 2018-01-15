package compute

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return "", err
	}

	if os.Getenv("MCP_TEST_DUMP_REQUEST_BODY") != "" {
		fmt.Printf("RequestBody:\n%s", string(requestBody))
	}

	return string(requestBody), nil
}

func readRequestBodyAsJSON(request *http.Request, target interface{}) error {
	if request.Body == nil {
		return fmt.Errorf("request body is missing")
	}

	defer request.Body.Close()
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}

	if os.Getenv("MCP_TEST_DUMP_REQUEST_BODY") != "" {
		fmt.Printf("RequestBody:\n%s", string(requestBody))
	}

	return json.Unmarshal(requestBody, target)
}

func readRequestBodyAsXML(request *http.Request, target interface{}) error {
	if request.Body == nil {
		return fmt.Errorf("request body is missing")
	}

	defer request.Body.Close()
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}

	if os.Getenv("MCP_TEST_DUMP_REQUEST_BODY") != "" {
		fmt.Printf("RequestBody:\n%s", string(requestBody))
	}

	return xml.Unmarshal(requestBody, target)
}
