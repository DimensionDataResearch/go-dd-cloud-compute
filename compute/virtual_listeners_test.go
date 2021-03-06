package compute

import (
	"net/http"
	"testing"
)

// Get virtual listener by Id (successful).
func TestClient_GetVirtualListener_ById_Success(test *testing.T) {
	testClientRequest(test, &ClientTestConfig{
		Request: func(test *testing.T, client *Client) {
			virtualListener, err := client.GetVirtualListener("6115469d-a8bb-445b-bb23-d23b5283f2b9")
			if err != nil {
				test.Fatal(err)
			}

			verifyGetVirtualListenerTestResponse(test, virtualListener)
		},
		Respond: testRespondOK(getVirtualListenerTestResponse),
	})
}

// Create virtual listener (successful).
func TestClient_CreateVirtualListener_Success(test *testing.T) {
	testClientRequest(test, &ClientTestConfig{
		Request: func(test *testing.T, client *Client) {
			virtualListenerID, err := client.CreateVirtualListener(NewVirtualListenerConfiguration{
				NetworkDomainID:              "553f26b6-2a73-42c3-a78b-6116f11291d0",
				Name:                         "Production.Load.Balancer",
				Description:                  "Used as the load balancer for the production applications.",
				Type:                         VirtualListenerTypeStandard,
				Protocol:                     VirtualListenerStandardProtocolHTTP,
				ListenerIPAddress:            stringToPtr("165.180.12.22"),
				Port:                         80,
				Enabled:                      true,
				ConnectionLimit:              25000,
				ConnectionRateLimit:          2000,
				SourcePortPreservation:       SourcePortPreservationEnabled,
				PoolID:                       stringToPtr("afb1fb1a-eab9-43f4-95c2-36a4cdda6cb8"),
				ClientClonePoolID:            stringToPtr("033a97dc-ee9b-4808-97ea-50b06624fd16"),
				PersistenceProfileID:         stringToPtr("a34ca25c-f3db-11e4-b010-005056806999"),
				FallbackPersistenceProfileID: stringToPtr("6f2f5d7b-cdd9-4d84-8ad7-999b64a87978"),
				SSLOffloadProfileID:          stringToPtr("233a97dc-ee9b-4808-97ea-50b06624fd18"),
				IRuleIDs:                     []string{"2b20abd9-ffdc-11e4-b010-005056806999"},
				OptimizationProfile:          stringToPtr("TCP"),
			})
			if err != nil {
				test.Fatal(err)
			}

			expect(test).EqualsString("VirtuaListenerID", "43a445f1-9ac9-4f13-8b0d-a2d1fad231c3", virtualListenerID)
		},
		Respond: testValidateJSONRequestAndRespond(http.StatusOK, createVirtualListenerTestResponse, &NewVirtualListenerConfiguration{}, func(test *testing.T, requestBody interface{}) {
			verifyCreateVirtualListenerTestRequest(test, requestBody.(*NewVirtualListenerConfiguration))
		}),
	})
}

/*
 * Test requests.
 */

const createVirtualListenerTestRequest = `
{
	"networkDomainId": "553f26b6-2a73-42c3-a78b-6116f11291d0",
	"name": "Production.Load.Balancer",
	"description": "Used as the load balancer for the production applications.",
	"type": "STANDARD",
	"protocol": "TCP",
	"listenerIpAddress": "165.180.12.22",
	"port": 80,
	"enabled": true,
	"connectionLimit": 25000,
	"connectionRateLimit": 2000,
	"sourcePortPreservation": "PRESERVE",
	"poolId": "afb1fb1a-eab9-43f4-95c2-36a4cdda6cb8",
	"clientClonePoolId": "033a97dc-ee9b-4808-97ea-50b06624fd16",
	"persistenceProfileId": "a34ca25c-f3db-11e4-b010-005056806999",
	"fallbackPersistenceProfileId": "6f2f5d7b-cdd9-4d84-8ad7-999b64a87978",
	"sslOffloadProfileId": "233a97dc-ee9b-4808-97ea-50b06624fd18",
	"iruleId": [
		"2b20abd9-ffdc-11e4-b010-005056806999"
	],
	"optimizationProfile": "TCP"
}
`

