package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Get public IPv4 address block by Id (successful).
func TestClient_GetPublicIPBlock_ById_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, getPublicIPBlockResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	block, err := client.GetPublicIPBlock("cacc028a-7f12-11e4-a91c-0030487e0302")
	if err != nil {
		test.Fatal(err)
	}

	verifyGetPublicIPBlockResponse(test, block)
}

// Get public IPv4 address block by Id (successful).
func TestClient_AddPublicIPBlock_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, addPublicIPBlockResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	blockID, err := client.AddPublicIPBlock("484174a2-ae74-4658-9e56-50fc90e086cf")
	if err != nil {
		test.Fatal(err)
	}

	expect(test).EqualsString("PublicIPBlockID", "4487241a-f0ca-11e3-9315-d4bed9b167ba", blockID)
}

// List reserved public IPv4 addresses (successful).
func TestClient_ListReservedPublicIPAddresses_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, listReservedPublicIPAddressesResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	networkDomains, err := client.ListReservedPublicIPAddresses("802abc9f-45a7-4efb-9d5a-810082368708", nil)
	if err != nil {
		test.Fatal(err)
	}

	verifyListReservedPublicIPAddressesResponse(test, networkDomains)
}

/*
 * Test responses.
 */

const getPublicIPBlockResponse = `
	{
		"networkDomainId": "802abc9f-45a7-4efb-9d5a-810082368708",
		"baseIp": "165.180.12.12",
		"size": 2,
		"createTime": "2014-12-15T16:35:07.000Z",
		"state": "NORMAL",
		"id": "cacc028a-7f12-11e4-a91c-0030487e0302",
		"datacenterId": "NA9"
	}
`

func verifyGetPublicIPBlockResponse(test *testing.T, block *PublicIPBlock) {
	expect := expect(test)

	expect.NotNil("PublicIPBlock", block)

	expect.EqualsString("PublicIPBlock.ID", "cacc028a-7f12-11e4-a91c-0030487e0302", block.ID)
	expect.EqualsString("PublicIPBlock.NetworkDomainID", "802abc9f-45a7-4efb-9d5a-810082368708", block.NetworkDomainID)
}

const addPublicIPBlockResponse = `
	{
		"operation": "ADD_PUBLIC_IP_BLOCK",
		"responseCode": "OK",
		"message": "Public IPv4 Address Block has been added successfully to Network Domain '484174a2-ae74-4658-9e56-50fc90e086cf'.",
		"info": [
			{
				"name": "ipBlockId",
				"value": "4487241a-f0ca-11e3-9315-d4bed9b167ba"
			}
		],
		"warning": [],
		"error": [],
		"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
	}
`

const listReservedPublicIPAddressesResponse = `
	{
		"ip": [
			{
				"value": "165.180.12.12",
				"datacenterId": "NA9",
				"ipBlockId": "cacc028a-7f12-11e4-a91c-0030487e0302",
				"networkDomainId": "802abc9f-45a7-4efb-9d5a-810082368708"
			}
		],
		"pageNumber": 1,
		"pageCount": 1,
		"totalCount": 1,
		"pageSize": 250
	}
`

func verifyListReservedPublicIPAddressesResponse(test *testing.T, reservedIPs *ReservedPublicIPs) {
	expect := expect(test)

	expect.NotNil("ReservedPublicIPs", reservedIPs)

	expect.EqualsInt("ReservedPublicIPs.PageCount", 1, reservedIPs.PageCount)
	expect.EqualsInt("ReservedPublicIPs.IPs size", 1, len(reservedIPs.IPs))

	ip1 := reservedIPs.IPs[0]
	expect.EqualsString("ReservedPublicIPs.IPs[0].Address", "165.180.12.12", ip1.Address)
	expect.EqualsString("ReservedPublicIPs.IPs[0].IPBlockID", "cacc028a-7f12-11e4-a91c-0030487e0302", ip1.IPBlockID)
	expect.EqualsString("ReservedPublicIPs.IPs[0].NetworkDomainID", "802abc9f-45a7-4efb-9d5a-810082368708", ip1.NetworkDomainID)
	expect.EqualsString("ReservedPublicIPs.IPs[0].DataCenterID", "NA9", ip1.DataCenterID)
}
