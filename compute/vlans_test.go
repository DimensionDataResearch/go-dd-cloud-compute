package compute

import "testing"

// Get VLAN by Id (successful).
func TestClient_GetVLAN_ById_Success(test *testing.T) {
	testClientRequest(test, &ClientTestConfig{
		Request: func(test *testing.T, client *Client) {
			vlan, err := client.GetVLAN("0e56433f-d808-4669-821d-812769517ff8")
			if err != nil {
				test.Fatal(err)
			}

			verifyGetVLANTestResponse(test, vlan)
		},
		Respond: testRespondOK(getVLANTestResponse),
	})
}

// List VLANs (successful).
func TestClient_ListVLANs_Success(test *testing.T) {
	testClientRequest(test, &ClientTestConfig{
		Request: func(test *testing.T, client *Client) {
			vlans, err := client.ListVLANs("484174a2-ae74-4658-9e56-50fc90e086cf", nil)
			if err != nil {
				test.Fatal(err)
			}

			verifyListVLANsTestResponse(test, vlans)
		},
		Respond: testRespondOK(listVLANsTestResponse),
	})
}

// Deploy VLAN (successful).
func TestClient_DeployVlan_Success(test *testing.T) {
	expect := expect(test)

	testClientRequest(test, &ClientTestConfig{
		Request: func(test *testing.T, client *Client) {
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

			expect.EqualsString("VLANID", "0e56433f-d808-4669-821d-812769517ff8", vlanID)
		},
		Respond: testValidateJSONRequestAndRespondOK(deployVLANTestResponse, &DeployVLAN{}, func(test *testing.T, requestBody interface{}) {
			verifyDeployVLANTestRequest(test, requestBody.(*DeployVLAN))
		}),
	})
}

// Edit VLAN (successful).
func TestClient_EditVlan_Success(test *testing.T) {
	testClientRequest(test, &ClientTestConfig{
		Request: func(test *testing.T, client *Client) {
			name := "Production VLAN"
			description := "For hosting our Production Cloud Servers"
			err := client.EditVLAN("0e56433f-d808-4669-821d-812769517ff8", &name, &description)
			if err != nil {
				test.Fatal(err)
			}

			// Pass
		},
		Respond: testValidateJSONRequestAndRespondOK(editVLANTestResponse, &EditVLAN{}, func(test *testing.T, requestBody interface{}) {
			verifyEditVLANTestRequest(test, requestBody.(*EditVLAN))
		}),
	})
}

// Delete VLAN (successful).
func TestClient_DeleteVLAN_Success(test *testing.T) {
	testClientRequest(test, &ClientTestConfig{
		Request: func(test *testing.T, client *Client) {
			err := client.DeleteVLAN("0e56433f-d808-4669-821d-812769517ff8")
			if err != nil {
				test.Fatal(err)
			}

			// Pass
		},
		Respond: testValidateJSONRequestAndRespondOK(deleteVLANTestResponse, &DeleteVLAN{}, func(test *testing.T, requestBody interface{}) {
			verifyDeleteVLANTestRequest(test, requestBody.(*DeleteVLAN))
		}),
	})
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

	expect.NotNil("DeployVLAN", request)
	expect.EqualsString("DeployVLAN.ID", "484174a2-ae74-4658-9e56-50fc90e086cf", request.VLANID)
	expect.EqualsString("DeployVLAN.Name", "Production VLAN", request.Name)
	expect.EqualsString("DeployVLAN.Description", "For hosting our Production Cloud Servers", request.Description)
	expect.EqualsString("DeployVLAN.IPv4BaseAddress", "10.0.3.0", request.IPv4BaseAddress)
	expect.EqualsInt("DeployVLAN.IPv4PrefixSize", 23, request.IPv4PrefixSize)
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

	expect.NotNil("EditVLAN", request)
	expect.EqualsString("EditVLAN.ID", "0e56433f-d808-4669-821d-812769517ff8", request.ID)
	expect.NotNil("EditVLAN.Name", request.Name)
	expect.EqualsString("EditVLAN.Name", "Production VLAN", *request.Name)
	expect.NotNil("EditVLAN.Description", request.Description)
	expect.EqualsString("EditVLAN.Description", "For hosting our Production Cloud Servers", *request.Description)
}

var deleteVLANTestRequest = `
	{
		"id":"0e56433f-d808-4669-821d-812769517ff8"
	}
`

