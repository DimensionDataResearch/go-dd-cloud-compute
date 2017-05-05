package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetServer_ById_Success(test *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, getServerTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	server, err := client.GetServer("5a32d6e4-9707-4813-a269-56ab4d989f4d")
	if err != nil {
		test.Fatal("Unable to deploy server: ", err)
	}

	verifyGetServerTestResponse(test, server)
}

// Deploy server (successful).
func TestClient_DeployServer_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		deploymentConfiguration := &ServerDeploymentConfiguration{}
		err := readRequestBodyAsJSON(request, deploymentConfiguration)
		if err != nil {
			test.Fatal(err.Error())
		}

		verifyDeployServerTestRequest(test, deploymentConfiguration)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, deployServerTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	serverConfiguration := ServerDeploymentConfiguration{
		Name:                  "Production FTPS Server",
		Description:           "This is the main FTPS Server",
		ImageID:               "02250336-de2b-4e99-ab96-78511b7f8f4b",
		AdministratorPassword: "P$$ssWwrrdGoDd!",
		CPU: VirtualMachineCPU{
			Count:          4,
			CoresPerSocket: 1,
			Speed:          "STANDARD",
		},
		MemoryGB:     4,
		PrimaryDNS:   "10.20.255.12",
		SecondaryDNS: "10.20.255.13",
		Network: VirtualMachineNetwork{
			NetworkDomainID: "484174a2-ae74-4658-9e56-50fc90e086cf",
			PrimaryAdapter: VirtualMachineNetworkAdapter{
				VLANID: stringToPtr("0e56433f-d808-4669-821d-812769517ff8"),
			},
			AdditionalNetworkAdapters: []VirtualMachineNetworkAdapter{
				VirtualMachineNetworkAdapter{
					PrivateIPv4Address: stringToPtr("172.16.0.14"),
				},
				VirtualMachineNetworkAdapter{
					VLANID:      stringToPtr("e0b4d43c-c648-11e4-b33a-72802a5322b2"),
					AdapterType: stringToPtr(NetworkAdapterTypeVMXNET3),
				},
			},
		},
		SCSIControllers: []VirtualMachineSCSIController{
			VirtualMachineSCSIController{
				BusNumber:   0,
				AdapterType: StorageControllerAdapterTypeLSILogicParallel,
				Disks: []VirtualMachineDisk{
					VirtualMachineDisk{
						SCSIUnitID: 0,
						Speed:      "STANDARD",
					},
					VirtualMachineDisk{
						SCSIUnitID: 1,
						Speed:      "HIGHPERFORMANCE",
					},
				},
			},
		},
	}

	serverID, err := client.DeployServer(serverConfiguration)
	if err != nil {
		test.Fatal(err)
	}

	expect.EqualsString("serverID", "7b62aae5-bdbe-4595-b58d-c78f95db2a7f", serverID)
}

// Deploy uncustomised server (successful).
func TestClient_DeployUncustomizedServer_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		deploymentConfiguration := &UncustomizedServerDeploymentConfiguration{}
		err := readRequestBodyAsJSON(request, deploymentConfiguration)
		if err != nil {
			test.Fatal(err.Error())
		}

		verifyDeployUncustomizedServerTestRequest(test, deploymentConfiguration)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, deployUncustomizedServerTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	serverConfiguration := UncustomizedServerDeploymentConfiguration{
		Name:        "Production Server",
		Description: "Uncustomized appliance server.",
		ImageID:     "e926545b-1b9c-4068-8cef-076830a9a0bc",
		CPU: VirtualMachineCPU{
			Count:          9,
			CoresPerSocket: 3,
			Speed:          "ECONOMY",
		},
		MemoryGB: 2,
		Network: VirtualMachineNetwork{
			NetworkDomainID: "e926545b-1b9c-4068-8cef-076830a9a0bc",
			PrimaryAdapter: VirtualMachineNetworkAdapter{
				PrivateIPv4Address: stringToPtr("10.0.1.15"),
				AdapterType:        stringToPtr(NetworkAdapterTypeVMXNET3),
			},
			AdditionalNetworkAdapters: []VirtualMachineNetworkAdapter{
				VirtualMachineNetworkAdapter{
					VLANID:      stringToPtr("e0b4d43c-c648-11e4-b33a-72802a5322b2"),
					AdapterType: stringToPtr(NetworkAdapterTypeVMXNET3),
				},
			},
		},
		Disks: []VirtualMachineDisk{
			VirtualMachineDisk{
				ID:    "d99e4d2a-24c0-4c54-b491-e56697b8f004",
				Speed: "ECONOMY",
			},
			VirtualMachineDisk{
				ID:    "e6a3c0b7-cd32-4224-b8ec-5f1359940204",
				Speed: "HIGHPERFORMANCE",
			},
		},
	}

	serverID, err := client.DeployUncustomizedServer(serverConfiguration)
	if err != nil {
		test.Fatal(err)
	}

	expect.EqualsString("serverID", "7b62aae5-bdbe-4595-b58d-c78f95db2a7f", serverID)
}

