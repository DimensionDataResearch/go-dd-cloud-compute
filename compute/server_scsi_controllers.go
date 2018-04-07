package compute

import (
	"fmt"
	"net/http"
	"net/url"
)

// addSCSIControllerToServer represents the request body when adding a SCSI controller to a server.
type addSCSIControllerToServer struct {
	ServerID    string `json:"serverId"`
	AdapterType string `json:"adapterType"`
	BusNumber   int    `json:"busNumber"`
}

// removeSCSIControllerFromServer represents the request body when removing a SCSI controller from a server.
type removeSCSIControllerFromServer struct {
	ControllerID string `json:"id"`
}

// addDiskToSCSIController represents the request body when adding a new disk to a SCSI controller.
type addDiskToSCSIController struct {
	SCSIController scsiController `json:"scsiController"`
	SizeGB         int            `json:"sizeGb"`
	Speed          string         `json:"speed"`
}

// scsiController represents part of the request body when adding a disk to a SCSI controller
type scsiController struct {
	ControllerID string `json:"controllerId"`
	SCSIUnitID   int    `json:"scsiId"`
}

// expandDisk represents the request body when expamding a server disk.
type expandDisk struct {
	// The ID of the disk to expand.
	DiskID string `json:"id"`

	// The new disk size, in gigabytes.
	NewSizeGB int `json:"newSizeGb"`
}

// removeDisk represents the request body when removing an existing disk from a server or SCSI controller.
type removeDisk struct {
	DiskID string `json:"id"`
}

// AddSCSIControllerToServer adds a SCSI controller to a server.
func (client *Client) AddSCSIControllerToServer(serverID string, adapterType string, busNumber int) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/addScsiController",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV26(requestURI, http.MethodPost, &addSCSIControllerToServer{
		ServerID:    serverID,
		AdapterType: adapterType,
		BusNumber:   busNumber,
	})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return apiResponse.ToError("Request to add SCSI controller for bus %d to server '%s' failed with status code %d (%s): %s", busNumber, serverID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// RemoveSCSIControllerFromServer removes an existing SCSI controller from a server.
func (client *Client) RemoveSCSIControllerFromServer(controllerID string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/removeScsiController",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV26(requestURI, http.MethodPost, &removeSCSIControllerFromServer{
		ControllerID: controllerID,
	})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return apiResponse.ToError("Request to remove SCSI controller '%s' failed with status code %d (%s): %s", controllerID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// AddDiskToSCSIController adds a disk to an existing SCSI controller.
func (client *Client) AddDiskToSCSIController(controllerID string, scsiUnitID int, sizeGB int, speed string) (diskID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/server/addDisk",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV26(requestURI, http.MethodPost, &addDiskToSCSIController{
		SCSIController: scsiController{
			ControllerID: controllerID,
			SCSIUnitID:   scsiUnitID,
		},
		SizeGB: sizeGB,
		Speed:  speed,
	})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return "", err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return "", err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return "", apiResponse.ToError("Request to add disk with SCSI Unit ID %d to controller '%s' failed with status code %d (%s): %s", scsiUnitID, controllerID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "diskId", "value": "the-Id-of-the-new-disk" }
	if len(apiResponse.FieldMessages) < 1 || apiResponse.FieldMessages[0].FieldName != "diskId" {
		return "", apiResponse.ToError("Received an unexpected response (missing 'diskId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return apiResponse.FieldMessages[0].Message, nil
}

// ExpandDisk requests expanding of a server / SCSI controller's disk.
func (client *Client) ExpandDisk(diskID string, newSizeGB int) (response *APIResponseV2, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return
	}

	requestURI := fmt.Sprintf("%s/server/expandDisk",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV26(requestURI, http.MethodPost, &expandDisk{
		DiskID:    diskID,
		NewSizeGB: newSizeGB,
	})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return
	}

	response, err = readAPIResponseAsJSON(responseBody, statusCode)

	return
}

// RemoveDisk removes an existing disk from a server or SCSI controller.
func (client *Client) RemoveDisk(diskID string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/removeDisk",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV26(requestURI, http.MethodPost, &removeDisk{
		DiskID: diskID,
	})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return apiResponse.ToError("Request to remove disk '%s' failed with status code %d (%s): %s", diskID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