func verifyDeleteVLANTestRequest(test *testing.T, request *DeleteVLAN) {
	expect := expect(test)

	expect.NotNil("DeleteVLAN", request)
	expect.EqualsString("DeleteVLAN.ID", "0e56433f-d808-4669-821d-812769517ff8", request.ID)
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

	expect.NotNil("VLAN", vlan)
	expect.EqualsString("VLAN.ID", "0e56433f-d808-4669-821d-812769517ff8", vlan.ID)
	expect.EqualsString("VLAN.Name", "Production VLAN", vlan.Name)
	expect.EqualsString("VLAN.Description", "For hosting our Production Cloud Servers", vlan.Description)
	expect.EqualsString("VLAN.IPv4Range.BaseAddress", "10.0.3.0", vlan.IPv4Range.BaseAddress)
	expect.EqualsInt("VLAN.IPv4Range.PrefixSize", 24, vlan.IPv4Range.PrefixSize)
	expect.EqualsString("VLAN.IPv4GatewayAddress", "10.0.3.1", vlan.IPv4GatewayAddress)
	expect.EqualsString("VLAN.IPv6Range.BaseAddress", "2607:f480:1111:1153:0:0:0:0", vlan.IPv6Range.BaseAddress)
	expect.EqualsInt("VLAN.IPv6Range.PrefixSize", 64, vlan.IPv6Range.PrefixSize)
	expect.EqualsString("VLAN.IPv6GatewayAddress", "2607:f480:1111:1153:0:0:0:1", vlan.IPv6GatewayAddress)
	expect.EqualsString("VLAN.CreateTime", "2016-06-09T07:21:34.000Z", vlan.CreateTime)
	expect.EqualsString("VLAN.State", "NORMAL", vlan.State)
	expect.EqualsString("VLAN.DataCenterID", "NA9", vlan.DataCenterID)
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

	expect.NotNil("VLANs", vlans)

	expect.EqualsInt("VLANs.PageCount", 1, vlans.PageCount)
	expect.EqualsInt("VLANs.VLANs size", 1, len(vlans.VLANs))

	vlan1 := vlans.VLANs[0]
	expect.EqualsString("VLANs.VLANs[0].ID", "0e56433f-d808-4669-821d-812769517ff8", vlan1.ID)
	expect.EqualsString("VLANs.VLANs[0].Name", "Production VLAN", vlan1.Name)
	expect.EqualsString("VLANs.VLANs[0].Description", "For hosting our Production Cloud Servers", vlan1.Description)
	expect.EqualsString("VLANs.VLANs[0].DataCenterID", "NA9", vlan1.DataCenterID)

	expect.EqualsString("VLANs.VLANs[0].NetworkDomain.ID", "484174a2-ae74-4658-9e56-50fc90e086cf", vlan1.NetworkDomain.ID)
	expect.EqualsString("VLANs.VLANs[0].NetworkDomain.Name", "Production Network Domain", vlan1.NetworkDomain.Name)

	expect.EqualsString("VLANs.VLANs[0].IPv4Range.BaseAddress", "10.0.3.0", vlan1.IPv4Range.BaseAddress)
	expect.EqualsInt("VLANs.VLANs[0].IPv4Range.PrefixSize", 24, vlan1.IPv4Range.PrefixSize)
	expect.EqualsString("VLANs.VLANs[0].IPv4GatewayAddress", "10.0.3.1", vlan1.IPv4GatewayAddress)

	expect.EqualsString("VLANs.VLANs[0].IPv6Range.BaseAddress", "2607:f480:1111:1153:0:0:0:0", vlan1.IPv6Range.BaseAddress)
	expect.EqualsInt("VLANs.VLANs[0].IPv6Range.PrefixSize", 64, vlan1.IPv6Range.PrefixSize)
	expect.EqualsString("VLANs.VLANs[0].IPv6GatewayAddress", "2607:f480:1111:1153:0:0:0:1", vlan1.IPv6GatewayAddress)

	expect.EqualsString("VLANs.VLANs[0].CreateTime", "2016-06-09T07:21:34.000Z", vlan1.CreateTime)
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

func verifyDeployVLANTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "DEPLOY_VLAN", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResponseCodeInProgress, response.ResponseCode)
	expect.EqualsString("Response.Message", "Request to deploy VLAN 'Production VLAN' has been accepted and is being processed.", response.Message)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
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

func verifyEditVLANTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "EDIT_VLAN", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResponseCodeOK, response.ResponseCode)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
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

func verifyDeleteVLANTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "DELETE_VLAN", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResponseCodeInProgress, response.ResponseCode)
	expect.EqualsString("Response.Message", "Request to VLAN (Id: 0e56433f-d808-4669-821d-812769517ff8) has been accepted and is being processed.", response.Message)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}