// Add disk to server (successful).
func TestClient_AddServerDisk_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestBody := &addDiskToServer{}
		err := readRequestBodyAsJSON(request, requestBody)
		if err != nil {
			test.Fatal(err.Error())
		}

		verifyAddDiskToServerTestRequest(test, requestBody)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, addDiskToServerTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	diskID, err := client.AddDiskToServer("7b62aae5-bdbe-4595-b58d-c78f95db2a7f", 4, 20, "ECONOMY")
	if err != nil {
		test.Fatal(err)
	}

	expect.EqualsString("diskID", "9e6b496d-5261-4542-91aa-b50c7f569c54", diskID)
}

// Resize server disk (successful).
func TestClient_ResizeServerDisk_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		expect.EqualsString(
			"Request.URL",
			"/oec/0.9/dummy-organization-id/server/7b62aae5-bdbe-4595-b58d-c78f95db2a7f/disk/92b1819e-6f91-4abe-88c7-607841959f90/changeSize",
			request.URL.Path,
		)

		requestBody := &resizeServerDisk{}
		err := readRequestBodyAsXML(request, requestBody)
		if err != nil {
			test.Fatal(err.Error())
		}

		verifyResizeServerDiskRequest(test, requestBody)

		writer.Header().Set("Content-Type", "application/xml")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, resizeServerDiskTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	response, err := client.ResizeServerDisk("7b62aae5-bdbe-4595-b58d-c78f95db2a7f", "92b1819e-6f91-4abe-88c7-607841959f90", 23)
	if err != nil {
		test.Fatal(err)
	}

	verifyResizeServerDiskTestResponse(test, response)
}

// Change server disk speed (successful).
func TestClient_ChangeServerDiskSpeed_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		expect.EqualsString(
			"Request.URL",
			"/oec/0.9/dummy-organization-id/server/7b62aae5-bdbe-4595-b58d-c78f95db2a7f/disk/92b1819e-6f91-4abe-88c7-607841959f90/changeSpeed",
			request.URL.Path,
		)

		requestBody := &changeServerDiskSpeed{}
		err := readRequestBodyAsXML(request, requestBody)
		if err != nil {
			test.Fatal(err.Error())
		}

		verifyChangeServerDiskSpeedRequest(test, requestBody)

		writer.Header().Set("Content-Type", "application/xml")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, changeServerDiskSpeedTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	response, err := client.ChangeServerDiskSpeed("7b62aae5-bdbe-4595-b58d-c78f95db2a7f", "92b1819e-6f91-4abe-88c7-607841959f90", ServerDiskSpeedStandard)
	if err != nil {
		test.Fatal(err)
	}

	verifyChangeServerDiskSpeedTestResponse(test, response)
}

// Add Nic (successful).
func TestClient_AddServerNic_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestBody := &addNicConfiguration{}
		err := readRequestBodyAsJSON(request, requestBody)
		if err != nil {
			test.Fatal(err.Error())
		}

		verifyAddNicToServerTestRequest(test, requestBody)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, addNicToServerTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	nicID, err := client.AddNicToServer("1c7762ca-f379-4eef-b08e-aa526d602589", "10.0.3.18", "2e312054-532a-46aa-ab4f-226660bfba6d")
	if err != nil {
		test.Fatal(err)
	}

	expect.EqualsString("nicID", "5999db1d-725c-46ba-9d4e-d33991e61ab1", nicID)
}

// Remove Server Nic (successful).
func TestClient_RemoveServerNic_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestBody, err := readRequestBodyAsString(request)
		if err != nil {
			test.Fatal(err.Error())
		}

		expect.EqualsString("Request.Body",
			`{"id":"5999db1d-725c-46ba-9d4e-d33991e61ab1"}`,
			requestBody,
		)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, removeNicFromServerTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	err := client.RemoveNicFromServer("5999db1d-725c-46ba-9d4e-d33991e61ab1")
	if err != nil {
		test.Fatal(err)
	}

}

