package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Deploy network domain (successful).
func TestClient_DeployVlan_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestBody, err := readRequestBodyAsString(request)
		if err != nil {
			test.Fatal("Failed to read request body: ", err)
		}

		expect.equalsString("Request.Body",
			`{"networkDomainId":"484174a2-ae74-4658-9e56-50fc90e086cf","name":"Production VLAN","description":"For hosting our Production Cloud Servers","privateIpv4BaseAddress":"10.0.3.0","privateIpv4PrefixSize":23}`,
			requestBody,
		)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, deployVLANTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	vlanID, err := client.DeployVLAN(
		"484174a2-ae74-4658-9e56-50fc90e086cf",
		"Production VLAN",
		"For hosting our Production Cloud Servers",
		"10.0.3.0",
		23,
	)
	if err != nil {
		test.Fatal(err)
	}

	expect.equalsString("VLANID", "0e56433f-d808-4669-821d-812769517ff8", vlanID)
}

/*
 * Test requests.
 */

var deployVLANTestRequest = `
	{
		"networkDomainId": "484174a2-ae74-4658-9e56-50fc90e086cf",
		"name": "Production VLAN",
		"description": "For hosting our Production Cloud Servers",
		"privateIpv4BaseAddress": "10.0.3.0",
		"privateIpv4PrefixSize": 23
	}
`

func verifyDeployVLANTestRequest(test *testing.T, request *DeployVLAN) {
	expect := expect(test)

	expect.notNil("DeployVLAN", request)
	expect.equalsString("DeployVLAN.NetworkDomainID", "484174a2-ae74-4658-9e56-50fc90e086cf", request.NetworkDomainID)
	expect.equalsString("DeployVLAN.Name", "Production VLAN", request.Name)
	expect.equalsString("DeployVLAN.Description", "For hosting our Production Cloud Servers", request.Description)
	expect.equalsString("DeployVLAN.IPv4BaseAddress", "10.0.3.0", request.IPv4BaseAddress)
	expect.equalsInt("DeployVLAN.IPv4PrefixSize", 23, request.IPv4PrefixSize)
}

/*
 * Test responses.
 */

var deployVLANTestResponse = `
	{
		"operation": "DEPLOY_VLAN",
		"responseCode": "IN_PROGRESS",
		"message": "Request to deploy VLAN 'Production VLAN' has been accepted and is being processed.",
		"info": [
			{
				"name": "vlanId",
				"value": "0e56433f-d808-4669-821d-812769517ff8"
			}
		],
		"warning": [],
		"error": [],
		"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
	}
`

func verifyDeployVLANTestResponse(test *testing.T, response *APIResponse) {
	expect := expect(test)

	expect.notNil("APIResponse", response)
	expect.equalsString("Response.Operation", "DEPLOY_VLAN", response.Operation)
	expect.equalsString("Response.ResponseCode", ResponseCodeInProgress, response.ResponseCode)
	expect.equalsString("Response.Message", "Request to deploy VLAN 'Production VLAN' has been accepted and is being processed.", response.Message)
	expect.equalsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}
