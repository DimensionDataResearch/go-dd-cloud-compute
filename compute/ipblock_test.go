package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// List reserved public IPv4 addresses (successful).
func TestClient_ListReservedPublicIPAddresses_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, listReservedPublicIPAddressesTestResponse)
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

	verifyListReservedPublicIPAddressesTestResponse(test, networkDomains)
}

/*
 * Test responses.
 */

const listReservedPublicIPAddressesTestResponse = `
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

func verifyListReservedPublicIPAddressesTestResponse(test *testing.T, reservedIPs *ReservedPublicIPs) {
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