// Delete Server (successful).
func TestClient_DeleteServer_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestBody, err := readRequestBodyAsString(request)
		if err != nil {
			test.Fatal("Failed to read request body: ", err)
		}

		expect.EqualsString("Request.Body",
			`{"id":"5b00a2ab-c665-4cd6-8291-0b931374fb3d"}`,
			requestBody,
		)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, deleteServerTestResponse)
	}))
	defer testServer.Close()

	client := NewClientWithBaseAddress(testServer.URL, "user1", "password")
	client.setAccount(&Account{
		OrganizationID: "dummy-organization-id",
	})

	err := client.DeleteServer("5b00a2ab-c665-4cd6-8291-0b931374fb3d")
	if err != nil {
		test.Fatal(err)
	}

	// Pass
}

/*
 * Test requests.
 */

const deployServerTestRequest = `
	{
		"name": "Production FTPS Server",
		"description": "This is the main FTPS Server",
		"imageId": "02250336-de2b-4e99-ab96-78511b7f8f4b",
		"start": true,
		"administratorPassword": "P$$ssWwrrdGoDd!",
		"cpu": {
			"count": 4,
			"coresPerSocket": 1,
			"speed": "STANDARD"
		},
		"memoryGb": 4,
		"primaryDns": "10.20.255.12",
		"secondaryDns": "10.20.255.13",
		"networkInfo": {
			"networkDomainId": "484174a2-ae74-4658-9e56-50fc90e086cf",
			"primaryNic": {
				"vlanId": "0e56433f-d808-4669-821d-812769517ff8"
			},
			"additionalNic": [
				{
					"privateIpv4": "172.16.0.14"
				},
				{
					"vlanId": "e0b4d43c-c648-11e4-b33a-72802a5322b2",
					"networkAdapter": "VMXNET3"
				}
			]
		},
		"disk": [
			{
				"scsiId": "0",
				"speed": "STANDARD"
			},
			{
				"scsiId": "1",
				"speed": "HIGHPERFORMANCE"
			}
		]
	}
`

func verifyDeployServerTestRequest(test *testing.T, deploymentConfiguration *ServerDeploymentConfiguration) {
	expect := expect(test)

	expect.NotNil("ServerDeploymentConfiguration", deploymentConfiguration)
	expect.EqualsString("ServerDeploymentConfiguration.Name", "Production FTPS Server", deploymentConfiguration.Name)
	expect.EqualsString("ServerDeploymentConfiguration.Description", "This is the main FTPS Server", deploymentConfiguration.Description)
	expect.EqualsString("ServerDeploymentConfiguration.ImageID", "02250336-de2b-4e99-ab96-78511b7f8f4b", deploymentConfiguration.ImageID)
	expect.EqualsString("ServerDeploymentConfiguration.AdministratorPassword", "P$$ssWwrrdGoDd!", deploymentConfiguration.AdministratorPassword)

	// CPU
	expect.EqualsInt("ServerDeploymentConfiguration.CPU.Count", 4, deploymentConfiguration.CPU.Count)
	expect.EqualsInt("ServerDeploymentConfiguration.CPU.CoresPerSocket", 1, deploymentConfiguration.CPU.CoresPerSocket)
	expect.EqualsString("ServerDeploymentConfiguration.CPU.Speed", "STANDARD", deploymentConfiguration.CPU.Speed)

	// Memory
	expect.EqualsInt("ServerDeploymentConfiguration.MemoryGB", 4, deploymentConfiguration.MemoryGB)

	// Network.
	network := deploymentConfiguration.Network
	expect.EqualsString("ServerDeploymentConfiguration.Network.NetworkDomainID", "484174a2-ae74-4658-9e56-50fc90e086cf", network.NetworkDomainID)

	expect.EqualsString("ServerDeploymentConfiguration.PrimaryDNS", "10.20.255.12", deploymentConfiguration.PrimaryDNS)
	expect.EqualsString("ServerDeploymentConfiguration.Secondary", "10.20.255.13", deploymentConfiguration.SecondaryDNS)

	expect.NotNil("ServerDeploymentConfiguration.Network.PrimaryAdapter.VLANID", network.PrimaryAdapter.VLANID)
	expect.EqualsString("ServerDeploymentConfiguration.Network.PrimaryAdapter.VLANID", "0e56433f-d808-4669-821d-812769517ff8", *network.PrimaryAdapter.VLANID)

	// Network adapters.
	expect.EqualsInt("ServerDeploymentConfiguration.Network.AdditionalNetworkAdapters.Length", 2, len(network.AdditionalNetworkAdapters))

	expect.NotNil("ServerDeploymentConfiguration.Network.AdditionalNetworkAdapters[0].PrivateIPv4Address", network.AdditionalNetworkAdapters[0].PrivateIPv4Address)
	expect.EqualsString("ServerDeploymentConfiguration.Network.AdditionalNetworkAdapters[0].PrivateIPv4Address", "172.16.0.14", *network.AdditionalNetworkAdapters[0].PrivateIPv4Address)

	expect.NotNil("ServerDeploymentConfiguration.Network.AdditionalNetworkAdapters[1].VLANID", network.AdditionalNetworkAdapters[1].VLANID)
	expect.EqualsString("ServerDeploymentConfiguration.Network.AdditionalNetworkAdapters[1].VLANID", "e0b4d43c-c648-11e4-b33a-72802a5322b2", *network.AdditionalNetworkAdapters[1].VLANID)
	expect.NotNil("ServerDeploymentConfiguration.Network.AdditionalNetworkAdapters[1].AdapterType", network.AdditionalNetworkAdapters[1].AdapterType)
	expect.EqualsString("ServerDeploymentConfiguration.Network.AdditionalNetworkAdapters[1].AdapterType", "VMXNET3", *network.AdditionalNetworkAdapters[1].AdapterType)

	// Disks.
	expect.EqualsInt("ServerDeploymentConfiguration.SCSIControllers.Length", 1, len(deploymentConfiguration.SCSIControllers))
	expect.EqualsInt("ServerDeploymentConfiguration.SCSIControllers[0].Disks.Length", 2, len(deploymentConfiguration.SCSIControllers[0].Disks))

	expect.EqualsInt("ServerDeploymentConfiguration.SCSIControllers[0].Disks[0].SCSIUnitID", 0, deploymentConfiguration.SCSIControllers[0].Disks[0].SCSIUnitID)
	expect.EqualsString("ServerDeploymentConfiguration.SCSIControllers[0].Disks[0].Speed", "STANDARD", deploymentConfiguration.SCSIControllers[0].Disks[0].Speed)

	expect.EqualsInt("ServerDeploymentConfiguration.SCSIControllers[0].Disks[1].SCSIUnitID", 1, deploymentConfiguration.SCSIControllers[0].Disks[1].SCSIUnitID)
	expect.EqualsString("ServerDeploymentConfiguration.SCSIControllers[0].Disks[1].Speed", "HIGHPERFORMANCE", deploymentConfiguration.SCSIControllers[0].Disks[1].Speed)
}

