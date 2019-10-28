package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_CreateStaticRoute_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestBody, err := readRequestBodyAsString(request)
		if err != nil {
			test.Fatal("Failed to read request body: ", err)
		}

		expect.EqualsString("Request.Body",
			`{"networkDomainId":"e926545b-1b9c-4068-8cef-076830a9a0bc","name":"ClientStaticRoute","description":"This is a client static route.","ipVersion":"IPV4","destinationNetworkAddress":"10.0.0.0","destinationPrefixSize":24,"nextHopAddress":"10.10.10.1"}`,
			requestBody,
		)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, createStaticRouteResponse)
	}))

	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	staticRouteID, err := client.CreateStaticRoute("e926545b-1b9c-4068-8cef-076830a9a0bc", "ClientStaticRoute",
		"This is a client static route.", "IPV4", "10.0.0.0", 24, "10.10.10.1")

	if err != nil {
		test.Fatal(err)
	}

	expect.EqualsString("staticRouteId", "0e56433f-d808-4669-821d-812769517ff8", staticRouteID)

}

var createStaticRouteResponse = `
{
	"operation": "CREATE_STATIC_ROUTE", "responseCode": "OK",
	"message": "Static Route has been created.", "info": [
			{
				"name": "staticRouteId",
				"value": "0e56433f-d808-4669-821d-812769517ff8"
	} ],
		"warning": [],
		"error": [],
		"requestId": "na9_20180321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
}
`

func TestClient_DeleteStaticRoute_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestBody, err := readRequestBodyAsString(request)
		if err != nil {
			test.Fatal("Failed to read request body: ", err)
		}

		expect.EqualsString("Request.Body",
			`{"id":"0e56433f-d808-4669-821d-812769517ff8"}`,
			requestBody,
		)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, deleteStaticRouteResponse)
	}))

	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	err := client.DeleteStaticRoute("0e56433f-d808-4669-821d-812769517ff8")

	if err != nil {
		test.Fatal(err)
	}

	// Pass
}

var deleteStaticRouteResponse = `
{
    "operation": "DELETE_STATIC_ROUTE",
	"responseCode": "OK",
	"message": "Static Route has been deleted.",
	"requestId": "na/2018-04-14T13:37:20/62f06368-c3fb-11e3-b29c-001517c4643e" 
}
`

func TestClient_ListStaticRoute_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, listStaticRouteTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	staticRoutes, err := client.ListStaticRoute(nil)
	if err != nil {
		test.Fatal(err)
	}

	verifyListStaticRouteTestResponse(test, staticRoutes)
}

var listStaticRouteTestResponse = `
{
    "staticRoute": [
	{
	"id": "9e6b496d-5261-4542-91aa-b50c7f569c54",
	"datacenterId": "AU1",
	"networkDomainId": "9888c372-eb4b-11e3-b29c-001517c4643e", "name": "ClientStaticRoute",
	"description": "This is a Client Static Route",
	"type": "CLIENT",
	"ipVersion": "IPV4",
	"destinationNetworkAddress": "132.15.2.0", "destinationPrefixSize": 24,
	"nextHopAddress": "132.15.3.2",
	"nextHopAddressVlanId": "9888c372-eb4b-11e3-b29c-001517c54433", "state": "NORMAL",
	"createTime": "2018-05-16T12:05:10.000Z"
	} ],
    "pageNumber": 1,
    "pageCount": 1,
    "totalCount": 1,
    "pageSize": 250
}
`

func verifyListStaticRouteTestResponse(test *testing.T, staticRoutes *StaticRoutes) {
	expect := expect(test)

	expect.NotNil("StaticRoute", staticRoutes)
	expect.EqualsInt("StaticRoute.PageCount", 1, staticRoutes.PageCount)
	expect.EqualsInt("StaticRoute.Routes size", 1, len(staticRoutes.Routes))

	route1 := staticRoutes.Routes[0]
	expect.EqualsString("StaticRoutes.Routes[0].Name", "ClientStaticRoute", route1.Name)
	expect.EqualsString("StaticRoutes.Routes[0].Description", "This is a Client Static Route", route1.Description)
	expect.EqualsString("StaticRoutes.Routes[0].IpVersion", "IPV4", route1.IpVersion)
	expect.EqualsString("StaticRoutes.Routes[0].DestinationNetworkAddress", "132.15.2.0", route1.DestinationNetworkAddress)
	expect.EqualsInt("StaticRoutes.Routes[0].DestinationPrefixSize", 24, route1.DestinationPrefixSize)
	expect.EqualsString("StaticRoutes.Routes[0].NextHopAddress", "132.15.3.2", route1.NextHopAddress)
}

