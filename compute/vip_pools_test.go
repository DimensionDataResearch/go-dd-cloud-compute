package compute

import "testing"

// Create VIP pool (successful).
func TestClient_CreateVIPPool_Success(test *testing.T) {
	expect := expect(test)

	testClientRequest(test, &ClientTestConfig{
		Request: func(test *testing.T, client *Client) {
			vipPoolID, err := client.CreateVIPPool(NewVIPPoolConfiguration{
				NetworkDomainID:   "553f26b6-2a73-42c3-a78b-6116f11291d0",
				Name:              "myDevelopmentPool.1",
				Description:       "Pool for load balancing development application servers.",
				LoadBalanceMethod: LoadBalanceMethodRoundRobin,
				ServiceDownAction: ServiceDownActionReselect,
				SlowRampTime:      10,
				HealthMonitorIDs: []string{
					"01683574-d487-11e4-811f-005056806999",
					"0168546c-d487-11e4-811f-005056806999",
				},
			})
			if err != nil {
				test.Fatal(err)
			}

			expect.EqualsString("VIPPoolID", "4d360b1f-bc2c-4ab7-9884-1f03ba2768f7", vipPoolID)
		},
		Respond: testValidateJSONRequestAndRespondOK(createVIPPoolTestResponse, &NewVIPPoolConfiguration{}, func(test *testing.T, requestBody interface{}) {
			verifyCreateVIPPoolTestRequest(test, requestBody.(*NewVIPPoolConfiguration))
		}),
	})
}

/*
 * Test requests.
 */

var createVIPPoolTestRequest = `
{
	"networkDomainId": "553f26b6-2a73-42c3-a78b-6116f11291d0",
	"name": "myDevelopmentPool.1",
	"description": "Pool for load balancing development application servers.",
	"loadBalanceMethod": "ROUND_ROBIN",
	"healthMonitorId": [
		"01683574-d487-11e4-811f-005056806999",
		"0168546c-d487-11e4-811f-005056806999"
	],
	"serviceDownAction": "RESELECT",
	"slowRampTime": 10
}
`

func verifyCreateVIPPoolTestRequest(test *testing.T, request *NewVIPPoolConfiguration) {
	expect := expect(test)

	expect.NotNil("NewVIPPoolConfiguration", request)
	expect.EqualsString("NewVIPPoolConfiguration.NetworkDomainID", "553f26b6-2a73-42c3-a78b-6116f11291d0", request.NetworkDomainID)
	expect.EqualsString("NewVIPPoolConfiguration.Name", "myDevelopmentPool.1", request.Name)
	expect.EqualsString("NewVIPPoolConfiguration.Description", "Pool for load balancing development application servers.", request.Description)
	expect.EqualsString("NewVIPPoolConfiguration.LoadBalanceMethod", LoadBalanceMethodRoundRobin, request.LoadBalanceMethod)
	expect.EqualsInt("NewVIPPoolConfiguration.HealthMonitorIDs.Length", 2, len(request.HealthMonitorIDs))
	expect.EqualsString("NewVIPPoolConfiguration.HealthMonitorIDs[0]", "01683574-d487-11e4-811f-005056806999", request.HealthMonitorIDs[0])
	expect.EqualsString("NewVIPPoolConfiguration.HealthMonitorIDs[1]", "0168546c-d487-11e4-811f-005056806999", request.HealthMonitorIDs[1])
	expect.EqualsString("NewVIPPoolConfiguration.ServiceDownAction", ServiceDownActionReselect, request.ServiceDownAction)
	expect.EqualsInt("NewVIPPoolConfiguration.SlowRampTime", 10, request.SlowRampTime)
}

/*
 * Test responses.
 */

var createVIPPoolTestResponse = `
{
	"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad",
	"operation": "CREATE_POOL",
	"responseCode": "OK",
	"message": "Pool 'myDevelopmentPool.1' has been created.",
	"info": [
		{
			"name": "poolId",
			"value": "4d360b1f-bc2c-4ab7-9884-1f03ba2768f7"
		},
		{
			"name": "name",
			"value": "myDevelopmentPool.1"
		}
	]
}
 `

func verifyDeployVIPPoolTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "CREATE_POOL", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResponseCodeOK, response.ResponseCode)
	expect.EqualsString("Response.Message", "Pool 'myDevelopmentPool.1' has been created.", response.Message)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
	expect.EqualsInt("Response.Message.Length", 2, len(response.Message))
	expect.EqualsString("Response.FieldMessages[0].FieldName", "poolId", response.FieldMessages[0].FieldName)
	expect.EqualsString("Response.FieldMessages[0].Message", "4d360b1f-bc2c-4ab7-9884-1f03ba2768f7", response.FieldMessages[0].Message)
	expect.EqualsString("Response.FieldMessages[1].FieldName", "name", response.FieldMessages[1].FieldName)
	expect.EqualsString("Response.FieldMessages[1].Message", "myDevelopmentPool.1", response.FieldMessages[1].Message)
}