const deployUncustomizedServerTestRequest = `
	{
		"name": "Production Server",
		"description": "Uncustomized appliance server.",
		"imageId": "e926545b-1b9c-4068-8cef-076830a9a0bc",
		"start": false,
		"cpu": {
			"speed": "ECONOMY",
			"count": "9",
			"coresPerSocket": "3"
		},
		"memoryGb": "2",
		"clusterId": "NA9-01",
		"networkInfo": {
			"networkDomainId": "e926545b-1b9c-4068-8cef-076830a9a0bc",
			"primaryNic": {
				"privateIpv4": "10.0.1.15",
				"networkAdapter": "VMXNET3"
			},
			"additionalNic": [
				{
					"vlanId": "e0b4d43c-c648-11e4-b33a-72802a5322b2",
					"networkAdapter": "VMXNET3"
				}
			]
		},
		"disk": [
			{
				"id": "d99e4d2a-24c0-4c54-b491-e56697b8f004",
				"speed": "ECONOMY"
			},
			{
				"id": "e6a3c0b7-cd32-4224-b8ec-5f1359940204",
				"speed": "HIGHPERFORMANCE"
			}
		],
		"tag": [
			{
				"tagKeyName": "department",
				"value": "IT"
			},
			{
				"tagKeyName": "backup",
				"value": "nope"
			}
		]
	}
`

