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

	networkDomains, err := client.ListReservedPublicIPAddresses("802abc9f-45a7-4efb-9d5a-810082368708")
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

	expect.notNil("PublicIPBlock", block)

	expect.equalsString("PublicIPBlock.ID", "cacc028a-7f12-11e4-a91c-0030487e0302", block.ID)
	expect.equalsString("PublicIPBlock.NetworkDomainID", "802abc9f-45a7-4efb-9d5a-810082368708", block.NetworkDomainID)
}

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

	expect.notNil("ReservedPublicIPs", reservedIPs)

	expect.equalsInt("ReservedPublicIPs.PageCount", 1, reservedIPs.PageCount)
	expect.equalsInt("ReservedPublicIPs.IPs size", 1, len(reservedIPs.IPs))

	ip1 := reservedIPs.IPs[0]
	expect.equalsString("ReservedPublicIPs.IPs[0].Address", "165.180.12.12", ip1.Address)
	expect.equalsString("ReservedPublicIPs.IPs[0].IPBlockID", "cacc028a-7f12-11e4-a91c-0030487e0302", ip1.IPBlockID)
	expect.equalsString("ReservedPublicIPs.IPs[0].NetworkDomainID", "802abc9f-45a7-4efb-9d5a-810082368708", ip1.NetworkDomainID)
	expect.equalsString("ReservedPublicIPs.IPs[0].DataCenterID", "NA9", ip1.DataCenterID)
}
