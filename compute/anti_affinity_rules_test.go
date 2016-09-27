package compute

import "testing"

// List anti-affinity rules (successful).
func TestClient_ListAntityAffinityRules_Success(test *testing.T) {
	expect := expect(test)

	testClientRequest(test, &ClientTestConfig{
		Request: func(test *testing.T, client *Client) {
			page := DefaultPaging()
			rules, err := client.ListServerAntiAffinityRules("553f26b6-2a73-42c3-a78b-6116f11291d0", page)
			if err != nil {
				test.Fatal(err)
			}

			expect.NotNil("ServerAntiAffinityRules", rules)
			expect.EqualsInt("ServerAntiAffinityRules.Length", 1, len(rules.Items))

			rule := rules.Items[0]
			expect.EqualsString("ServerAntiAffinityRules[0].State", "NORMAL", rule.State)
			expect.EqualsString("ServerAntiAffinityRules[0].DatacenterID", "NA9", rule.DatacenterID)

			expect.EqualsInt("ServerAntiAffinityRules[0].Servers.Length", 2, len(rule.Servers))
			expect.EqualsString("ServerAntiAffinityRules[0].Servers[0].ID", "681a6db2-9c7c-4d98-a0c4-7b3d7c1619ba", rule.Servers[0].ID)
			expect.EqualsString("ServerAntiAffinityRules[0].Servers[1].ID", "5783e93f-5370-44fc-a772-cd3c29a2ecaa", rule.Servers[1].ID)
		},
		Respond: testRespondOK(listServerAntiAffinityRulesTestResponse),
	})
}

// Create anti-affinity rule (successful).
func TestClient_CreateAntityAffinityRule_Success(test *testing.T) {
	expect := expect(test)

	testClientRequest(test, &ClientTestConfig{
		Request: func(test *testing.T, client *Client) {
			ruleID, err := client.CreateServerAntiAffinityRule(
				"40285f24-300f-11e2-b574-1a6dd6e90d84",
				"00616730-faca-4cb7-860d-07c553f4c41e",
			)
			if err != nil {
				test.Fatal(err)
			}

			expect.EqualsString("RuleID", "20ce6bee-a4ed-11e1-a91c-0030487e0302", ruleID)
		},
		Respond: testValidateXMLRequestAndRespondOK(createServerAntiAffinityRuleTestResponse, &newServerAntiAffinityRule{}, func(test *testing.T, requestBody interface{}) {
			verifyCreateServerAntiAffinityTestRequest(test, requestBody.(*newServerAntiAffinityRule))
		}),
	})
}

/*
 * Test requests.
 */

const createServerAntiAffinityRuleTestRequest = `
<NewAntiAffinityRule xmlns="http://oec.api.opsource.net/schemas/server">
	<serverId>40285f24-300f-11e2-b574-1a6dd6e90d84</serverId>
	<serverId>00616730-faca-4cb7-860d-07c553f4c41e</serverId>
</NewAntiAffinityRule>
`

func verifyCreateServerAntiAffinityTestRequest(test *testing.T, request *newServerAntiAffinityRule) {
	expect := expect(test)

	expect.NotNil("NewServerAntiAffinityRule", request)

	expect.EqualsInt("NewServerAntiAffinityRule.ServerIDs.Length", 2, len(request.ServerIds))
	expect.EqualsString("NewServerAntiAffinityRule.ServerIDs[0]", "40285f24-300f-11e2-b574-1a6dd6e90d84", request.ServerIds[0])
	expect.EqualsString("NewServerAntiAffinityRule.ServerIDs[1]", "00616730-faca-4cb7-860d-07c553f4c41e", request.ServerIds[1])
}

/*
 * Test responses.
 */

const listServerAntiAffinityRulesTestResponse = `
{
    "antiAffinityRule": [
        {
            "serverSummary": [
                {
                    "name": "Production Server 1",
                    "description": "",
                    "networkingDetails": {
                        "networkInfo": {
                            "primaryNic": {
                                "id": "a6a16a86-7e5b-4138-8c94-3c09f5195a98",
                                "privateIpv4": "10.0.0.13",
                                "ipv6": "2607:f480:1111:1348:5909:96d3:29f5:5e4d",
                                "vlanId": "a6d3e2d7-0092-4f87-b00c-f276127bc26d",
                                "vlanName": "Main VLAN"
                            },
                            "additionalNic": [
                                {
                                    "id": "1e35ea03-a8e5-4771-abc0-8269a45c5735",
                                    "privateIpv4": "10.0.3.12",
                                    "ipv6": "2607:f480:1111:1351:73e4:7d49:93f1:31a3",
                                    "vlanId": "8f19ad6d-ebbb-4393-9ef3-060e5f0b0618",
                                    "vlanName": "Secondary VLAN"
                                }
                            ],
                            "networkDomainId": "553f26b6-2a73-42c3-a78b-6116f11291d0",
                            "networkDomainName": "Production Network Domani"
                        }
                    },
                    "id": "681a6db2-9c7c-4d98-a0c4-7b3d7c1619ba"
                },
                {
                    "name": "Production Server 2",
                    "description": "",
                    "networkingDetails": {
                        "networkInfo": {
                            "primaryNic": {
                                "id": "b3b9261e-50d8-4919-bbe6-866b65b223a1",
                                "privateIpv4": "10.0.3.13",
                                "ipv6": "2607:f480:1111:1351:5793:ebd8:d01c:53f3",
                                "vlanId": "8f19ad6d-ebbb-4393-9ef3-060e5f0b0618",
                                "vlanName": "Production VLAN"
                            },
                            "additionalNic": [
                                {
                                    "id": "c9975e8c-77e8-4c0c-8aa1-17df726c37cc",
                                    "privateIpv4": "10.0.2.13",
                                    "ipv6": "2607:f480:1111:1350:c7f:bef6:bc89:f0e4",
                                    "vlanId": "5404a06d-a084-4f15-a18a-9dfea006d00c",
                                    "vlanName": "Bus VLAN"
                                }
                            ],
                            "networkDomainId": "553f26b6-2a73-42c3-a78b-6116f11291d0",
                            "networkDomainName": "Production Network Domain"
                        }
                    },
                    "id": "5783e93f-5370-44fc-a772-cd3c29a2ecaa"
                }
            ],
            "id": "d4ebfdd1-ec03-45c7-b0be-fbcc0861e9bf",
            "state": "NORMAL",
            "created": "2015-06-05T14:44:54.000Z",
            "datacenterId": "NA9"
        }
    ],
    "pageNumber": 1,
    "pageCount": 1,
    "totalCount": 1,
    "pageSize": 250
}
`

const createServerAntiAffinityRuleTestResponse = `
<Status>
	<operation>Create Anti Affinity Rule</operation>
	<result>SUCCESS</result>
	<resultDetail>Success message</resultDetail>
	<resultCode>RESULT_0</resultCode>
	<additionalInformation name="antiaffinityrule.id">
		<value>20ce6bee-a4ed-11e1-a91c-0030487e0302</value>
	</additionalInformation>
</Status>
`