func verifyDeployUncustomizedServerTestRequest(test *testing.T, deploymentConfiguration *UncustomizedServerDeploymentConfiguration) {
	expect := expect(test)

	expect.NotNil("UncustomizedServerDeploymentConfiguration", deploymentConfiguration)
	expect.EqualsString("UncustomizedServerDeploymentConfiguration.Name", "Production Server", deploymentConfiguration.Name)
	expect.EqualsString("UncustomizedServerDeploymentConfiguration.Description", "Uncustomized appliance server.", deploymentConfiguration.Description)
	expect.EqualsString("UncustomizedServerDeploymentConfiguration.ImageID", "e926545b-1b9c-4068-8cef-076830a9a0bc", deploymentConfiguration.ImageID)

	// CPU
	expect.EqualsInt("UncustomizedServerDeploymentConfiguration.CPU.Count", 9, deploymentConfiguration.CPU.Count)
	expect.EqualsInt("UncustomizedServerDeploymentConfiguration.CPU.CoresPerSocket", 3, deploymentConfiguration.CPU.CoresPerSocket)
	expect.EqualsString("UncustomizedServerDeploymentConfiguration.CPU.Speed", "ECONOMY", deploymentConfiguration.CPU.Speed)

	// Memory
	expect.EqualsInt("UncustomizedServerDeploymentConfiguration.MemoryGB", 2, deploymentConfiguration.MemoryGB)

	// Network.
	network := deploymentConfiguration.Network
	expect.EqualsString("UncustomizedServerDeploymentConfiguration.Network.NetworkDomainID", "e926545b-1b9c-4068-8cef-076830a9a0bc", network.NetworkDomainID)

	expect.IsNil("UncustomizedServerDeploymentConfiguration.Network.PrimaryAdapter.VLANID", network.PrimaryAdapter.VLANID)
	expect.NotNil("UncustomizedServerDeploymentConfiguration.Network.PrimaryAdapter.PrivateIPv4Address", network.PrimaryAdapter.PrivateIPv4Address)
	expect.EqualsString("UncustomizedServerDeploymentConfiguration.Network.PrimaryAdapter.PrivateIPv4Address", "10.0.1.15", *network.PrimaryAdapter.PrivateIPv4Address)
	expect.EqualsString("UncustomizedServerDeploymentConfiguration.Network.PrimaryAdapter.AdapterType", "VMXNET3", *network.AdditionalNetworkAdapters[0].AdapterType)

	// Network adapters.
	expect.EqualsInt("UncustomizedServerDeploymentConfiguration.Network.AdditionalNetworkAdapters.Length", 1, len(network.AdditionalNetworkAdapters))

	expect.NotNil("UncustomizedServerDeploymentConfiguration.Network.AdditionalNetworkAdapters[0].VLANID", network.AdditionalNetworkAdapters[0].VLANID)
	expect.EqualsString("UncustomizedServerDeploymentConfiguration.Network.AdditionalNetworkAdapters[0].VLANID", "e0b4d43c-c648-11e4-b33a-72802a5322b2", *network.AdditionalNetworkAdapters[0].VLANID)
	expect.IsNil("UncustomizedServerDeploymentConfiguration.Network.AdditionalNetworkAdapters[0].PrivateIPv4Address", network.AdditionalNetworkAdapters[0].PrivateIPv4Address)
	expect.NotNil("UncustomizedServerDeploymentConfiguration.Network.AdditionalNetworkAdapters[0].AdapterType", network.AdditionalNetworkAdapters[0].AdapterType)
	expect.EqualsString("UncustomizedServerDeploymentConfiguration.Network.AdditionalNetworkAdapters[0].AdapterType", "VMXNET3", *network.AdditionalNetworkAdapters[0].AdapterType)

	// Disks.
	expect.EqualsInt("UncustomizedServerDeploymentConfiguration.Disks.Length", 2, len(deploymentConfiguration.Disks))

	expect.EqualsString("UncustomizedServerDeploymentConfiguration.Disks[0].ID", "d99e4d2a-24c0-4c54-b491-e56697b8f004", deploymentConfiguration.Disks[0].ID)
	expect.EqualsString("UncustomizedServerDeploymentConfiguration.Disks[0].Speed", "ECONOMY", deploymentConfiguration.Disks[0].Speed)

	expect.EqualsString("UncustomizedServerDeploymentConfiguration.Disks[1].ID", "e6a3c0b7-cd32-4224-b8ec-5f1359940204", deploymentConfiguration.Disks[1].ID)
	expect.EqualsString("UncustomizedServerDeploymentConfiguration.Disks[1].Speed", "HIGHPERFORMANCE", deploymentConfiguration.Disks[1].Speed)
}

const addDiskToServerTestRequest = `
	{
		"id": "7b62aae5-bdbe-4595-b58d-c78f95db2a7f",
		"sizeGb": 20,
		"speed": "ECONOMY",
		"scsiId": 4
	}
`

func verifyAddDiskToServerTestRequest(test *testing.T, request *addDiskToServer) {
	expect := expect(test)

	expect.EqualsString("addDiskToServer.ServerID", "7b62aae5-bdbe-4595-b58d-c78f95db2a7f", request.ServerID)
	expect.EqualsString("addDiskToServer.Speed", "ECONOMY", request.Speed)
	expect.EqualsInt("addDiskToServer.SizeGB", 20, request.SizeGB)
	expect.EqualsInt("addDiskToServer.SCSIUnitID", 4, request.SCSIUnitID)
}

const addNicToServerTestRequest = `
	{
		"serverId": "1c7762ca-f379-4eef-b08e-aa526d602589",
		"nic":
		{
			"vlanId": "2e312054-532a-46aa-ab4f-226660bfba6d"
			"privateIpv4": "10.0.3.18",
			"networkAdapter":"E1000"
		}
	}
`

