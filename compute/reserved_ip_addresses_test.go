package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// List reserved private IPv4 addresses in VLAN (successful).
func TestClient_ListReservedPrivateIPv4AddressesInVLAN_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, listReservedPrivateIPv4AddressesInVLANTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	server, err := client.ListReservedPrivateIPv4AddressesInVLAN("c8c92ea3-2da8-4d51-8153-f39bec794d69")
	if err != nil {
		test.Fatal("Unable to retrieve reserved IPv4 address: ", err)
	}

	verifyListReservedPrivateIPv4AddressesInVLANTestResponse(test, server)
}

// List reserved IPv6 addresses in VLAN (successful).
func TestClient_ListReservedIPv6AddressesInVLAN_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, listReservedIPv6AddressesInVLANTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	server, err := client.ListReservedIPv6AddressesInVLAN("c8c92ea3-2da8-4d51-8153-f39bec794d69")
	if err != nil {
		test.Fatal("Unable to retrieve reserved IPv6 address: ", err)
	}

	verifyListReservedIPv6AddressesInVLANTestResponse(test, server)
}

/*
 * Test responses.
 */

const listReservedPrivateIPv4AddressesInVLANTestResponse = `
{
    "ipv4": [
        {
            "value": "10.0.0.11",
            "datacenterId": "NA9",
            "vlanId": "5d1d62c4-0627-4dc9-83a3-985fbd82ff29"
        },
        {
            "value": "10.0.0.12",
            "datacenterId": "NA9",
            "vlanId": "5d1d62c4-0627-4dc9-83a3-985fbd82ff29"
        }
    ],
    "pageNumber": 1,
    "pageCount": 2,
    "totalCount": 2,
    "pageSize": 250
}`

func verifyListReservedPrivateIPv4AddressesInVLANTestResponse(test *testing.T, reservedIPv4Addresses *ReservedIPv4Addresses) {
	expect := expect(test)

	expect.NotNil("ReservedPrivateIPv4Addresses", reservedIPv4Addresses)
	expect.EqualsInt("ReservedPrivateIPv4Addresses.PageCount", 2, reservedIPv4Addresses.PageCount)
	expect.EqualsInt("ReservedPrivateIPv4Addresses.PageSize", 250, reservedIPv4Addresses.PageSize)
	expect.EqualsInt("ReservedPrivateIPv4Addresses.TotalCount", 2, reservedIPv4Addresses.TotalCount)
	expect.EqualsInt("ReservedPrivateIPv4Addresses.Length", 2, len(reservedIPv4Addresses.Items))

	address1 := reservedIPv4Addresses.Items[0]
	expect.EqualsString("ReservedPrivateIPv4Addresses.Items[0].IPAddress", "10.0.0.11", address1.IPAddress)
	expect.EqualsString("ReservedPrivateIPv4Addresses.Items[0].VLANID", "5d1d62c4-0627-4dc9-83a3-985fbd82ff29", address1.VLANID)
	expect.EqualsString("ReservedPrivateIPv4Addresses.Items[0].DatacenterID", "NA9", address1.DatacenterID)

	address2 := reservedIPv4Addresses.Items[1]
	expect.EqualsString("ReservedPrivateIPv4Addresses.Items[1].IPAddress", "10.0.0.12", address2.IPAddress)
	expect.EqualsString("ReservedPrivateIPv4Addresses.Items[1].VLANID", "5d1d62c4-0627-4dc9-83a3-985fbd82ff29", address2.VLANID)
	expect.EqualsString("ReservedPrivateIPv4Addresses.Items[1].DatacenterID", "NA9", address2.DatacenterID)
}

const listReservedIPv6AddressesInVLANTestResponse = `
{
    "reservedIpv6Address": [
        {
            "datacenterId": "NA1",
            "vlanId": "f79b5d92-6594-4659-9776-86b5264130a4",
            "value": "2001:cdba:0000:0000:0000:0000:3257:9652"
        },
        {
            "datacenterId": "NA1",
            "vlanId": "f79b5d92-6594-4659-9776-86b5264130a4",
            "value": "2607:f0d0:1002:0051:0000:0000:0000:0004"
        },
        {
            "datacenterId": "NA1",
            "vlanId": "f79b5d92-6594-4659-9776-86b5264130a4",
            "value": "2001:0000:3238:DFE1:0063:0000:0000:FEFB"
        }
    ],
    "pageNumber": 1,
    "pageCount": 3,
    "totalCount": 3,
    "pageSize": 250
}`

func verifyListReservedIPv6AddressesInVLANTestResponse(test *testing.T, reservedIPv6Addresses *ReservedIPv6Addresses) {
	expect := expect(test)

	expect.NotNil("ReservedIPv6Addresses", reservedIPv6Addresses)
	expect.EqualsInt("ReservedIPv6Addresses.PageCount", 3, reservedIPv6Addresses.PageCount)
	expect.EqualsInt("ReservedIPv6Addresses.PageSize", 250, reservedIPv6Addresses.PageSize)
	expect.EqualsInt("ReservedIPv6Addresses.TotalCount", 3, reservedIPv6Addresses.TotalCount)
	expect.EqualsInt("ReservedIPv6Addresses.Length", 3, len(reservedIPv6Addresses.Items))

	address1 := reservedIPv6Addresses.Items[0]
	expect.EqualsString("ReservedIPv6Addresses.Items[0].IPAddress", "2001:cdba:0000:0000:0000:0000:3257:9652", address1.IPAddress)
	expect.EqualsString("ReservedIPv6Addresses.Items[0].VLANID", "f79b5d92-6594-4659-9776-86b5264130a4", address1.VLANID)
	expect.EqualsString("ReservedIPv6Addresses.Items[0].DatacenterID", "NA1", address1.DatacenterID)

	address2 := reservedIPv6Addresses.Items[1]
	expect.EqualsString("ReservedIPv6Addresses.Items[1].IPAddress", "2607:f0d0:1002:0051:0000:0000:0000:0004", address2.IPAddress)
	expect.EqualsString("ReservedIPv6Addresses.Items[1].VLANID", "f79b5d92-6594-4659-9776-86b5264130a4", address2.VLANID)
	expect.EqualsString("ReservedIPv6Addresses.Items[1].DatacenterID", "NA1", address2.DatacenterID)

	address3 := reservedIPv6Addresses.Items[2]
	expect.EqualsString("ReservedIPv6Addresses.Items[2].IPAddress", "2001:0000:3238:DFE1:0063:0000:0000:FEFB", address3.IPAddress)
	expect.EqualsString("ReservedIPv6Addresses.Items[2].VLANID", "f79b5d92-6594-4659-9776-86b5264130a4", address3.VLANID)
	expect.EqualsString("ReservedIPv6Addresses.Items[2].DatacenterID", "NA1", address3.DatacenterID)
}
