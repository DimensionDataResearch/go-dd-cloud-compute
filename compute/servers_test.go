package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Deploy server (successful).
func TestClient_DeployServer_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		deploymentConfiguration := &ServerDeploymentConfiguration{}
		err := readRequestBodyAsJSON(request, deploymentConfiguration)
		if err != nil {
			test.Fatal(err.Error())
		}

		verifyDeployServerRequest(test, deploymentConfiguration)

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
	}

	serverID, err := client.DeployServer(serverConfiguration)
	if err != nil {
		test.Fatal(err)
	}

	expect.equalsString("serverID", "7b62aae5-bdbe-4595-b58d-c78f95db2a7f", serverID)
}

/*
 * Test requests.
 */

var deployServerTestRequest = `
	{
		"name":"Production FTPS Server",
		"description":"This is the main FTPS Server",
		"imageId":"02250336-de2b-4e99-ab96-78511b7f8f4b",
		"start":true,
		"administratorPassword":"P$$ssWwrrdGoDd!",
		"cpu": {
			"count":4,
			"coresPerSocket":1,
			"speed":"STANDARD"
		},
		"memoryGb":4,
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

func verifyDeployServerRequest(test *testing.T, deploymentConfiguration *ServerDeploymentConfiguration) {
	expect := expect(test)

	expect.notNil("ServerDeploymentConfiguration", deploymentConfiguration)
	expect.equalsString("ServerDeploymentConfiguration.Name", "Production FTPS Server", deploymentConfiguration.Name)
	expect.equalsString("ServerDeploymentConfiguration.Description", "This is the main FTPS Server", deploymentConfiguration.Description)
	expect.equalsString("ServerDeploymentConfiguration.ImageID", "02250336-de2b-4e99-ab96-78511b7f8f4b", deploymentConfiguration.ImageID)
	expect.equalsString("ServerDeploymentConfiguration.AdministratorPassword", "password", deploymentConfiguration.AdministratorPassword)
}

/*
 * Test responses.
 */

var deployServerTestResponse = `
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

func verifyDeployServerTestResponse(test *testing.T, response *APIResponse) {
	expect := expect(test)

	expect.notNil("APIResponse", response)
	expect.equalsString("Response.Operation", "DEPLOY_SERVER", response.Operation)
	expect.equalsString("Response.ResponseCode", ResponseCodeInProgress, response.ResponseCode)
	expect.equalsString("Response.Message", "Request to deploy Server 'Production FTPS Server' has been accepted and is being processed.", response.Message)
	expect.equalsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}