func verifyAddNicToServerTestRequest(test *testing.T, request *addNicConfiguration) {
	expect := expect(test)
	expect.EqualsString("addNicConfiguration.ServerID", "1c7762ca-f379-4eef-b08e-aa526d602589", request.ServerID)
	expect.EqualsString("addNicConfiguration.Nic.PrivateIPv4", "10.0.3.18", request.Nic.PrivateIPv4)
	// VLANID will not be submitted because private IPv4 has been submitted.
	expect.EqualsString("addNicConfiguration.VlanID", "", request.Nic.VlanID)
}

const resizeServerDiskTestRequest = `
	<ChangeDiskSize xmlns="http://oec.api.opsource.net/schemas/server">
		<newSizeGb>23</newSizeGb>
	</ChangeDiskSize>
`

func verifyResizeServerDiskRequest(test *testing.T, request *resizeServerDisk) {
	expect := expect(test)

	expect.NotNil("ReconfigureServer", request)
	expect.EqualsInt("ReconfigureServer.NewSizeGB", 23, request.NewSizeGB)
}

const changeServerDiskSpeedTestRequest = `
	<ChangeServerDiskSpeed xmlns="http://oec.api.opsource.net/schemas/server">
		<speed>STANDARD</speed>
	</ChangeServerDiskSpeed>
`

func verifyChangeServerDiskSpeedRequest(test *testing.T, request *changeServerDiskSpeed) {
	expect := expect(test)

	expect.NotNil("ChangeServerDiskSpeed", request)
	expect.EqualsString("ChangeServerDiskSpeed.Speed", ServerDiskSpeedStandard, request.Speed)
}

const notifyServerIPAddressChangeTestRequest = `
	{
		"nicId": "5999db1d-725c-46ba-9d4e-d33991e61ab1",
		"privateIpv4": "10.0.1.5",
		"ipv6": "fdfe::5a55:caff:fefa::1:9089"
	}
`

func verifyNotifyServerIPAddressChangeTestRequest(test *testing.T, request *notifyServerIPAddressChange) {
	expect := expect(test)

	expect.NotNil("NotifyServerIPAddressChange", request)
	expect.EqualsString("NotifyServerIPAddressChange.AdapterID", "5999db1d-725c-46ba-9d4e-d33991e61ab1", request.AdapterID)

	expect.NotNil("NotifyServerIPAddressChange.IPv4Address", request.IPv4Address)
	expect.EqualsString("NotifyServerIPAddressChange.IPv4Address", "10.0.1.5", *request.IPv4Address)

	expect.NotNil("NotifyServerIPAddressChange.IPv4Address", request.IPv6Address)
	expect.EqualsString("NotifyServerIPAddressChange.IPv6Address", "fdfe::5a55:caff:fefa::1:9089", *request.IPv6Address)
}

const reconfigureServerTestRequest = `
	{
		"memoryGb": 8,
		"cpuCount": 5,
		"cpuSpeed": "STANDARD",
		"coresPerSocket": 1,
		"id": "f8fe7965-3b7c-4cee-827e-f1e0b40a72e0"
	}
`

func verifyReconfigureServerTestRequest(test *testing.T, request *reconfigureServer) {
	expect := expect(test)

	expect.NotNil("ReconfigureServer", request)
	expect.EqualsString("ReconfigureServer.ServerID", "5999db1d-725c-46ba-9d4e-d33991e61ab1", request.ServerID)

	expect.NotNil("ReconfigureServer.MemoryGB", request.MemoryGB)
	expect.EqualsInt("ReconfigureServer.MemoryGB", 8, *request.MemoryGB)

	expect.NotNil("ReconfigureServer.CPUCount", request.CPUCount)
	expect.EqualsInt("ReconfigureServer.CPUCount", 5, *request.CPUCount)
}

/*
 * Test responses.
 */