func verifyCreateVirtualListenerTestRequest(test *testing.T, request *NewVirtualListenerConfiguration) {
	expect := expect(test)

	expect.NotNil("VirtualListener", request)

	expect.EqualsString("VirtualListener.Name", "Production.Load.Balancer", request.Name)
	expect.EqualsString("VirtualListener.Description", "Used as the load balancer for the production applications.", request.Description)
	expect.EqualsString("VirtualListener.Type", VirtualListenerTypeStandard, request.Type)

	expect.EqualsString("VirtualListener.Protocol", VirtualListenerStandardProtocolHTTP, request.Protocol)

	expect.NotNil("VirtualListener.ListenerIPAddress", request.ListenerIPAddress)
	expect.EqualsString("VirtualListener.ListenerIPAddress", "165.180.12.22", *request.ListenerIPAddress)

	expect.EqualsInt("VirtualListener.Port", 80, request.Port)

	expect.IsTrue("VirtualListener.Enabled", request.Enabled)

	expect.EqualsInt("VirtualListener.ConnectionLimit", 25000, request.ConnectionLimit)
	expect.EqualsInt("VirtualListener.ConnectionRateLimit", 2000, request.ConnectionRateLimit)

	expect.EqualsString("VirtualListener.SourcePortPreservation", SourcePortPreservationEnabled, request.SourcePortPreservation)

	expect.NotNil("VirtualListener.PersistenceProfileID", request.PersistenceProfileID)
	expect.EqualsString("VirtualListener.PersistenceProfileID", "a34ca25c-f3db-11e4-b010-005056806999", *request.PersistenceProfileID)

	expect.NotNil("VirtualListener.FallbackPersistenceProfileID", request.FallbackPersistenceProfileID)
	expect.EqualsString("VirtualListener.FallbackPersistenceProfileID", "6f2f5d7b-cdd9-4d84-8ad7-999b64a87978", *request.FallbackPersistenceProfileID)

	expect.NotNil("VirtualListener.SSLOffloadProfileID", request.SSLOffloadProfileID)
	expect.EqualsString("VirtualListener.SSLOffloadProfileID", "233a97dc-ee9b-4808-97ea-50b06624fd18", *request.SSLOffloadProfileID)

	expect.EqualsString("VirtualListener.NetworkDomainID", "553f26b6-2a73-42c3-a78b-6116f11291d0", request.NetworkDomainID)
}

/*
 * Test responses.
 */

const createVirtualListenerTestResponse = `
{
    "operation": "CREATE_VIRTUAL_LISTENER",
    "responseCode": "OK",
    "message": "Virtual Listener 'Production.Load.Balancer' has been created on Public IP Address 165.180.12.22.",
    "info": [{
            "name": "virtualListenerId",
            "value": "43a445f1-9ac9-4f13-8b0d-a2d1fad231c3"
        },
        {
            "name": "name",
            "value": "Production.Load.Balancer"
        },
        {
            "name": "listenerIpAddress",
            "value": "165.180.12.22"
        }
    ],
    "warning": [],
    "error": [],
    "requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
}
`

const editVirtualListenerTestResponse = `
{
	"id": "6e42868b-e013-41c3-ac38-5f7b50d54808",
	"description": "Used as the load balancer for the production applications.",
	"enabled": true,
	"connectionLimit": 25000,
	"connectionRateLimit": 2000,
	"sourcePortPreservation": "PRESERVE",
	"poolId": "afb1fb1a-eab9-43f4-95c2-36a4cdda6cb8",
	"persistenceProfileId": "a34ca25c-f3db-11e4-b010-005056806999",
	"iruleId": [
		"2b20abd9-ffdc-11e4-b010-005056806999"
	],
	"optimizationProfile": [
		"TCP"
	]
}
`

