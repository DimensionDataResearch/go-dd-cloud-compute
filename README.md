# Dimension Data Cloud Compute client

The Go client for Dimension Data Cloud Compute targets the CloudControl (MCP) 2.0 API.

API reference documentation can be found [here](https://godoc.org/github.com/DimensionDataResearch/go-dd-cloud-compute/compute).

The client's methods (where possible) match the operations defined in the [CloudControl 2.2 to 2.13 API](https://community.opsourcecloud.net/Public?goto=DocumentRevision.jsp%3FdocId%3D144faed4cf556401f5b94fb1afafa9b5%26from%3DBrowse_e5b1a66815188ad439f76183b401f026).


To get started, run:
```
go get github.com/DimensionDataResearch/go-dd-cloud-compute
```

You can now create a client for the API, and retrieve a listing of your network domains.

```go
package main

import (
	"fmt"
	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
)

region := "AU"
username := "my_user"
password := "my_password"

client := compute.NewClient(region, username, password)
networkDomains, err := client.ListNetworkDomains()
if err != nil {
	panic(err)
}

for _, networkDomain = range(networkDomains) {
	fmt.Printf("Found network domain: '%s'\n", networkDomain.Name)
}

```

To deploy a new server and wait for its deployment to complete:

```go
networkDomainID := "20e05553-226e-4ce4-b953-b837d816a087"
networkDomain, err := client.GetNetworkDomain(networkDomainID)
if err != nil {
	return err
}
if networkDomain == nil {
	return fmt.Errorf("No network domain was found with Id '%s'.", networkDomainID)
}

deploymentConfiguration := compute.ServerDeploymentConfiguration{
	Name:                  "my-server",
	Description:           "This is my server",
	AdministratorPassword: "sn4u$ag3s!",
	Start:                 true,

	// CPU and Memory should be configured after the OS image has been applied to this configuration.

	Network: compute.VirtualMachineNetwork{
		NetworkDomainID: networkDomainID,

		PrimaryAdapter: compute.VirtualMachineNetworkAdapter{
			VLANID:             "5be19198-f270-4b86-b60a-e7787c4d67e4",
			PrivateIPv4Address: "10.0.3.12",
		},
	},

	PrimaryDNS:   "8.8.8.8",
	SecondaryDNS: "8.8.4.4",
}

// Retrieve image details.
//
// The machine will be deployed in the data center where the OS image is located.
osImageName := "CentOS 7 64-bit 2 CPU"
osImage, err := apiClient.FindOSImage(*osImageName, networkDomain.DatacenterID)
if err != nil {
	return err
}
if osImage == nil {
	return fmt.Errorf("Unable to find an OS image named '%s' in data centre '%s' (which is where the target network domain, '%s', is located).", *osImageName, dataCenterID, networkDomainID)
}

// Apply the OS image configuration to your server.
err = deploymentConfiguration.ApplyImage(osImage)
if err != nil {
	return err
}

// Customise memory and / or CPU (if required).
deploymentConfiguration.MemoryGB = 8
deploymentConfiguration.CPU.Count = 4
deploymentConfiguration.CPU.CoresPerSocket = 2

// Initiate the deployment.
serverID, err := apiClient.DeployServer(deploymentConfiguration)
if err != nil {
	return err
}

// Now wait for the deployment to complete.
resource, err := apiClient.WaitForDeploy(compute.ResourceTypeServer, serverID, 25 * time.Minute)
if err != nil {
	return err
}

// When deployment is complete, resource can be cast to a Server to obtain server details (if required).
server := resource.(*compute.Server)
fmt.Printf("Server '%s' (%s) has been successfully deployed. ", server.Name, server.ID)
```