const getServerTestResponse = `
	{
		"name": "Production Web Server",
		"description": "Server to host our main web application.",
		"operatingSystem": {
			"id": "WIN2008S32",
			"displayName": "WIN2008S/32",
			"family": "WINDOWS"
		},
		"cpu": {
			"count": 2,
			"speed": "STANDARD",
			"coresPerSocket": 1
		},
		"memoryGb": 4,
		"scsiController": [
			{
				"id": "00cbc4-1b3b-49c4-a4e6-697caff4b872",
				"adapterType": "BUS_LOGIC",
				"key": 1000,
				"state": "NORMAL",
				"busNumber": 0,
				"disk": [
					{
						"id": "c2e1f199-116e-4dbc-9960-68720b832b0a",
						"scsiId": 0,
						"sizeGb": 50,
						"speed": "STANDARD",
						"state": "NORMAL"
					}
				]
			}
		],
		"networkInfo": {
			"primaryNic": {
			"id": "5e869800-df7b-4626-bcbf-8643b8be11fd",
			"privateIpv4": "10.0.4.8",
			"ipv6": "2607:f480:1111:1282:2960:fb72:7154:6160",
			"vlanId": "bc529e20-dc6f-42ba-be20-0ffe44d1993f",
			"vlanName": "Production Server",
			"state": "NORMAL"
		},
		"additionalNic": [],
		"networkDomainId": "553f26b6-2a73-42c3-a78b-6116f11291d0" },
		"backup": {
			"assetId": "91002e08-8dc1-47a1-ad33-04f501c06f87",
			"servicePlan": "Advanced",
			"state": "NORMAL"
		},
		"monitoring": {
			"monitoringId": "11049",
			"servicePlan": "ESSENTIALS",
			"state": "NORMAL"
		},
		"softwareLabel": [
			"MSSQL2008R2S"
		],
		"sourceImageId": "3ebf3c0f-90fe-4a8b-8585-6e65b316592c",
		"createTime": "2015-12-02T10:31:33.000Z",
		"deployed": true,
		"started": true,
		"state": "PENDING_CHANGE",
		"progress": {
			"action": "SHUTDOWN_SERVER",
			"requestTime": "2015-12-02T11:07:40.000Z",
			"userName": "devuser1"
		},
		"vmwareTools": {
			"versionStatus": "CURRENT",
			"runningStatus": "RUNNING",
			"apiVersion": 9354
		},
		"virtualHardware": {
			"version": "vmx-08",
			"upToDate": false
		},
		"id": "5a32d6e4-9707-4813-a269-56ab4d989f4d",
		"datacenterId": "NA9"
	}
`

func verifyGetServerTestResponse(test *testing.T, server *Server) {
	expect := expect(test)

	expect.NotNil("Server", server)
	expect.EqualsString("Server.Name", "Production Web Server", server.Name)
	expect.EqualsString("Server.State", ResourceStatusPendingChange, server.State)

	expect.EqualsInt("Server.SCSIControllers.Length", 1, len(server.SCSIControllers))

	controller1 := server.SCSIControllers[0]
	expect.EqualsString("Server.SCSIControllers[0].ID", "00cbc4-1b3b-49c4-a4e6-697caff4b872", controller1.ID)
	expect.EqualsInt("Server.SCSIControllers[0].BusNumber", 0, controller1.BusNumber)
	expect.EqualsString("Server.SCSIControllers[0].AdapterType", "BUS_LOGIC", controller1.AdapterType)
	expect.EqualsString("Server.SCSIControllers[0].State", ResourceStatusNormal, controller1.State)

	controller1Disks := controller1.Disks
	expect.EqualsInt("Server.SCSIControllers[0].Disks.Length", 1, len(controller1Disks))

	disk1 := controller1Disks[0]
	expect.EqualsString("Server.SCSIControllers[0].Disks[0].ID", "c2e1f199-116e-4dbc-9960-68720b832b0a", disk1.ID)
	expect.EqualsInt("Server.SCSIControllers[0].Disks[0].SCSIUnitID", 0, disk1.SCSIUnitID)
	expect.EqualsInt("Server.SCSIControllers[0].Disks[0].SizeGB", 50, disk1.SizeGB)
	expect.EqualsString("Server.SCSIControllers[0].Disks[0].Speed", ServerDiskSpeedStandard, disk1.Speed)
	expect.EqualsString("Server.SCSIControllers[0].Disks[0].State", ResourceStatusNormal, disk1.State)
}

const deployServerTestResponse = `
	{
		"operation": "DEPLOY_SERVER",
		"responseCode": "IN_PROGRESS",
		"message": "Request to deploy Server 'Production FTPS Server' has been accepted and is being processed.",
		"info": [
			{
				"name": "serverId",
				"value": "7b62aae5-bdbe-4595-b58d-c78f95db2a7f"
			}
		],
		"warning": [],
		"error": [],
		"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
	}
`

const deployUncustomizedServerTestResponse = `
	{
        "operation": "DEPLOY_UNCUSTOMIZED_SERVER",
        "responseCode": "IN_PROGRESS",
        "message": "Request to deploy uncustomized Server 'Production Server' has been accepted and is being processed.",
        "info": [
            {
                "name": "serverId",
                "value": "7b62aae5-bdbe-4595-b58d-c78f95db2a7f"
            }
        ],
        "warning": [],
        "error": [],
        "requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
    }
`

func verifyDeployServerTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "DEPLOY_SERVER", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResponseCodeInProgress, response.ResponseCode)
	expect.EqualsString("Response.Message", "Request to deploy Server 'Production FTPS Server' has been accepted and is being processed.", response.Message)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}

