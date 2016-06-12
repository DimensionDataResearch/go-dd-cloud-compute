package compute

import (
	"fmt"
	"net/http"
)

// ServerDeploymentConfiguration represents the configuration for deploying a virtual machine.
type ServerDeploymentConfiguration struct {
	Name                  string                `json:"name"`
	Description           string                `json:"description"`
	ImageID               string                `json:"imageId"`
	AdministratorPassword string                `json:"administratorPassword"`
	CPU                   *VirtualMachineCPU    `json:"cpu,omitempty"`
	MemoryGB              *int                  `json:"memoryGb,omitempty"`
	Disks                 []VirtualMachineDisk  `json:"disk"`
	Network               VirtualMachineNetwork `json:"networkInfo"`
	PrimaryDNS            string                `json:"primaryDns"`
	SecondaryDNS          string                `json:"secondaryDns"`
	Start                 bool                  `json:"start"`
}

// ApplyImage applies the specified image (and its default values for CPU, memory, and disks) to the ServerDeploymentConfiguration.
func (config *ServerDeploymentConfiguration) ApplyImage(image *OSImage) error {
	if image == nil {
		return fmt.Errorf("Cannot apply image defaults (no image was supplied).")
	}

	config.ImageID = image.ID
	config.CPU = &image.CPU
	config.MemoryGB = &image.MemoryGB
	config.Disks = make([]VirtualMachineDisk, len(image.Disks))
	for index, disk := range image.Disks {
		config.Disks[index] = disk
	}

	return nil
}

// DeployServer deploys a new virtual machine.
func (client *Client) DeployServer(serverConfiguration ServerDeploymentConfiguration) (serverID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/server/deployServer", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &serverConfiguration)
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return "", err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return "", err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return "", fmt.Errorf("Request to deploy server '%s' failed with status code %d (%s): %s", serverConfiguration.Name, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "serverId", "value": "the-Id-of-the-new-server" }
	if len(apiResponse.FieldMessages) != 1 || apiResponse.FieldMessages[0].FieldName != "serverId" {
		return "", fmt.Errorf("Received an unexpected response (missing 'serverId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return apiResponse.FieldMessages[0].Message, nil
}
