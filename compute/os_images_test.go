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

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
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
				"description": "CentOS Release 7.2 64-bit",
				"cpu": {
					"count": 2,
					"speed": "STANDARD",
					"coresPerSocket": 1
				},
				"memoryGb": 4,
				"scsiController": [
					{
						"id": "eec5a912-0fb0-11e7-b626-001b21cfdbe0",
						"adapterType": "LSI_LOGIC_PARALLEL",
						"key": 1000,
						"disk": [
							{
								"id": "fab8c94e-5f8a-4617-8be9-6514fc779bf1",
								"sizeGb": 10,
								"speed": "STANDARD",
								"scsiId": 0
							}
						],
						"busNumber": 0
					}
				],
				"sataController": [],
				"floppy": [],
				"nic": [],
				"softwareLabel": [],
				"createTime": "2016-08-10T14:05:12.000Z",
				"id": "e1b4e0cc-35ba-47be-a2d7-1b5601b87119",
				"datacenterId": "AU9",
				"osImageKey": "T-CENT-7-64-2-4-10",
				"sortOrder": 65,
				"guest": {
					"operatingSystem": {
						"id": "CENTOS764",
						"displayName": "CENTOS7/64",
						"family": "UNIX"
					},
					"osCustomization": true
				}
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
	expect.EqualsString("OSImage.ID", "e1b4e0cc-35ba-47be-a2d7-1b5601b87119", image.ID)
	expect.EqualsString("OSImage.Name", "CentOS 7 64-bit 2 CPU", image.Name)
	expect.EqualsString("OSImage.Description", "CentOS Release 7.2 64-bit", image.Description)

	expect.EqualsString("OSImage.Guest.OperatingSystem.ID", "CENTOS764", image.Guest.OperatingSystem.ID)
	expect.EqualsString("OSImage.Guest.OperatingSystem.DisplayName", "CENTOS7/64", image.Guest.OperatingSystem.DisplayName)
	expect.EqualsString("OSImage.Guest.OperatingSystem.Family", "UNIX", image.Guest.OperatingSystem.Family)

	expect.EqualsInt("OSImage.CPU.Count", 2, image.CPU.Count)
	expect.EqualsString("OSImage.CPU.Speed", "STANDARD", image.CPU.Speed)
	expect.EqualsInt("OSImage.CPU.CoresPerSocket", 1, image.CPU.CoresPerSocket)

	expect.EqualsInt("OSImage.SCSIControllers.Length", 1, len(image.SCSIControllers))

	controller1 := image.SCSIControllers[0]
	expect.EqualsString("OSImage.SCSIControllers[0].ID", "eec5a912-0fb0-11e7-b626-001b21cfdbe0", controller1.ID)
	expect.EqualsString("OSImage.SCSIControllers[0].AdapterType", StorageControllerAdapterTypeLSILogicParallel, controller1.AdapterType)
	expect.EqualsInt("OSImage.SCSIControllers[0].BusNumber", 0, controller1.BusNumber)

	expect.EqualsInt("OSImage.SCSIControllers[0].Disks.Length", 1, len(controller1.Disks))

	disk1 := controller1.Disks[0]
	expect.EqualsString("OSImage.Disks[0].ID", "fab8c94e-5f8a-4617-8be9-6514fc779bf1", disk1.ID)
	expect.EqualsInt("OSImage.Disks[0].SCSIUnitID", 0, disk1.SCSIUnitID)
	expect.EqualsInt("OSImage.Disks[0].SizeGB", 10, disk1.SizeGB)
	expect.EqualsString("OSImage.Disks[0].Speed", "STANDARD", disk1.Speed)

	expect.EqualsString("OSImage.CreateTime", "2016-08-10T14:05:12.000Z", image.CreateTime)
	expect.EqualsString("OSImage.OSImageKey", "T-CENT-7-64-2-4-10", image.OSImageKey)
}
