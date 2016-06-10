package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Get VLAN by Id (successful).
func TestClient_GetVLAN_ById_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, getVLANTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	networkDomain, err := client.GetVLAN("0e56433f-d808-4669-821d-812769517ff8")
	if err != nil {
		test.Fatal(err)
	}

	verifyGetVLANTestResponse(test, networkDomain)
}

// List VLANs (successful).
func TestClient_ListVLANs_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, listVLANsTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	vlans, err := client.ListVLANs("484174a2-ae74-4658-9e56-50fc90e086cf")
	if err != nil {
		test.Fatal(err)
	}

	verifyListVLANsTestResponse(test, vlans)
}

// Deploy VLAN (successful).
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

// Edit VLAN (successful).
func TestClient_EditVlan_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestBody, err := readRequestBodyAsString(request)
		if err != nil {
			test.Fatal("Failed to read request body: ", err)
		}

		expect.equalsString("Request.Body",
			`{"id":"0e56433f-d808-4669-821d-812769517ff8","name":"Production VLAN","description":"For hosting our Production Cloud Servers"}`,
			requestBody,
		)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, editVLANTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	name := "Production VLAN"
	description := "For hosting our Production Cloud Servers"
	err := client.EditVLAN("0e56433f-d808-4669-821d-812769517ff8", &name, &description)
	if err != nil {
		test.Fatal(err)
	}

	// Pass
}

// Delete VLAN (successful).
func TestClient_DeleteVLAN_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestBody, err := readRequestBodyAsString(request)
		if err != nil {
			test.Fatal("Failed to read request body: ", err)
		}

		expect.equalsString("Request.Body",
			`{"id":"0e56433f-d808-4669-821d-812769517ff8"}`,
			requestBody,
		)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, deleteVLANTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	err := client.DeleteVLAN("0e56433f-d808-4669-821d-812769517ff8")
	if err != nil {
		test.Fatal(err)
	}

	// Pass
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
	expect.equalsString("DeployVLAN.ID", "484174a2-ae74-4658-9e56-50fc90e086cf", request.VLANID)
	expect.equalsString("DeployVLAN.Name", "Production VLAN", request.Name)
	expect.equalsString("DeployVLAN.Description", "For hosting our Production Cloud Servers", request.Description)
	expect.equalsString("DeployVLAN.IPv4BaseAddress", "10.0.3.0", request.IPv4BaseAddress)
	expect.equalsInt("DeployVLAN.IPv4PrefixSize", 23, request.IPv4PrefixSize)
}

var editVLANTestRequest = `
	{
		"id": "0e56433f-d808-4669-821d-812769517ff8",
		"name": "Production VLAN",
		"description": "For hosting our Production Cloud Servers"
	}
`

func verifyEditVLANTestRequest(test *testing.T, request *EditVLAN) {
	expect := expect(test)

	expect.notNil("EditVLAN", request)
	expect.equalsString("EditVLAN.ID", "0e56433f-d808-4669-821d-812769517ff8", request.ID)
	expect.notNil("EditVLAN.Name", request.Name)
	expect.equalsString("EditVLAN.Name", "Production VLAN", *request.Name)
	expect.notNil("EditVLAN.Description", request.Description)
	expect.equalsString("EditVLAN.Description", "For hosting our Production Cloud Servers", *request.Description)
}

/*
 * Test responses.
 */

var getVLANTestResponse = `
	{
		"networkDomain": {
			"id": "484174a2-ae74-4658-9e56-50fc90e086cf",
			"name": "Production Network Domain"
		},
		"name": "Production VLAN",
		"description": "For hosting our Production Cloud Servers",
		"privateIpv4Range": {
			"address": "10.0.3.0",
			"prefixSize": 24
		},
		"ipv4GatewayAddress": "10.0.3.1",
		"ipv6Range": {
			"address": "2607:f480:1111:1153:0:0:0:0",
			"prefixSize": 64
		},
		"ipv6GatewayAddress": "2607:f480:1111:1153:0:0:0:1",
		"createTime": "2016-06-09T07:21:34.000Z",
		"state": "NORMAL",
		"id": "0e56433f-d808-4669-821d-812769517ff8",
		"datacenterId": "NA9"
	}
`

func verifyGetVLANTestResponse(test *testing.T, vlan *VLAN) {
	expect := expect(test)

	expect.notNil("VLAN", vlan)
	expect.equalsString("VLAN.ID", "0e56433f-d808-4669-821d-812769517ff8", vlan.ID)
	expect.equalsString("VLAN.Name", "Production VLAN", vlan.Name)
	expect.equalsString("VLAN.Description", "For hosting our Production Cloud Servers", vlan.Description)
	expect.equalsString("VLAN.IPv4Range.BaseAddress", "10.0.3.0", vlan.IPv4Range.BaseAddress)
	expect.equalsInt("VLAN.IPv4Range.PrefixSize", 24, vlan.IPv4Range.PrefixSize)
	expect.equalsString("VLAN.IPv4GatewayAddress", "10.0.3.1", vlan.IPv4GatewayAddress)
	expect.equalsString("VLAN.IPv6Range.BaseAddress", "2607:f480:1111:1153:0:0:0:0", vlan.IPv6Range.BaseAddress)
	expect.equalsInt("VLAN.IPv6Range.PrefixSize", 64, vlan.IPv6Range.PrefixSize)
	expect.equalsString("VLAN.IPv6GatewayAddress", "2607:f480:1111:1153:0:0:0:1", vlan.IPv6GatewayAddress)
	expect.equalsString("VLAN.CreateTime", "2016-06-09T07:21:34.000Z", vlan.CreateTime)
	expect.equalsString("VLAN.State", "NORMAL", vlan.State)
	expect.equalsString("VLAN.DataCenterID", "NA9", vlan.DataCenterID)
}

