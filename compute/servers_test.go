package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetServer_ById_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, getServerTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	server, err := client.GetServer("5a32d6e4-9707-4813-a269-56ab4d989f4d")
	if err != nil {
		test.Fatal("Unable to deploy server: ", err)
	}

	verifyGetServerTestResponse(test, server)
}

// Deploy server (successful).
func TestClient_DeployServer_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		deploymentConfiguration := &ServerDeploymentConfiguration{}
		err := readRequestBodyAsJSON(request, deploymentConfiguration)
		if err != nil {
			test.Fatal(err.Error())
		}

		verifyDeployServerTestRequest(test, deploymentConfiguration)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, deployServerTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	serverConfiguration := ServerDeploymentConfiguration{
		Name:                  "Production FTPS Server",
		Description:           "This is the main FTPS Server",
		ImageID:               "02250336-de2b-4e99-ab96-78511b7f8f4b",
		AdministratorPassword: "password",
		CPU: VirtualMachineCPU{Count: 2},
	}

	serverID, err := client.DeployServer(serverConfiguration)
	if err != nil {
		test.Fatal(err)
	}

	expect.EqualsString("serverID", "7b62aae5-bdbe-4595-b58d-c78f95db2a7f", serverID)
}

// Add disk to server (successful).
func TestClient_AddServerDisk_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestBody := &addDiskToServer{}
		err := readRequestBodyAsJSON(request, requestBody)
		if err != nil {
			test.Fatal(err.Error())
		}

		verifyAddDiskToServerTestRequest(test, requestBody)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, addDiskToServerTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	diskID, err := client.AddDiskToServer("7b62aae5-bdbe-4595-b58d-c78f95db2a7f", 4, 20, "ECONOMY")
	if err != nil {
		test.Fatal(err)
	}

	expect.EqualsString("diskID", "9e6b496d-5261-4542-91aa-b50c7f569c54", diskID)
}

// Resize server disk (successful).
func TestClient_ResizeServerDisk_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		expect.EqualsString(
			"Request.URL",
			"/oec/0.9/dummy-organization-id/server/7b62aae5-bdbe-4595-b58d-c78f95db2a7f/disk/92b1819e-6f91-4abe-88c7-607841959f90/changeSize",
			request.URL.Path,
		)

		requestBody := &resizeServerDisk{}
		err := readRequestBodyAsXML(request, requestBody)
		if err != nil {
			test.Fatal(err.Error())
		}

		verifyResizeServerDiskRequest(test, requestBody)

		writer.Header().Set("Content-Type", "application/xml")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, resizeServerDiskTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	response, err := client.ResizeServerDisk("7b62aae5-bdbe-4595-b58d-c78f95db2a7f", "92b1819e-6f91-4abe-88c7-607841959f90", 23)
	if err != nil {
		test.Fatal(err)
	}

	verifyResizeServerDiskTestResponse(test, response)
}

// Delete Server (successful).
func TestClient_DeleteServer_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestBody, err := readRequestBodyAsString(request)
		if err != nil {
			test.Fatal("Failed to read request body: ", err)
		}

		expect.EqualsString("Request.Body",
			`{"id":"5b00a2ab-c665-4cd6-8291-0b931374fb3d"}`,
			requestBody,
		)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, deleteServerTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	err := client.DeleteServer("5b00a2ab-c665-4cd6-8291-0b931374fb3d")
	if err != nil {
		test.Fatal(err)
	}

	// Pass
}

/*
 * Test requests.
 */

const deployServerTestRequest = `
	{
		"name":"Production FTPS Server",
		"description":"This is the main FTPS Server",
		"imageId":"02250336-de2b-4e99-ab96-78511b7f8f4b",
		"start":true,
		"administratorPassword":"P$$ssWwrrdGoDd!",
		"memoryGb": 4,
		"cpu": {
			"count": 2
		},
		"primaryDns":"10.20.255.12",
		"secondaryDns":"10.20.255.13",
		"networkInfo": {
			"networkDomainId":"484174a2-ae74-4658-9e56-50fc90e086cf",
			"primaryNic" : {
				"vlanId":"0e56433f-d808-4669-821d-812769517ff8"
			},
			"additionalNic" : [
				{
					"privateIpv4" : "172.16.0.14"
				},
				{
					"vlanId":"e0b4d43c-c648-11e4-b33a-72802a5322b2"
				}
			]
		},
		"disk" : [
			{
				"scsiId" :"0",
				"speed" :"STANDARD"
			},
			{
				"scsiId" :"1" ,
				"speed" :"HIGHPERFORMANCE"
			}
		],
		"microsoftTimeZone":"035"
	}
`

