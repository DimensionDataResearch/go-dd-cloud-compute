package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Get user account details (successful).
func TestClient_GetMyAccount_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/xml")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, accountTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.SetBaseAddress(testServer.URL)

	account, err := client.GetMyAccount()
	if err != nil {
		test.Fatal(err)
	}

	verifyAccountTestResponse(test, account)
}

// Get user account details (access denied).
func TestClient_GetMyAccount_AccessDenied(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		http.Error(writer, "Invalid credentials.", http.StatusUnauthorized)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user", "password")
	client.SetBaseAddress(testServer.URL)

	_, err := client.GetMyAccount()
	if err == nil {
		test.Fatal("Client did not return expected access-denied error.")

		return
	}
	if err.Error() != "Cannot connect to compute API (invalid credentials)." {
		test.Fatal("Unexpected error: ", err)
	}
}

/*
 * Test responses.
 */

var accountTestResponse = `
    <?xml version="1.0" encoding="UTF-8" standalone="yes"?>
        <ns3:Account xmlns="http://oec.api.opsource.net/schemas/organization" xmlns:ns2="http://oec.api.opsource.net/schemas/admin" xmlns:ns4="http://oec.api.opsource.net/schemas/server" xmlns:ns3="http://oec.api.opsource.net/schemas/directory" xmlns:ns6="http://oec.api.opsource.net/schemas/datacenter" xmlns:ns5="http://oec.api.opsource.net/schemas/whitelabel" xmlns:ns8="http://oec.api.opsource.net/schemas/backup" xmlns:ns7="http://oec.api.opsource.net/schemas/general" xmlns:ns13="http://oec.api.opsource.net/schemas/serverbootstrap" xmlns:ns9="http://oec.api.opsource.net/schemas/storage" xmlns:ns12="http://oec.api.opsource.net/schemas/vip" xmlns:ns11="http://oec.api.opsource.net/schemas/network" xmlns:ns10="http://oec.api.opsource.net/schemas/manualimport" xmlns:ns16="http://oec.api.opsource.net/schemas/reset" xmlns:ns15="http://oec.api.opsource.net/schemas/multigeo" xmlns:ns14="http://oec.api.opsource.net/schemas/support">
            <ns3:userName>user1</ns3:userName>
            <ns3:fullName>User One</ns3:fullName>
            <ns3:firstName>User</ns3:firstName>
            <ns3:lastName>One</ns3:lastName>
            <ns3:emailAddress>user1@corp.com</ns3:emailAddress>
            <ns3:department>Some Department</ns3:department>
            <ns3:customDefined1></ns3:customDefined1>
            <ns3:customDefined2></ns3:customDefined2>
            <ns3:orgId>cc309bfe-2710-43b7-a6a6-2b7a1965cf63</ns3:orgId>
            <ns3:roles>
                <ns3:role>
                    <ns3:name>server</ns3:name>
                </ns3:role>
                <ns3:role>
                    <ns3:name>network</ns3:name>
                </ns3:role>
                <ns3:role>
                    <ns3:name>create image</ns3:name>
                </ns3:role>
            </ns3:roles>
        </ns3:Account>
`

func verifyAccountTestResponse(test *testing.T, account *Account) {
	if account == nil {
		test.Fatal("Account was nil.")
	}

	if account.UserName != "user1" {
		test.Fatalf("UserName field is '%s' (expected '%s').", account.UserName, "user1")
	}

	if account.FullName != "User One" {
		test.Fatalf("FullName field is '%s' (expected '%s').", account.FullName, "User One")
	}

	if account.FirstName != "User" {
		test.Fatalf("FirstName field is '%s' (expected '%s').", account.FirstName, "User")
	}

	if account.LastName != "One" {
		test.Fatalf("LastName field is '%s' (expected '%s').", account.LastName, "One")
	}

	if account.Department != "Some Department" {
		test.Fatalf("FullName field is '%s' (expected '%s').", account.Department, "Some Department")
	}

	if account.EmailAddress != "user1@corp.com" {
		test.Fatalf("EmailAddress field is '%s' (expected '%s').", account.EmailAddress, "user1@corp.com")
	}

	if account.OrganizationID != "cc309bfe-2710-43b7-a6a6-2b7a1965cf63" {
		test.Fatalf("OrganizationID field is '%s' (expected '%s').", account.OrganizationID, "cc309bfe-2710-43b7-a6a6-2b7a1965cf63")
	}

	if account.AssignedRoles == nil {
		test.Fatal("AssignedRoles.Roles field is nil.")
	}

	if len(account.AssignedRoles) != 3 {
		test.Fatalf("AssignedRoles.Roles field has length %d (expected %d).", len(account.AssignedRoles), 3)
	}

	if account.AssignedRoles[0].Name != "server" {
		test.Fatalf("AssignedRoles[0].Name is '%s' (expected '%s').", account.AssignedRoles[0].Name, "server")
	}

	if account.AssignedRoles[1].Name != "network" {
		test.Fatalf("AssignedRoles[1].Name is '%s' (expected '%s').", account.AssignedRoles[1].Name, "network")
	}

	if account.AssignedRoles[2].Name != "create image" {
		test.Fatalf("AssignedRoles[2].Name is '%s' (expected '%s').", account.AssignedRoles[2].Name, "create image")
	}
}
