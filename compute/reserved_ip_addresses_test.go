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

	server, err := client.ListReservedIPv6AddressesInVLAN("efa6f2fc-9d43-11e7-8991-0389a5a13529")
	if err != nil {
		test.Fatal("Unable to retrieve reserved IPv6 address.", err)
	}

	verifyListReservedIPv6AddressesInVLANTestResponse(test, server)
}

/*
 * Test responses.
 */

//const listReservedPrivateIPv4AddressesInVLANTestResponse = `
//{
//   "ipv4": [
//       {
//           "value": "10.0.0.11",
//           "datacenterId": "NA9",
//           "vlanId": "5d1d62c4-0627-4dc9-83a3-985fbd82ff29"
//       },
//       {
//           "value": "10.0.0.12",
//           "datacenterId": "NA9",
//           "vlanId": "5d1d62c4-0627-4dc9-83a3-985fbd82ff29"
//       }
//   ],
//   "pageNumber": 1,
//   "pageCount": 2,
//   "totalCount": 2,
//   "pageSize": 250
//}`

const listReservedPrivateIPv4AddressesInVLANTestResponse = `
{
	"ipv4": [
	{
		"datacenterid": "NA1",
		"networkId": "8e4515fa-9d54-11e7-8991-0389a5a13529", 
		"exclusive": "true",
		"ipAddress": "10.0.0.20",
		"description": "this is an exclusively reserved IPv4 address",
		"vlanId": "5d1d62c4-0627-4dc9-83a3-985fbd82ff29"
	}, 
	{
		"datacenterid": "NA1",
		"networkId": "80e94d18-9d54-11e7-8991-0389a5a13529", 
		"exclusive": "false",
		"ipAddress": "10.0.0.21",
		"vlanId": "5d1d62c4-0627-4dc9-83a3-985fbd82ff29"
	}, 
    {
		"datacenterid": "NA1",
		"networkId": "783aa7b6-9d54-11e7-8991-0389a5a13529", 
    	"exclusive": "true",
		"ipAddress": "10.0.0.22",
		"description": "This is an exclusively reserved IPv4 address"
	},
	{
		"datacenterid": "NA1",
		"networkId": "87f0566a-9d54-11e7-8991-0389a5a13529", 
		"exclusive": "false",
		"ipAddress": "10.0.0.20"
	} ],
	"pageNumber": 1,
	"pageCount": 4,
	"totalCount": 4,
	"pageSize": 250
}`

func verifyListReservedPrivateIPv4AddressesInVLANTestResponse(test *testing.T, reservedIPv4Addresses *ReservedIPv4Addresses) {
	expect := expect(test)

	expect.NotNil("ReservedPrivateIPv4Addresses", reservedIPv4Addresses)
	expect.EqualsInt("ReservedPrivateIPv4Addresses.PageCount", 4, reservedIPv4Addresses.PageCount)
	expect.EqualsInt("ReservedPrivateIPv4Addresses.PageSize", 250, reservedIPv4Addresses.PageSize)
	expect.EqualsInt("ReservedPrivateIPv4Addresses.TotalCount", 4, reservedIPv4Addresses.TotalCount)
	expect.EqualsInt("ReservedPrivateIPv4Addresses.Length", 4, len(reservedIPv4Addresses.Items))

	address1 := reservedIPv4Addresses.Items[0]
	expect.EqualsString("ReservedPrivateIPv4Addresses.Items[0].IPAddress", "10.0.0.20", address1.IPAddress)
	expect.EqualsString("ReservedPrivateIPv4Addresses.Items[0].VLANID", "5d1d62c4-0627-4dc9-83a3-985fbd82ff29", address1.VLANID)
	expect.EqualsString("ReservedPrivateIPv4Addresses.Items[0].Description",
		"this is an exclusively reserved IPv4 address", address1.Description)

	address2 := reservedIPv4Addresses.Items[1]
	expect.EqualsString("ReservedPrivateIPv4Addresses.Items[1].IPAddress", "10.0.0.21", address2.IPAddress)
	expect.EqualsString("ReservedPrivateIPv4Addresses.Items[1].VLANID", "5d1d62c4-0627-4dc9-83a3-985fbd82ff29", address2.VLANID)
	expect.EqualsString("ReservedPrivateIPv4Addresses.Items[1].Description",
		"", address2.Description)
}

const listReservedIPv6AddressesInVLANTestResponse = `
{
    "pageNumber": 1,
    "pageCount": 3,
    "totalCount": 3,
    "pageSize": 250,
    "reservedIpv6Address": [
		{
			"datacenterid": "NA1",
			"vlanId": "efa6f2fc-9d43-11e7-8991-0389a5a13529", 
			"exclusive": "true",
			"ipAddress": "2001:cdba:0000:0000:0000:0000:3257:9652", 
			"description": "this is an exclusively reserved IPV6 address"
		}, 
		{
			"datacenterid": "NA1",
			"vlanId": "f8c54e10-9d43-11e7-8991-0389a5a13529", 
			"exclusive": "false",
			"ipAddress": "2607:f0d0:1002:0051:0000:0000:0000:0004"
		}, 
		{
			"datacenterid": "NA1",
			"vlanId": "fd65c4ea-9d43-11e7-8991-0389a5a13529", 
			"exclusive": "true",
			"ipAddress": "2001:0000:3238:DFE1:0063:0000:0000:FEFB", 
			"description": "this is an exclusively reserved IPV6 address"
		} 
	]
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
	expect.EqualsString("ReservedIPv6Addresses.Items[0].VLANID", "efa6f2fc-9d43-11e7-8991-0389a5a13529", address1.VLANID)
	expect.EqualsString("ReservedIPv6Addresses.Items[0].Description", "this is an exclusively reserved IPV6 address", address1.Description)

	address2 := reservedIPv6Addresses.Items[1]
	expect.EqualsString("ReservedIPv6Addresses.Items[1].IPAddress", "2607:f0d0:1002:0051:0000:0000:0000:0004", address2.IPAddress)
	expect.EqualsString("ReservedIPv6Addresses.Items[1].VLANID", "f8c54e10-9d43-11e7-8991-0389a5a13529", address2.VLANID)
	expect.EqualsString("ReservedIPv6Addresses.Items[1].Description", "", address2.Description)

	address3 := reservedIPv6Addresses.Items[2]
	expect.EqualsString("ReservedIPv6Addresses.Items[2].IPAddress", "2001:0000:3238:DFE1:0063:0000:0000:FEFB", address3.IPAddress)
	expect.EqualsString("ReservedIPv6Addresses.Items[2].VLANID", "fd65c4ea-9d43-11e7-8991-0389a5a13529", address3.VLANID)
	expect.EqualsString("ReservedIPv6Addresses.Items[2].Description",
		"this is an exclusively reserved IPV6 address", address3.Description)
}