func verifyDeployServerTestRequest(test *testing.T, deploymentConfiguration *ServerDeploymentConfiguration) {
	expect := expect(test)

	expect.NotNil("ServerDeploymentConfiguration", deploymentConfiguration)
	expect.EqualsString("ServerDeploymentConfiguration.Name", "Production FTPS Server", deploymentConfiguration.Name)
	expect.EqualsString("ServerDeploymentConfiguration.Description", "This is the main FTPS Server", deploymentConfiguration.Description)
	expect.EqualsString("ServerDeploymentConfiguration.ImageID", "02250336-de2b-4e99-ab96-78511b7f8f4b", deploymentConfiguration.ImageID)
	expect.EqualsString("ServerDeploymentConfiguration.AdministratorPassword", "password", deploymentConfiguration.AdministratorPassword)

	expect.EqualsInt("ServerDeploymentConfiguration.CPU.Count", 2, deploymentConfiguration.CPU.Count)
}

const addDiskToServerTestRequest = `
	{
		"id": "7b62aae5-bdbe-4595-b58d-c78f95db2a7f",
		"sizeGb": 20,
		"speed": "ECONOMY",
		"scsiId": 4
	}
`

func verifyAddDiskToServerTestRequest(test *testing.T, request *addDiskToServer) {
	expect := expect(test)

	expect.EqualsString("addDiskToServer.ServerID", "7b62aae5-bdbe-4595-b58d-c78f95db2a7f", request.ServerID)
	expect.EqualsString("addDiskToServer.Speed", "ECONOMY", request.Speed)
	expect.EqualsInt("addDiskToServer.SizeGB", 20, request.SizeGB)
	expect.EqualsInt("addDiskToServer.SCSIUnitID", 4, request.SCSIUnitID)
}

const resizeServerDiskTestRequest = `
	<ChangeDiskSize xmlns="http://oec.api.opsource.net/schemas/server">
		<newSizeGb>23</newSizeGb>
	</ChangeDiskSize>
`

func verifyResizeServerDiskRequest(test *testing.T, request *resizeServerDisk) {
	expect := expect(test)

	expect.NotNil("ReconfigureServer", request)
	expect.EqualsInt("ReconfigureServer.NewSizeGB", 23, request.NewSizeGB)
}

const notifyServerIPAddressChangeTestRequest = `
	{
		"nicId": "5999db1d-725c-46ba-9d4e-d33991e61ab1",
		"privateIpv4": "10.0.1.5",
		"ipv6": "fdfe::5a55:caff:fefa::1:9089"
	}
`

func verifyNotifyServerIPAddressChangeTestRequest(test *testing.T, request *notifyServerIPAddressChange) {
	expect := expect(test)

	expect.NotNil("NotifyServerIPAddressChange", request)
	expect.EqualsString("NotifyServerIPAddressChange.AdapterID", "5999db1d-725c-46ba-9d4e-d33991e61ab1", request.AdapterID)

	expect.NotNil("NotifyServerIPAddressChange.IPv4Address", request.IPv4Address)
	expect.EqualsString("NotifyServerIPAddressChange.IPv4Address", "10.0.1.5", *request.IPv4Address)

	expect.NotNil("NotifyServerIPAddressChange.IPv4Address", request.IPv6Address)
	expect.EqualsString("NotifyServerIPAddressChange.IPv6Address", "fdfe::5a55:caff:fefa::1:9089", *request.IPv6Address)
}

const reconfigureServerTestRequest = `
	{
		"memoryGb": 8,
		"cpuCount": 5,
		"cpuSpeed": "STANDARD",
		"coresPerSocket": 1,
		"id": "f8fe7965-3b7c-4cee-827e-f1e0b40a72e0"
	}
`

func verifyReconfigureServerTestRequest(test *testing.T, request *reconfigureServer) {
	expect := expect(test)

	expect.NotNil("ReconfigureServer", request)
	expect.EqualsString("ReconfigureServer.ServerID", "5999db1d-725c-46ba-9d4e-d33991e61ab1", request.ServerID)

	expect.NotNil("ReconfigureServer.MemoryGB", request.MemoryGB)
	expect.EqualsInt("ReconfigureServer.MemoryGB", 8, *request.MemoryGB)

	expect.NotNil("ReconfigureServer.CPUCount", request.CPUCount)
	expect.EqualsInt("ReconfigureServer.CPUCount", 5, *request.CPUCount)
}

/*
 * Test responses.
 */

