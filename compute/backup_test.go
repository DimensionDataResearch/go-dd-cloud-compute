package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
)

// Get server backup details (successful).
func TestClient_GetServerBackupDetails_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/xml")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, getServerBackupDetailsResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	backupDetails, err := client.GetServerBackupDetails("5a32d6e4-9707-4813-a269-56ab4d989f4d")
	if err != nil {
		test.Fatal("Failed to retrieve server backup details: ", err)
	}

	verifyGetServerBackupDetailsResponse(test, backupDetails)
}

// Enable backup for server (successful).
func TestClient_EnableServerBackup_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/xml")
		writer.WriteHeader(http.StatusOK)

		enableBackup := &newBackup{}
		err := readRequestBodyAsXML(request, enableBackup)
		if err != nil {
			test.Fatal(
				errors.Wrap(err, "failed to deserialise request body"),
			)
		}

		fmt.Fprintln(writer, enableBackupForServerResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	err := client.EnableServerBackup("5a32d6e4-9707-4813-a269-56ab4d989f4d", BackupServicePlanAdvanced)
	if err != nil {
		test.Fatal("Failed to enable backup for server: ", err)
	}
}

// Get server configured backup client types (successful).
func TestClient_GetServerBackupClientTypes_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/xml")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, getServerBackupClientTypesResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	backupClientTypes, err := client.GetServerBackupClientTypes("5a32d6e4-9707-4813-a269-56ab4d989f4d")
	if err != nil {
		test.Fatal("Failed to retrieve server backup client types: ", err)
	}

	verifyGetServerBackupClientTypesResponse(test, backupClientTypes)
}

// Get server configured backup storage policies (successful).
func TestClient_GetServerBackupStoragePolicies_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/xml")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, getServerBackupStoragePoliciesResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	backupStoragePolicies, err := client.GetServerBackupStoragePolicies("5a32d6e4-9707-4813-a269-56ab4d989f4d")
	if err != nil {
		test.Fatal("Failed to retrieve server backup storage policies: ", err)
	}

	verifyGetServerBackupStoragePoliciesResponse(test, backupStoragePolicies)
}

/*
 * Test requests.
 */

const enableServerBackupRequest = `
<NewBackup xmlns="http://oec.api.opsource.net/schemas/backup" servicePlan="Advanced"/>
`

func verifyEnableServerBackupTestRequest(test *testing.T, request *newBackup) {
	expect := expect(test)

	expect.EqualsString("NewBackup.ServicePlan", BackupServicePlanAdvanced, request.ServicePlan)
}

/*
 * Test responses.
 */

const getServerBackupDetailsResponse = `
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<BackupDetails assetId="30fea016-db07-4935-919d-bfd18caaae62" servicePlan="Enterprise" state="NORMAL" xmlns="http://oec.api.opsource.net/schemas/backup">
	<backupClient id="f1921082-94bc-491d-9136-0e7072df4823" type="FA.Linux" isFileSystem="true" status="Active">
		<description>Linux File Agent</description>
		<schedulePolicyName>12AM - 6AM</schedulePolicyName>
		<storagePolicyName>14 Day Storage Policy + Secondary Copy</storagePolicyName>
		<alerting trigger="ON_SUCCESS">
			<emailAddress>backup@example.com</emailAddress>
		</alerting>
		<times lastBackup="2013-11-07T14:42:56" nextBackup="2013-1108T06:00:00"/>
		<totalBackupSizeGb>1</totalBackupSizeGb>
		<downloadUrl>http://10.161.212.55:8081/PCS/BackupClientInstallerDownload/7786e cd024ed47900a068f27e7867a08049d6d8c</downloadUrl>
	</backupClient>
</BackupDetails>
`