var listVLANsTestResponse = `
	{
		"vlan": [
			{
				"networkDomain": {
					"id": "484174a2-ae74-4658-9e56-50fc90e086cf",
					"name": "Production Network Domain"
				},
				"name": "Production VLAN",
				"description": "For hosting our Production Cloud Servers",
				"privateIpv4Range": {
					"address": "10.0.3.0",
					"prefixSize": 24
				},
				"ipv4GatewayAddress": "10.0.3.1",
				"ipv6Range": {
					"address": "2607:f480:1111:1153:0:0:0:0",
					"prefixSize": 64
				},
				"ipv6GatewayAddress": "2607:f480:1111:1153:0:0:0:1",
				"createTime": "2016-06-09T07:21:34.000Z",
				"state": "NORMAL",
				"id": "0e56433f-d808-4669-821d-812769517ff8",
				"datacenterId": "NA9"
			}
		],
		"pageNumber": 1,
		"pageCount": 1,
		"totalCount": 1,
		"pageSize": 250
	}
`

func verifyListVLANsTestResponse(test *testing.T, vlans *VLANs) {
	expect := expect(test)

	expect.notNil("VLANs", vlans)

	expect.equalsInt("VLANs.PageCount", 1, vlans.PageCount)
	expect.equalsInt("VLANs.VLANs size", 1, len(vlans.VLANs))

	vlan1 := vlans.VLANs[0]
	expect.equalsString("VLANs.VLANs[0].ID", "0e56433f-d808-4669-821d-812769517ff8", vlan1.ID)
	expect.equalsString("VLANs.VLANs[0].Name", "Production VLAN", vlan1.Name)
	expect.equalsString("VLANs.VLANs[0].Description", "For hosting our Production Cloud Servers", vlan1.Description)
	expect.equalsString("VLANs.VLANs[0].DataCenterID", "NA9", vlan1.DataCenterID)

	expect.equalsString("VLANs.VLANs[0].VLAN.ID", "484174a2-ae74-4658-9e56-50fc90e086cf", vlan1.VLAN.ID)
	expect.equalsString("VLANs.VLANs[0].VLAN.Name", "Production Network Domain", vlan1.VLAN.Name)

	expect.equalsString("VLANs.VLANs[0].IPv4Range.BaseAddress", "10.0.3.0", vlan1.IPv4Range.BaseAddress)
	expect.equalsInt("VLANs.VLANs[0].IPv4Range.PrefixSize", 24, vlan1.IPv4Range.PrefixSize)
	expect.equalsString("VLANs.VLANs[0].IPv4GatewayAddress", "10.0.3.1", vlan1.IPv4GatewayAddress)

	expect.equalsString("VLANs.VLANs[0].IPv6Range.BaseAddress", "2607:f480:1111:1153:0:0:0:0", vlan1.IPv6Range.BaseAddress)
	expect.equalsInt("VLANs.VLANs[0].IPv6Range.PrefixSize", 64, vlan1.IPv6Range.PrefixSize)
	expect.equalsString("VLANs.VLANs[0].IPv6GatewayAddress", "2607:f480:1111:1153:0:0:0:1", vlan1.IPv6GatewayAddress)

	expect.equalsString("VLANs.VLANs[0].CreateTime", "2016-06-09T07:21:34.000Z", vlan1.CreateTime)
}

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

var editVLANTestResponse = `
	{
		"operation": "EDIT_VLAN",
		"responseCode": "OK",
		"message": "VLAN 'Production VLAN' was edited successfully.",
		"info": [],
		"warning": [],
		"error": [],
		"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
	}
`

func verifyEditVLANTestResponse(test *testing.T, response *APIResponse) {
	expect := expect(test)

	expect.notNil("APIResponse", response)
	expect.equalsString("Response.Operation", "EDIT_VLAN", response.Operation)
	expect.equalsString("Response.ResponseCode", ResponseCodeOK, response.ResponseCode)
	expect.equalsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}

var deleteVLANTestResponse = `
	{
		"operation": "DELETE_VLAN",
		"responseCode": "IN_PROGRESS",
		"message": "Request to Delete VLAN (Id: 0e56433f-d808-4669-821d-812769517ff8) has been accepted and is being processed.",
		"info": [],
		"warning": [],
		"error": [],
		"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
	}
`

func verifyDeleteVLANTestResponse(test *testing.T, response *APIResponse) {
	expect := expect(test)

	expect.notNil("APIResponse", response)
	expect.equalsString("Response.Operation", "DELETE_VLAN", response.Operation)
	expect.equalsString("Response.ResponseCode", ResponseCodeInProgress, response.ResponseCode)
	expect.equalsString("Response.Message", "Request to VLAN (Id: 0e56433f-d808-4669-821d-812769517ff8) has been accepted and is being processed.", response.Message)
	expect.equalsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}
