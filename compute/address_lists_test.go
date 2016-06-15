package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Get IP address list by Id (successful).
func TestClient_GetIPAddressList_ById_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, getIPAddressListTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	server, err := client.GetIPAddressList("5a32d6e4-9707-4813-a269-56ab4d989f4d")
	if err != nil {
		test.Fatal("Unable to deploy server: ", err)
	}

	verifyGetIPAddressListTestResponse(test, server)
}

/*
 * Test responses.
 */

const getIPAddressListTestResponse = `
	{
		"id": "c8c92ea3-2da8-4d51-8153-f39bec794d69",
		"name": "ProductionIPAddressList",
		"description": "For our production web servers",
		"ipVersion": "IPV4",
		"ipAddress": [
			{
				"begin": "1.1.1.1",
				"end": "2.2.2.2"
			},
			{
				"begin": "192.168.1.1"
			},
			{
				"begin": "192.168.1.1",
				"prefixSize": 24
			}
		],
		"childIpAddressList": [
			{
				"id": "c8c92ea3-2da8-4d51-8153-f39bec794d68",
				"name": "tomcatIpAddresses"
			},
			{
				"id": "c8c92ea3-2da8-4d51-8153-f39bec794d67",
				"name": "mySqlIpAddresses"
			}
		],
		"state": "NORMAL",
		"createTime": "2015-09-29T02:49:45"
	}
`

func verifyGetIPAddressListTestResponse(test *testing.T, addressList *IPAddressList) {
	expect := expect(test)

	expect.NotNil("IPAddressList", addressList)
	expect.EqualsString("IPAddressList.Name", "ProductionIPAddressList", addressList.Name)
	expect.EqualsString("IPAddressList.Description", "For our production web servers", addressList.Description)
	expect.EqualsString("IPAddressList.IPVersion", "IPV4", addressList.IPVersion)
	expect.EqualsString("IPAddressList.Name", "ProductionIPAddressList", addressList.Name)
	expect.EqualsString("IPAddressList.Name", "ProductionIPAddressList", addressList.Name)
	expect.EqualsString("IPAddressList.State", ResourceStatusNormal, addressList.State)
	expect.EqualsString("IPAddressList.CreateTime", "2015-09-29T02:49:45", addressList.CreateTime)

	expect.EqualsInt("IPAddressList.Addresses.Length", 3, len(addressList.Addresses))

	address1 := addressList.Addresses[0]
	expect.EqualsString("IPAddressList.Addresses[0].Begin", "1.1.1.1", address1.Begin)
	expect.NotNil("IPAddressList.Addresses[0].End", address1.End)
	expect.EqualsString("IPAddressList.Addresses[0].End", "2.2.2.2", *address1.End)
	expect.IsNil("IPAddressList.Addresses[0].PrefixSize", address1.PrefixSize)

	address2 := addressList.Addresses[1]
	expect.EqualsString("IPAddressList.Addresses[1].Begin", "192.168.1.1", address2.Begin)
	expect.IsNil("IPAddressList.Addresses[1].End", address2.End)
	expect.IsNil("IPAddressList.Addresses[1].PrefixSize", address2.PrefixSize)

	address3 := addressList.Addresses[2]
	expect.EqualsString("IPAddressList.Addresses[2].Begin", "192.168.1.1", address3.Begin)
	expect.IsNil("IPAddressList.Addresses[2].End", address3.End)
	expect.NotNil("IPAddressList.Addresses[2].PrefixSize", address3.PrefixSize)
	expect.EqualsInt("IPAddressList.Addresses[2].PrefixSize", 24, *address3.PrefixSize)
}
