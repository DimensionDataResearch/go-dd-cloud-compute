package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Find OS image by name and data centre (successful).
func TestClient_FindOSImage_By_NameAndDataCenter_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, findOSImageTestResponse)
	}))
	defer testServer.Close()

	client := NewClient("au1", "user1", "password")
	client.setBaseAddress(testServer.URL)
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	image, err := client.FindOSImage("CentOS 7 64-bit 2 CPU", "AU9")
	if err != nil {
		test.Fatal(err)
	}

	verifyFindOSImageTestResponse(test, image)
}

/*
 * Test responses.
 */

var findOSImageTestResponse = `
	{
		"osImage": [
			{
			  "name": "CentOS 7 64-bit 2 CPU",
			  "description": "CentOS Release 7.1 64-bit",
			  "operatingSystem": {
				"id": "CENTOS764",
				"displayName": "CENTOS7/64",
				"family": "UNIX"
			  },
			  "cpu": {
				"count": 2,
				"speed": "STANDARD",
				"coresPerSocket": 1
			  },
			  "memoryGb": 4,
			  "disk": [
				{
				  "id": "55f6780c-bcc6-49d5-8e9b-26c26b6381fa",
				  "scsiId": 0,
				  "sizeGb": 10,
				  "speed": "STANDARD"
				}
			  ],
			  "softwareLabel": [],
			  "createTime": "2015-10-26T10:34:40.000Z",
			  "id": "7e68acb4-bbb8-4206-b30b-0e6c878056bc",
			  "datacenterId": "AU9",
			  "osImageKey": "T-CENT-7-64-2-4-10"
			}
		],
		"pageNumber": 1,
		"pageCount": 1,
		"totalCount": 1,
		"pageSize": 250
	}
`

func verifyFindOSImageTestResponse(test *testing.T, image *OSImage) {
	expect := expect(test)

	expect.NotNil("OSImage", image)
	expect.EqualsString("OSImage.ID", "7e68acb4-bbb8-4206-b30b-0e6c878056bc", image.ID)
	expect.EqualsString("OSImage.Name", "CentOS 7 64-bit 2 CPU", image.Name)
	expect.EqualsString("OSImage.Description", "CentOS Release 7.1 64-bit", image.Description)

	expect.EqualsString("OSImage.OperatingSystem.ID", "CENTOS764", image.OperatingSystem.ID)
	expect.EqualsString("OSImage.OperatingSystem.DisplayName", "CENTOS7/64", image.OperatingSystem.DisplayName)
	expect.EqualsString("OSImage.OperatingSystem.Family", "UNIX", image.OperatingSystem.Family)

	expect.EqualsInt("OSImage.CPU.Count", 2, image.CPU.Count)
	expect.EqualsString("OSImage.CPU.Speed", "STANDARD", image.CPU.Speed)
	expect.EqualsInt("OSImage.CPU.CoresPerSocket", 1, image.CPU.CoresPerSocket)

	expect.NotNil("OSImage.Disks", image.Disks)

	disk1 := image.Disks[0]
	expect.NotNil("OSImage.Disks[0].ID", disk1.ID)
	expect.EqualsString("OSImage.Disks[0].ID", "55f6780c-bcc6-49d5-8e9b-26c26b6381fa", *disk1.ID)
	expect.EqualsInt("OSImage.Disks[0].SCSIUnitID", 0, disk1.SCSIUnitID)
	expect.EqualsInt("OSImage.Disks[0].SizeGB", 10, disk1.SizeGB)
	expect.EqualsString("OSImage.Disks[0].Speed", "STANDARD", disk1.Speed)

	expect.EqualsString("OSImage.CreateTime", "2015-10-26T10:34:40.000Z", image.CreateTime)
	expect.EqualsString("OSImage.OSImageKey", "T-CENT-7-64-2-4-10", image.OSImageKey)
}
