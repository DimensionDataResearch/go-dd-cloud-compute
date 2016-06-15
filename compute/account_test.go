package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Get user account details (successful).
func TestClient_GetAccount_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/xml")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, accountTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)

	account, err := client.GetAccount()
	if err != nil {
		test.Fatal(err)
	}

	verifyAccountTestResponse(test, account)
}

// Get user account details (access denied).
func TestClient_GetAccount_AccessDenied(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		http.Error(writer, "Invalid credentials.", http.StatusUnauthorized)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user", "password")
	client.setBaseAddress(testServer.URL)

	_, err := client.GetAccount()
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
         <ns3:orgId>cc309bfe-1234-43b7-a6a6-2b7a1965cf63</ns3:orgId>
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
	expect := expect(test)

	expect.NotNil("Account", account)

	expect.EqualsString("Account.UserName", "user1", account.UserName)
	expect.EqualsString("Account.FullName", "User One", account.FullName)
	expect.EqualsString("Account.FirstName", "User", account.FirstName)
	expect.EqualsString("Account.LastName", "One", account.LastName)
	expect.EqualsString("Account.Department", "Some Department", account.Department)
	expect.EqualsString("Account.EmailAddress", "user1@corp.com", account.EmailAddress)
	expect.EqualsString("Account.OrganizationID", "cc309bfe-1234-43b7-a6a6-2b7a1965cf63", account.OrganizationID)

	expect.NotNil("Account.AssignedRoles", account.AssignedRoles)
	expect.EqualsInt("Account.AssignedRoles size", 3, len(account.AssignedRoles))

	role1 := account.AssignedRoles[0]
	expect.EqualsString("AssignedRoles[0].Name", "server", role1.Name)

	role2 := account.AssignedRoles[1]
	expect.EqualsString("AssignedRoles[1].Name", "network", role2.Name)

	role3 := account.AssignedRoles[2]
	expect.EqualsString("AssignedRoles[2].Name", "create image", role3.Name)
}
