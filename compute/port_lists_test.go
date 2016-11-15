package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Get port list by Id (successful).
func TestClient_GetPortList_ById_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, getPortListTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	server, err := client.GetPortList("c8c92ea3-2da8-4d51-8153-f39bec794d69")
	if err != nil {
		test.Fatal("Unable to retrieve port list: ", err)
	}

	verifyGetPortListTestResponse(test, server)
}

/*
 * Test responses.
 */

const getPortListTestResponse = `
{
	"id": "c8c92ea3-2da8-4d51-8153-f39bec794d69",
	"name": "MyPortList",
	"description": "Production Servers",
	"port": [
		{
			"begin": 8000,
			"end": 9600
		},
		{
		    "begin": 25
		},
		{
			"begin": 443
		}
	],
	"childPortList": [
		{
			"id": "c8c92ea3-2da8-4d51-8153-f39bec794d68",
			"name": "tomcatPorts"
		},
		{
			"id": "c8c92ea3-2da8-4d51-8153-f39bec794d67",
			"name": "mySqlPorts"
		}
	],
	"state": "NORMAL",
	"createTime": "2008-09-29T02:49:45"
}
`

func verifyGetPortListTestResponse(test *testing.T, portList *PortList) {
	expect := expect(test)

	expect.NotNil("PortList", portList)
	expect.EqualsString("PortList.ID", "c8c92ea3-2da8-4d51-8153-f39bec794d69", portList.ID)
	expect.EqualsString("PortList.Name", "MyPortList", portList.Name)
	expect.EqualsString("PortList.Description", "Production Servers", portList.Description)
	expect.EqualsString("PortList.State", ResourceStatusNormal, portList.State)
	expect.EqualsString("PortList.CreateTime", "2008-09-29T02:49:45", portList.CreateTime)

	expect.EqualsInt("PortList.Ports.Length", 3, len(portList.Ports))

	port1 := portList.Ports[0]
	expect.EqualsInt("PortList.Ports[0].Begin", 8000, port1.Begin)
	expect.NotNil("PortList.Ports[0].End", port1.End)
	expect.EqualsInt("PortList.Ports[0].End", 9600, *port1.End)

	port2 := portList.Ports[1]
	expect.EqualsInt("PortList.Ports[1].Begin", 25, port2.Begin)
	expect.IsNil("PortList.Ports[1].End", port2.End)

	port3 := portList.Ports[2]
	expect.EqualsInt("PortList.Ports[2].Begin", 443, port3.Begin)
	expect.IsNil("PortList.Ports[2].End", port3.End)

	expect.EqualsInt("PortList.ChildLists.Length", 2, len(portList.ChildLists))

	childList1 := portList.ChildLists[0]
	expect.EqualsString("PortList.ChildLists[0].ID", "c8c92ea3-2da8-4d51-8153-f39bec794d68", childList1.ID)
	expect.EqualsString("PortList.ChildLists[0].Name", "tomcatPorts", childList1.Name)

	childList2 := portList.ChildLists[1]
	expect.EqualsString("PortList.ChildLists[1].ID", "c8c92ea3-2da8-4d51-8153-f39bec794d67", childList2.ID)
	expect.EqualsString("PortList.ChildLists[1].Name", "mySqlPorts", childList2.Name)
}
