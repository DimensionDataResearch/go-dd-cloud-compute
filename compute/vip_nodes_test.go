package compute

import "testing"

// Create VIP node (successful).
func TestClient_CreateVIPNode_Success(test *testing.T) {
	expect := expect(test)

	testClientRequest(test, &ClientTestConfig{
		Request: func(test *testing.T, client *Client) {
			vipNodeID, err := client.CreateVIPNode(NewVIPNodeConfiguration{
				Name:                "myProductionNode.1",
				Description:         "Production Server 1",
				IPv4Address:         "10.5.2.14",
				Status:              VIPNodeStatusEnabled,
				HealthMonitorID:     "0168b83a-d487-11e4-811f-005056806999",
				ConnectionLimit:     20000,
				ConnectionRateLimit: 2000,
				NetworkDomainID:     "553f26b6-2a73-42c3-a78b-6116f11291d0",
			})
			if err != nil {
				test.Fatal(err)
			}

			expect.EqualsString("VIPNodeID", "9e6b496d-5261-4542-91aa-b50c7f569c54", vipNodeID)
		},
		Respond: testValidateJSONRequestAndRespondOK(createVIPNodeTestResponse, &NewVIPNodeConfiguration{}, func(test *testing.T, requestBody interface{}) {
			verifyCreateVIPNodeTestRequest(test, requestBody.(*NewVIPNodeConfiguration))
		}),
	})
}

/*
 * Test requests.
 */

var createVIPNodeTestRequest = `
	{
		"networkDomainId": "553f26b6-2a73-42c3-a78b-6116f11291d0",
		"name": "myProductionNode.1",
		"description": "Production Server 1",
		"ipv4Address": "10.5.2.14",
		"status": "ENABLED",
		"healthMonitorId": "0168b83a-d487-11e4-811f-005056806999",
		"connectionLimit": "20000",
		"connectionRateLimit": "2000"
	}
`

func verifyCreateVIPNodeTestRequest(test *testing.T, request *NewVIPNodeConfiguration) {
	expect := expect(test)

	expect.NotNil("NewVIPNodeConfiguration", request)
	expect.EqualsString("NewVIPNodeConfiguration.NetworkDomainID", "553f26b6-2a73-42c3-a78b-6116f11291d0", request.NetworkDomainID)
	expect.EqualsString("NewVIPNodeConfiguration.Name", "myProductionNode.1", request.Name)
	expect.EqualsString("NewVIPNodeConfiguration.Description", "Production Server 1", request.Description)
	expect.EqualsString("NewVIPNodeConfiguration.IPv4Address", "10.5.2.14", request.IPv4Address)
	expect.EqualsString("NewVIPNodeConfiguration.IPv6Address", "", request.IPv6Address)
	expect.EqualsString("NewVIPNodeConfiguration.Status", VIPNodeStatusEnabled, request.Status)
	expect.EqualsString("NewVIPNodeConfiguration.HealthMonitorID", "0168b83a-d487-11e4-811f-005056806999", request.HealthMonitorID)
	expect.EqualsInt("NewVIPNodeConfiguration.ConnectionLimit", 20000, request.ConnectionLimit)
	expect.EqualsInt("NewVIPNodeConfiguration.ConnectionRateLimit", 2000, request.ConnectionRateLimit)
}

/*
 * Test responses.
 */

var createVIPNodeTestResponse = `
	{
        "requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad",
        "operation": "CREATE_NODE",
        "responseCode": "OK",
        "message": "Node 'myProductionNode.1' has been created.",
        "info": [
            {
                "name": "nodeId",
                "value": "9e6b496d-5261-4542-91aa-b50c7f569c54"
            },
            {
                "name": "name",
                "value": "myProductionNode.1"
            }
        ]
    }
 `

func verifyDeployVIPNodeTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "CREATE_NODE", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResponseCodeOK, response.ResponseCode)
	expect.EqualsString("Response.Message", "Node 'myProductionNode.1' has been created.", response.Message)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
	expect.EqualsInt("Response.Message.Length", 2, len(response.Message))
	expect.EqualsString("Response.FieldMessages[0].FieldName", "nodeId", response.FieldMessages[0].FieldName)
	expect.EqualsString("Response.FieldMessages[0].Message", "9e6b496d-5261-4542-91aa-b50c7f569c54", response.FieldMessages[0].Message)
	expect.EqualsString("Response.FieldMessages[1].FieldName", "name", response.FieldMessages[1].FieldName)
	expect.EqualsString("Response.FieldMessages[1].Message", "myProductionNode.1", response.FieldMessages[1].Message)
}