const getVirtualListenerTestResponse = `
{
	"id": "6115469d-a8bb-445b-bb23-d23b5283f2b9",
    "name": "myProduction.Virtual.Listener",
    "description": "Virtual Listener for load balancing our test systems.",
    "type": "PERFORMANCE_LAYER_4",
    "protocol": "HTTP",
    "listenerIpAddress": "165.180.12.22",
    "port": 8899,
    "enabled": true,
    "connectionLimit": 10000,
    "connectionRateLimit": 400,
    "sourcePortPreservation": "PRESERVE",
    "pool": {
        "loadBalanceMethod": "ROUND_ROBIN",
        "serviceDownAction": "NONE",
        "slowRampTime": 10,
        "healthMonitor": [
            {
                "id": "01683574-d487-11e4-811f-005056806999",
                "name": "CCDEFAULT.Http"
            },
            {
                "id": "0168546c-d487-11e4-811f-005056806999",
                "name": "CCDEFAULT.Https"
            }
        ],
        "id": "afb1fb1a-eab9-43f4-95c2-36a4cdda6cb8",
        "name": "myProductionPool.1"
    },
    "clientClonePool": {
        "loadBalanceMethod": "ROUND_ROBIN",
        "serviceDownAction": "RESELECT",
        "slowRampTime": 10,
        "healthMonitor": [
            {
                "id": "01683574-d487-11e4-811f-005056806999",
                "name": "CCDEFAULT.Http"
            },
            {
                "id": "0168546c-d487-11e4-811f-005056806999",
                "name": "CCDEFAULT.Https"
            }
        ],
        "id": "6f2f5d7b-cdd9-4d84-8ad7-999b64a87978",
        "name": "myDevelopmentPool.1"
    },
    "persistenceProfile": {
        "id": "a34ca25c-f3db-11e4-b010-005056806999",
        "name": "CCDEFAULT.DestinationAddress"
    },
    "fallbackPersistenceProfile": {
        "id": "a34ca3f6-f3db-11e4-b010-005056806999",
        "name": "CCDEFAULT.SourceAddress"
    },
	"sslOffloadProfile": {
		"id": "9e6b496d-5261-4542-91aa-b50c7f569c55",
		"name": "mySslProfile"
	}, 
    "datacenterId": "NA9",
    "irule": [
        {
            "id": "2b20abd9-ffdc-11e4-b010-005056806999",
            "name": "CCDEFAULT.IpProtocolTimers"
        },
        {
            "id": "2b20e790-ffdc-11e4-b010-005056806999",
            "name": "CCDEFAULT.Ips"
        }
    ],
	"state": "NORMAL",
    "createTime": "2015-05-28T15:59:49.000Z",
	"networkDomainId": "553f26b6-2a73-42c3-a78b-6116f11291d0"
}
`

func verifyGetVirtualListenerTestResponse(test *testing.T, response *VirtualListener) {
	expect := expect(test)

	expect.NotNil("VirtualListener", response)
	expect.EqualsString("VirtualListener.ID", "6115469d-a8bb-445b-bb23-d23b5283f2b9", response.ID)
	expect.EqualsString("VirtualListener.Name", "myProduction.Virtual.Listener", response.Name)
	expect.EqualsString("VirtualListener.Description", "Virtual Listener for load balancing our test systems.", response.Description)
	expect.EqualsString("VirtualListener.Type", "PERFORMANCE_LAYER_4", response.Type)
	expect.EqualsString("VirtualListener.Protocol", "HTTP", response.Protocol)
	expect.EqualsString("VirtualListener.ListenerIPAddress", "165.180.12.22", response.ListenerIPAddress)
	expect.EqualsInt("VirtualListener.Port", 8899, response.Port)
	expect.IsTrue("VirtualListener.Enabled", response.Enabled)
	expect.EqualsInt("VirtualListener.ConnectionLimit", 10000, response.ConnectionLimit)
	expect.EqualsInt("VirtualListener.ConnectionRateLimit", 400, response.ConnectionRateLimit)
	expect.EqualsString("VirtualListener.SSLOffloadProfile.ID", "9e6b496d-5261-4542-91aa-b50c7f569c55", response.SSLOffloadProfile.ID)
	expect.EqualsString("VirtualListener.SSLOffloadProfile.Name", "mySslProfile", response.SSLOffloadProfile.Name)
	expect.EqualsString("VirtualListener.SourcePortPreservation", SourcePortPreservationEnabled, response.SourcePortPreservation)
	expect.EqualsString("VirtualListener.State", ResourceStatusNormal, response.State)
	expect.EqualsString("VirtualListener.NetworkDomainID", "553f26b6-2a73-42c3-a78b-6116f11291d0", response.NetworkDomainID)
}