func TestClient_GetStaticRoute_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, getStaticRouteTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	staticRoute, err := client.GetStaticRoute("9e6b496d-5261-4542-91aa-b50c7f569c54")
	if err != nil {
		test.Fatal(err)
	}

	verifyGetStaticRouteTestResponse(test, staticRoute)
}

var getStaticRouteTestResponse = `
{
	"id": "9e6b496d-5261-4542-91aa-b50c7f569c54",
	"datacenterId": "AU1",
	"networkDomainId": "9888c372-eb4b-11e3-b29c-001517c4643e", 
	"name": "ClientStaticRoute",
	"description": "This is a Client Static Route",
	"type": "CLIENT",
	"ipVersion": "IPV4",
	"destinationNetworkAddress": "132.15.2.0", "destinationPrefixSize": 24,
	"nextHopAddress": "132.15.3.2",
	"nextHopAddressVlanId": "9888c372-eb4b-11e3-b29c-001517c54433", 
	"state": "NORMAL",
	"createTime": "2018-05-16T12:05:10.000Z"
	}
`

func verifyGetStaticRouteTestResponse(test *testing.T, staticRoute *StaticRoute) {
	expect := expect(test)

	expect.NotNil("StaticRoute", staticRoute)

	expect.EqualsString("StaticRoutes.Routes[0].Name", "ClientStaticRoute", staticRoute.Name)
	expect.EqualsString("StaticRoutes.Routes[0].Description", "This is a Client Static Route", staticRoute.Description)
	expect.EqualsString("StaticRoutes.Routes[0].IpVersion", "IPV4", staticRoute.IpVersion)
	expect.EqualsString("StaticRoutes.Routes[0].DestinationNetworkAddress", "132.15.2.0", staticRoute.DestinationNetworkAddress)
	expect.EqualsInt("StaticRoutes.Routes[0].DestinationPrefixSize", 24, staticRoute.DestinationPrefixSize)
	expect.EqualsString("StaticRoutes.Routes[0].NextHopAddress", "132.15.3.2", staticRoute.NextHopAddress)
}

func TestClient_GetStaticRouteByName_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, getStaticRouteByNameTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	domainId := "660a4030-5051-4344-a5a3-f7cbf7c44832"
	staticRoute, err := client.GetStaticRouteByName("ClientStaticRoute", domainId)
	if err != nil {
		test.Fatal(err)
	}

	verifyGetStaticRouteByNameTestResponse(test, staticRoute)
}

var getStaticRouteByNameTestResponse = `
{
    "staticRoute": [
	{
	"id": "9e6b496d-5261-4542-91aa-b50c7f569c54",
	"datacenterId": "AU1",
	"networkDomainId": "9888c372-eb4b-11e3-b29c-001517c4643e", "name": "ClientStaticRoute",
	"description": "This is a Client Static Route",
	"type": "CLIENT",
	"ipVersion": "IPV4",
	"destinationNetworkAddress": "132.15.2.0", "destinationPrefixSize": 24,
	"nextHopAddress": "132.15.3.2",
	"nextHopAddressVlanId": "9888c372-eb4b-11e3-b29c-001517c54433", "state": "NORMAL",
	"createTime": "2018-05-16T12:05:10.000Z"
	} ],
    "pageNumber": 1,
    "pageCount": 1,
    "totalCount": 1,
    "pageSize": 250
}
`

func verifyGetStaticRouteByNameTestResponse(test *testing.T, staticRoute *StaticRoute) {
	expect := expect(test)

	expect.NotNil("StaticRoute", staticRoute)

	expect.EqualsString("StaticRoutes.Routes[0].Name", "ClientStaticRoute", staticRoute.Name)
	expect.EqualsString("StaticRoutes.Routes[0].Description", "This is a Client Static Route", staticRoute.Description)
	expect.EqualsString("StaticRoutes.Routes[0].IpVersion", "IPV4", staticRoute.IpVersion)
	expect.EqualsString("StaticRoutes.Routes[0].DestinationNetworkAddress", "132.15.2.0", staticRoute.DestinationNetworkAddress)
	expect.EqualsInt("StaticRoutes.Routes[0].DestinationPrefixSize", 24, staticRoute.DestinationPrefixSize)
	expect.EqualsString("StaticRoutes.Routes[0].NextHopAddress", "132.15.3.2", staticRoute.NextHopAddress)
}

func TestClient_RestoreStaticRoute_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, restoreStaticRouteTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	err := client.RestoreStaticRoute("9e6b496d-5261-4542-91aa-b50c7f569c54")
	if err != nil {
		test.Fatal(err)
	}

	// Pass
}

var restoreStaticRouteTestResponse = `
{
	"operation": "RESTORE_STATIC_ROUTES",
	"responseCode": "OK",
	"message": "Static Routes have been restored.",
	"requestId": "na/2018-04-14T13:37:20/62f06368-c3fb-11e3-b29c-001517c4643e"
	}
`