const getServerTestResponse = `
	{
		"name": "Production Web Server",
		"description": "Server to host our main web application.",
		"operatingSystem": {
			"id": "WIN2008S32",
			"displayName": "WIN2008S/32",
			"family": "WINDOWS"
		},
		"cpu": {
			"count": 2,
			"speed": "STANDARD",
			"coresPerSocket": 1
		},
		"memoryGb": 4,
		"disk": [
			{
				"id": "c2e1f199-116e-4dbc-9960-68720b832b0a",
				"scsiId": 0,
				"sizeGb": 50,
				"speed": "STANDARD",
				"state": "NORMAL"
			}
		],
		"networkInfo": {
			"primaryNic": {
			"id": "5e869800-df7b-4626-bcbf-8643b8be11fd",
			"privateIpv4": "10.0.4.8",
			"ipv6": "2607:f480:1111:1282:2960:fb72:7154:6160",
			"vlanId": "bc529e20-dc6f-42ba-be20-0ffe44d1993f",
			"vlanName": "Production Server",
			"state": "NORMAL"
		},
		"additionalNic": [],
		"networkDomainId": "553f26b6-2a73-42c3-a78b-6116f11291d0" },
		"backup": {
			"assetId": "91002e08-8dc1-47a1-ad33-04f501c06f87",
			"servicePlan": "Advanced",
			"state": "NORMAL"
		},
		"monitoring": {
			"monitoringId": "11049",
			"servicePlan": "ESSENTIALS",
			"state": "NORMAL"
		},
		"softwareLabel": [
			"MSSQL2008R2S"
		],
		"sourceImageId": "3ebf3c0f-90fe-4a8b-8585-6e65b316592c",
		"createTime": "2015-12-02T10:31:33.000Z",
		"deployed": true,
		"started": true,
		"state": "PENDING_CHANGE",
		"progress": {
			"action": "SHUTDOWN_SERVER",
			"requestTime": "2015-12-02T11:07:40.000Z",
			"userName": "devuser1"
		},
		"vmwareTools": {
			"versionStatus": "CURRENT",
			"runningStatus": "RUNNING",
			"apiVersion": 9354
		},
		"virtualHardware": {
			"version": "vmx-08",
			"upToDate": false
		},
		"id": "5a32d6e4-9707-4813-a269-56ab4d989f4d",
		"datacenterId": "NA9"
	}
`

func verifyGetServerTestResponse(test *testing.T, server *Server) {
	expect := expect(test)

	expect.NotNil("Server", server)
	expect.EqualsString("Server.Name", "Production Web Server", server.Name)
	// TODO: Verify the rest of these fields.
	expect.EqualsString("Server.State", ResourceStatusPendingChange, server.State)
}

const deployServerTestResponse = `
	{
		"operation": "DEPLOY_SERVER",
		"responseCode": "IN_PROGRESS",
		"message": "Request to deploy Server 'Production FTPS Server' has been accepted and is being processed.",
		"info": [
			{
				"name": "serverId",
				"value": "7b62aae5-bdbe-4595-b58d-c78f95db2a7f"
			}
		],
		"warning": [],
		"error": [],
		"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
	}
`

func verifyDeployServerTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "DEPLOY_SERVER", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResponseCodeInProgress, response.ResponseCode)
	expect.EqualsString("Response.Message", "Request to deploy Server 'Production FTPS Server' has been accepted and is being processed.", response.Message)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}

const addDiskToServerTestResponse = `
	{
		"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad",
		"operation": "ADD_DISK",
		"responseCode": "IN_PROGRESS",
		"message": "The request to add 20GB Standard Speed Disk to Server 'SERVER-1' has been accepted and is being processed.",
		"info": [
			{
				"name": "diskId",
				"value": "9e6b496d-5261-4542-91aa-b50c7f569c54"
			},
			{
				"name": "scsiId",
				"value": "4"
			},
			{
				"name": "speed",
				"value": "STANDARD"
			}
		]
	}
`

func verifyAddDiskToServerTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "ADD_DISK", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResponseCodeInProgress, response.ResponseCode)
	expect.EqualsString("Response.Message", "The request to add 20GB Standard Speed Disk to Server 'SERVER-1' has been accepted and is being processed.", response.Message)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}

const resizeServerDiskTestResponse = `
	<Status>
		<operation>Change Server Disk Size</operation>
		<result>SUCCESS</result>
		<resultDetail>Server 'Change Server Disk Size' Issued</resultDetail>
		<resultCode>RESULT_0</resultCode>
	</Status>
`

func verifyResizeServerDiskTestResponse(test *testing.T, response *APIResponseV1) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "Change Server Disk Size", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResultSuccess, response.Result)
	expect.EqualsString("Response.Message", "Server 'Change Server Disk Size' Issued", response.Message)
	expect.EqualsString("Response.ResultCode", "RESULT_0", response.ResultCode)
}

const deleteServerTestResponse = `
	{
		"operation": "DELETE_SERVER",
		"responseCode": "IN_PROGRESS",
		"message": "Request to Delete Server (Id:5b00a2ab-c665-4cd6-8291-0b931374fb3d) has been accepted and is being processed",
		"info": [],
		"warning": [],
		"error": [],
		"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
	}
`

func verifyDeleteServerTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "DELETE_SERVER", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResponseCodeInProgress, response.ResponseCode)
	expect.EqualsString("Response.Message", "Request to Delete Server (Id:5b00a2ab-c665-4cd6-8291-0b931374fb3d) has been accepted and is being processed.", response.Message)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}
