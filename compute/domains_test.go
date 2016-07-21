package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Deploy network domain (successful).
func TestClient_DeployNetworkDomain_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestBody, err := readRequestBodyAsString(request)
		if err != nil {
			test.Fatal("Failed to read request body: ", err)
		}

		expect.EqualsString("Request.Body",
			`{"name":"A Network Domain","description":"This is a network domain","type":"ESSENTIALS","datacenterId":"AU9"}`,
			requestBody,
		)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, deployNetworkDomainTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	networkDomainID, err := client.DeployNetworkDomain(
		"A Network Domain",
		"This is a network domain",
		"ESSENTIALS",
		"AU9",
	)
	if err != nil {
		test.Fatal(err)
	}

	expect.EqualsString("NetworkDomainID", "f14a871f-9a25-470c-aef8-51e13202e1aa", networkDomainID)
}

// Edit network domain (successful).
func TestClient_EditNetworkDomain_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestBody, err := readRequestBodyAsString(request)
		if err != nil {
			test.Fatal("Failed to read request body: ", err)
		}

		expect.EqualsString("Request.Body",
			`{"id":"f14a871f-9a25-470c-aef8-51e13202e1aa","name":"A Network Domain","description":"This is a network domain","type":"ESSENTIALS"}`,
			requestBody,
		)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, editNetworkDomainTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	name := "A Network Domain"
	description := "This is a network domain"
	plan := "ESSENTIALS"

	err := client.EditNetworkDomain("f14a871f-9a25-470c-aef8-51e13202e1aa", &name, &description, &plan)
	if err != nil {
		test.Fatal(err)
	}

	// Pass
}

// Delete network domain (successful).
func TestClient_DeleteNetworkDomain_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestBody, err := readRequestBodyAsString(request)
		if err != nil {
			test.Fatal("Failed to read request body: ", err)
		}

		expect.EqualsString("Request.Body",
			`{"id":"f14a871f-9a25-470c-aef8-51e13202e1aa"}`,
			requestBody,
		)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, deleteNetworkDomainTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	err := client.DeleteNetworkDomain("f14a871f-9a25-470c-aef8-51e13202e1aa")
	if err != nil {
		test.Fatal(err)
	}

	// Pass
}

// Get network domain by Id (successful).
func TestClient_GetNetworkDomain_ById_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, networkDomainTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	networkDomain, err := client.GetNetworkDomain("8cdfd607-f429-4df6-9352-162cfc0891be")
	if err != nil {
		test.Fatal(err)
	}

	verifyNetworkDomainTestResponse(test, networkDomain)
}

// List network domains (successful).
func TestClient_ListNetworkDomains_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, listNetworkDomainsTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	networkDomains, err := client.ListNetworkDomains(nil)
	if err != nil {
		test.Fatal(err)
	}

	verifyListNetworkDomainsTestResponse(test, networkDomains)
}

/*
 * Test responses.
 */

var listNetworkDomainsTestResponse = `
{
	  "networkDomain": [
	    {
	      "name": "Domain 1",
	      "description": "This is test domain 1",
	      "type": "ESSENTIALS",
	      "snatIpv4Address": "168.128.17.63",
	      "createTime": "2016-01-12T22:33:05.000Z",
	      "state": "NORMAL",
	      "id": "75ab2a57-b75e-4ec6-945a-e8c60164fdf6",
	      "datacenterId": "AU9"
	    },
	    {
	      "name": "Domain 2",
	      "description": "",
	      "type": "ESSENTIALS",
	      "snatIpv4Address": "168.128.7.18",
	      "createTime": "2016-01-18T08:56:16.000Z",
	      "state": "NORMAL",
	      "id": "b91e0ba4-322c-32ca-bbc7-50b9a72d5f98",
	      "datacenterId": "AU10"
	    }
	  ],
	  "pageNumber": 1,
	  "pageCount": 2,
	  "totalCount": 2,
	  "pageSize": 250
	}
`