func verifyGetServerBackupDetailsResponse(test *testing.T, response *ServerBackupDetails) {
	expect := expect(test)

	expect.EqualsString("ServerBackupDetails.AssetID", "30fea016-db07-4935-919d-bfd18caaae62", response.AssetID)
	expect.EqualsString("ServerBackupDetails.ServicePlan", BackupServicePlanEnterprise, response.ServicePlan)
	expect.EqualsString("ServerBackupDetails.State", ResourceStatusNormal, response.State)

	expect.EqualsInt("ServerBackupDetails.Clients.Length", 1, len(response.Clients))

	expect.EqualsString("ServerBackupDetails.Clients[0].ID", "f1921082-94bc-491d-9136-0e7072df4823", response.Clients[0].ID)
	expect.EqualsString("ServerBackupDetails.Clients[0].Type", "FA.Linux", response.Clients[0].Type)
	expect.EqualsString("ServerBackupDetails.Clients[0].Description", "Linux File Agent", response.Clients[0].Description)
	expect.EqualsString("ServerBackupDetails.Clients[0].SchedulePolicyName", "12AM - 6AM", response.Clients[0].SchedulePolicyName)
	expect.EqualsString("ServerBackupDetails.Clients[0].StoragePolicyName", "14 Day Storage Policy + Secondary Copy", response.Clients[0].StoragePolicyName)
	expect.EqualsInt("ServerBackupDetails.Clients[0].TotalBackupSizeGb", 1, response.Clients[0].TotalBackupSizeGb)
}

const enableBackupForServerResponse = `
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Status>
	<operation>Enable Backup for Server</operation>
	<result>SUCCESS</result>
	<resultDetail>Backup enabled for Server.</resultDetail>
	<resultCode>REASON_0</resultCode>
</Status>
`

const getServerBackupClientTypesResponse = `
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<BackupClientTypes xmlns="http://oec.api.opsource.net/schemas/backup">
	<backupClientType type="FA.Linux" isFileSystem="true" description="Linux File system"/>
</BackupClientTypes> 
`

func verifyGetServerBackupClientTypesResponse(test *testing.T, response *BackupClientTypes) {
	expect := expect(test)

	expect.EqualsInt("BackupClientTypes.Items.Length", 1, len(response.Items))

	expect.EqualsString("BackupClientTypes.Items[0].Type", "FA.Linux", response.Items[0].Type)
	expect.IsTrue("BackupClientTypes.Items[0].IsFileSystem", response.Items[0].IsFileSystem)
	expect.EqualsString("BackupClientTypes.Items[0].Description", "Linux File system", response.Items[0].Description)
}

const getServerBackupStoragePoliciesResponse = `
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<BackupStoragePolicies xmlns="http://oec.api.opsource.net/schemas/backup">
	<storagePolicy retentionPeriodInDays="10" name="10 Day Storage Policy" />
	<storagePolicy retentionPeriodInDays="30" name="30 Day Storage Policy + Secondary Copy" secondaryLocation="Primary"/>
</BackupStoragePolicies>
`

func verifyGetServerBackupStoragePoliciesResponse(test *testing.T, response *BackupStoragePolicies) {
	expect := expect(test)

	expect.EqualsInt("BackupStoragePolicies.Items.Length", 2, len(response.Items))

	storagePolicy := response.Items[0]
	expect.EqualsString("BackupStoragePolicies.Items[0].Name", "10 Day Storage Policy", storagePolicy.Name)
	expect.EqualsInt("BackupStoragePolicies.Items[0].RetentionPeriodInDays", 10, storagePolicy.RetentionPeriodInDays)
	expect.EqualsString("BackupStoragePolicies.Items[0].SecondaryLocation", "", storagePolicy.SecondaryLocation)

	storagePolicy = response.Items[1]
	expect.EqualsString("BackupStoragePolicies.Items[1].Name", "30 Day Storage Policy + Secondary Copy", storagePolicy.Name)
	expect.EqualsInt("BackupStoragePolicies.Items[1].RetentionPeriodInDays", 30, storagePolicy.RetentionPeriodInDays)
	expect.EqualsString("BackupStoragePolicies.Items[1].SecondaryLocation", "Primary", storagePolicy.SecondaryLocation)
}