const addDiskToServerTestResponse = `
	{
		"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad",
		"operation": "ADD_DISK",
		"responseCode": "IN_PROGRESS",
		"message": "The request to add 20GB Standard Speed Disk to Server 'SERVER-1' has been accepted and is being processed.",
		"info": [
			{
				"name": "diskId",
				"value": "9e6b496d-5261-4542-91aa-b50c7f569c54"
			},
			{
				"name": "scsiId",
				"value": "4"
			},
			{
				"name": "speed",
				"value": "STANDARD"
			}
		]
	}
`

func verifyAddDiskToServerTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "ADD_DISK", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResponseCodeInProgress, response.ResponseCode)
	expect.EqualsString("Response.Message", "The request to add 20GB Standard Speed Disk to Server 'SERVER-1' has been accepted and is being processed.", response.Message)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}

const resizeServerDiskTestResponse = `
	<Status>
		<operation>Change Server Disk Size</operation>
		<result>SUCCESS</result>
		<resultDetail>Server 'Change Server Disk Size' Issued</resultDetail>
		<resultCode>RESULT_0</resultCode>
	</Status>
`

func verifyResizeServerDiskTestResponse(test *testing.T, response *APIResponseV1) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "Change Server Disk Size", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResultSuccess, response.Result)
	expect.EqualsString("Response.Message", "Server 'Change Server Disk Size' Issued", response.Message)
	expect.EqualsString("Response.ResultCode", "RESULT_0", response.ResultCode)
}

const changeServerDiskSpeedTestResponse = `
	<Status>
		<operation>Change Server Disk Speed</operation>
		<result>SUCCESS</result>
		<resultDetail>Change Server Disk Speed Issued</resultDetail>
		<resultCode>RESULT_0</resultCode>
	</Status>
`

func verifyChangeServerDiskSpeedTestResponse(test *testing.T, response *APIResponseV1) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "Change Server Disk Speed", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResultSuccess, response.Result)
	expect.EqualsString("Response.Message", "Change Server Disk Speed Issued", response.Message)
	expect.EqualsString("Response.ResultCode", "RESULT_0", response.ResultCode)
}

const deleteServerTestResponse = `
	{
		"operation": "DELETE_SERVER",
		"responseCode": "IN_PROGRESS",
		"message": "Request to Delete Server (Id:5b00a2ab-c665-4cd6-8291-0b931374fb3d) has been accepted and is being processed",
		"info": [],
		"warning": [],
		"error": [],
		"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
	}
`

func verifyDeleteServerTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "DELETE_SERVER", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResponseCodeInProgress, response.ResponseCode)
	expect.EqualsString("Response.Message", "Request to Delete Server (Id:5b00a2ab-c665-4cd6-8291-0b931374fb3d) has been accepted and is being processed.", response.Message)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}

const addNicToServerTestResponse = `
	{

		"operation": "ADD_NIC",
		"responseCode": "IN_PROGRESS",
		"message": "The request to add NIC for VLAN 'Subsystem VLAN' on Server'Production Mail Server' has been accepted and is being processed",
		"info": [
			{
			"name": "nicId",
			"value": "5999db1d-725c-46ba-9d4e-d33991e61ab1"
			}
		],
		"warning": [],
		"error": [],
		"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
	}
`

func verifyAddNicToServerTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "ADD_NIC", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResponseCodeInProgress, response.ResponseCode)
	expect.EqualsString("Response.Message", "The request to add NIC for VLAN 'Subsystem VLAN' on Server'Production Mail Server' has been accepted and is being processed", response.Message)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}

const removeNicFromServerTestResponse = `
	{

		"operation": "REMOVE_NIC",
		"responseCode": "IN_PROGRESS",
		"message": "Request to Remove NIC 5999db1d-725c-46ba-9d4e-d33991e61ab1 for VLAN 'Subsystem VLAN' from Server 'Production Mail Server' has been accepted and is being processed.",
		"info": [],
		"warning": [],
		"error": [],
		"requestId": "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
	}
`

func verifyRemoveNicFromServerTestResponse(test *testing.T, response *APIResponseV2) {
	expect := expect(test)

	expect.NotNil("APIResponse", response)
	expect.EqualsString("Response.Operation", "REMOVE_NIC", response.Operation)
	expect.EqualsString("Response.ResponseCode", ResponseCodeInProgress, response.ResponseCode)
	expect.EqualsString("Response.Message", "Request to Remove NIC 5999db1d-725c-46ba-9d4e-d33991e61ab1 for VLAN 'Subsystem VLAN' from Server 'Production Mail Server' has been accepted and is being processed.", response.Message)
	expect.EqualsString("Response.RequestID", "na9_20160321T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad", response.RequestID)
}
