package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

/*
 * Integration test support
 */

// ClientTestConfig represents the configuration for a Client integration test.
type ClientTestConfig struct {
	Region         string
	User           string
	Password       string
	OrganizationID string
	ContentType    string

	Request ClientTestRequester
	Respond ClientTestResponder

	once sync.Once
}

// One-time initialisation for client test configuration.
func (clientTest *ClientTestConfig) initialize() {
	if clientTest.Request == nil {
		panic("ClientTest.Request is nil")
	}

	if clientTest.Respond == nil {
		panic("ClientTest.Respond is nil")
	}

	clientTest.Region = "AU"
	clientTest.User = "TestUser"
	clientTest.Password = "TestPassword"
	clientTest.OrganizationID = "my-organization-id"
	clientTest.ContentType = "application/json"
}

// EnsureInitialized ensures that one-time initialisation has been performed for the ClientTest.
func (clientTest *ClientTestConfig) EnsureInitialized() *ClientTestConfig {
	clientTest.once.Do(clientTest.initialize)

	return clientTest
}

// A function that invokes the request(s) for an integration test.
type ClientTestRequester func(test *testing.T, client *Client)

// A function that handles requests and generates responses for an integration test.
type ClientTestResponder func(test *testing.T, request *http.Request) (statusCode int, responseBody string)

// A function that validates a deserialised request body.
type ClientTestValidateRequestBody func(test *testing.T, requestBody interface{})

// Respond with HTTP OK (200) and the specified response body.
func testRespondOK(responseBody string) ClientTestResponder {
	return testRespond(http.StatusOK, responseBody)
}

// Respond with HTTP CREATED (201) and the specified response body.
func testRespondCreated(responseBody string) ClientTestResponder {
	return testRespond(http.StatusCreated, responseBody)
}

func testRespond(statusCode int, responseBody string) ClientTestResponder {
	return func(test *testing.T, request *http.Request) (int, string) {
		return statusCode, responseBody
	}
}

func testValidateJSONRequestAndRespondOK(responseBody string, requestBodyTemplate interface{}, validateRequestBody ClientTestValidateRequestBody) ClientTestResponder {
	return testValidateJSONRequestAndRespond(http.StatusOK, responseBody, requestBodyTemplate, validateRequestBody)
}

func testValidateJSONRequestAndRespond(statusCode int, responseBody string, requestBodyTemplate interface{}, validateRequestBody ClientTestValidateRequestBody) ClientTestResponder {
	return func(test *testing.T, request *http.Request) (int, string) {
		err := readRequestBodyAsJSON(request, requestBodyTemplate)
		if err != nil {
			test.Fatal("Failed to deserialise request body: ", err)
		}

		validateRequestBody(test, requestBodyTemplate)

		return statusCode, responseBody
	}
}

func testValidateXMLRequestAndRespondOK(responseBody string, requestBodyTemplate interface{}, validateRequestBody ClientTestValidateRequestBody) ClientTestResponder {
	return testValidateXMLRequestAndRespond(http.StatusOK, responseBody, requestBodyTemplate, validateRequestBody)
}

func testValidateXMLRequestAndRespond(statusCode int, responseBody string, requestBodyTemplate interface{}, validateRequestBody ClientTestValidateRequestBody) ClientTestResponder {
	return func(test *testing.T, request *http.Request) (int, string) {
		err := readRequestBodyAsXML(request, requestBodyTemplate)
		if err != nil {
			test.Fatal("Failed to deserialise request body: ", err)
		}

		validateRequestBody(test, requestBodyTemplate)

		return statusCode, responseBody
	}
}

func testClientRequest(test *testing.T, testConfiguration *ClientTestConfig) {
	testConfiguration.EnsureInitialized()

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		statusCode, response := testConfiguration.Respond(test, request)

		writer.Header().Set("Content-Type", testConfiguration.ContentType)
		writer.WriteHeader(statusCode)

		fmt.Fprintln(writer, response)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: testConfiguration.OrganizationID,
	})

	testConfiguration.Request(test, client)
}