func verifyListNetworkDomainsTestResponse(test *testing.T, networkDomains *NetworkDomains) {
	expect := expect(test)

	expect.NotNil("NetworkDomains", networkDomains)

	expect.EqualsInt("NetworkDomains.PageCount", 2, networkDomains.PageCount)
	expect.EqualsInt("NetworkDomains.Domains size", 2, len(networkDomains.Domains))

	domain1 := networkDomains.Domains[0]
	expect.EqualsString("NetworkDomains.Domains[0].Name", "Domain 1", domain1.Name)

	domain2 := networkDomains.Domains[1]
	expect.EqualsString("NetworkDomains.Domains[1].Name", "Domain 2", domain2.Name)
}

var networkDomainTestResponse = `
	{
		"name": "Development Network Domain",
		"description": "This is a new Network Domain",
		"type": "ESSENTIALS",
		"snatIpv4Address": "165.180.9.252",
		"createTime": "2015-02-24T10:47:21.000Z",
		"state": "NORMAL",
		"id": "8cdfd607-f429-4df6-9352-162cfc0891be",
		"datacenterId": "NA9"
	}
`

func verifyNetworkDomainTestResponse(test *testing.T, networkDomain *NetworkDomain) {
	expect := expect(test)

	expect.NotNil("NetworkDomain", networkDomain)
	expect.EqualsString("NetworkDomain.ID", "8cdfd607-f429-4df6-9352-162cfc0891be", networkDomain.ID)
	expect.EqualsString("NetworkDomain.Name", "Development Network Domain", networkDomain.Name)
	expect.EqualsString("NetworkDomain.Type", "ESSENTIALS", networkDomain.Type)
	expect.EqualsString("NetworkDomain.State", "NORMAL", networkDomain.State)
	expect.EqualsString("NetworkDomain.NatIPv4Address", "165.180.9.252", networkDomain.NatIPv4Address)
	expect.EqualsString("NetworkDomain.DatacenterID", "NA9", networkDomain.DatacenterID)
}

var deployNetworkDomainTestResponse = `
	{
		"operation": "DEPLOY_NETWORK_DOMAIN",
		"responseCode": "IN_PROGRESS",
		"message": "Request to deploy Network Domain 'A Network Domain' has been accepted and is being processed.",
		"info": [
			{
				"name": "networkDomainId",
				"value": "f14a871f-9a25-470c-aef8-51e13202e1aa"
			}
		],
		"warning": [],
		"error": [],
		"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
	}
`

func verifyDeployNetworkDomainTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "DEPLOY_NETWORK_DOMAIN", response.Operation)
	expect.EqualsString("Response.ResponseCode", "IN_PROGRESS", response.ResponseCode)
	expect.EqualsString("Response.Message", "Request to deploy Network Domain 'A Network Domain' has been accepted and is being processed.", response.Message)
	expect.EqualsInt("Response.FieldMessages size", 1, len(response.FieldMessages))
	expect.EqualsString("Response.FieldMessages[0].Name", "networkDomainId", response.FieldMessages[0].FieldName)
	expect.EqualsString("Response.FieldMessages[0].Message", "f14a871f-9a25-470c-aef8-51e13202e1aa", response.FieldMessages[0].Message)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}

var editNetworkDomainTestResponse = `
	{
		"operation": "EDIT_NETWORK_DOMAIN",
		"responseCode": "OK",
		"message": "Network Domain 'Development Network Domain' was edited successfully.",
		"info": [],
		"warning": [],
		"error": [],
		"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
	}
`

func verifyEditNetworkDomainTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "EDIT_NETWORK_DOMAIN", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResponseCodeOK, response.ResponseCode)
	expect.EqualsString("Response.Message", "Network Domain 'Development Network Domain' was edited successfully.", response.Message)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}

var deleteNetworkDomainTestResponse = `
	{
		"operation": "DELETE_NETWORK_DOMAIN",
		"responseCode": "IN_PROGRESS",
		"message": "Request to Delete Network Domain (Id: 8cdfd607-f429-4df6-9352-162cfc0891be) has been accepted and is being processed.",
		"info": [],
		"warning": [],
		"error": [],
		"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
	}
`

func verifyDeleteNetworkDomainTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "DELETE_NETWORK_DOMAIN", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResponseCodeInProgress, response.ResponseCode)
	expect.EqualsString("Response.Message", "Request to Delete Network Domain (Id: 8cdfd607-f429-4df6-9352-162cfc0891be) has been accepted and is being processed.", response.Message)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}
