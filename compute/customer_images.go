package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// CustomerImages represents a page of CustomerImage results.
type CustomerImages struct {
	// The current page of network domains.
	Images []CustomerImage `json:"customerImage"`

	// The current page number.
	PageNumber int `json:"pageNumber"`

	// The number of customer images in the current page of results.
	PageCount int `json:"pageCount"`

	// The total number of customer images that match the requested filter criteria (if any).
	TotalCount int `json:"totalCount"`

	// The maximum number of customer images per page.
	PageSize int `json:"pageSize"`
}

// CustomerImage represents a custom virtual machine image.
type CustomerImage struct {
	ID              string                        `json:"id"`
	Name            string                        `json:"name"`
	Description     string                        `json:"description"`
	DataCenterID    string                        `json:"datacenterId"`
	Guest           ImageGuestInformation         `json:"guest"`
	CPU             VirtualMachineCPU             `json:"cpu"`
	MemoryGB        int                           `json:"memoryGb"`
	SCSIControllers VirtualMachineSCSIControllers `json:"scsiControllers"`
	CreateTime      string                        `json:"createTime"`
	State           string                        `json:"state"`
}

// GetID retrieves the image ID.
func (image *CustomerImage) GetID() string {
	return image.ID
}

// GetName retrieves the image name.
func (image *CustomerImage) GetName() string {
	return image.Name
}

// ToEntityReference creates an EntityReference representing the CustomerImage.
func (image *CustomerImage) ToEntityReference() EntityReference {
	return EntityReference{
		ID:   image.ID,
		Name: image.Name,
	}
}

var _ NamedEntity = &CustomerImage{}

// GetResourceType retrieves the resource type.
func (image *CustomerImage) GetResourceType() ResourceType {
	return ResourceTypeCustomerImage
}

// GetState retrieves the resource's current state (e.g. ResourceStatusNormal, etc).
func (image *CustomerImage) GetState() string {
	return image.State
}

// IsDeleted determines whether the resource been deleted (i.e. the underlying struct is nil)?
func (image *CustomerImage) IsDeleted() bool {
	return image == nil
}

var _ Resource = &CustomerImage{}

// GetType determines the image type.
func (image *CustomerImage) GetType() ImageType {
	return ImageTypeCustomer
}

// GetDatacenterID retrieves Id of the datacenter where the image is located.
func (image *CustomerImage) GetDatacenterID() string {
	return image.DataCenterID
}

// GetOS retrieves information about the image's operating system.
func (image *CustomerImage) GetOS() OperatingSystem {
	return image.Guest.OperatingSystem
}

// RequiresCustomization determines whether the image requires guest OS customisation during deployment.
func (image *CustomerImage) RequiresCustomization() bool {
	return image.Guest.OSCustomization
}

// ApplyTo applies the CustomerImage to the specified ServerDeploymentConfiguration.
func (image *CustomerImage) ApplyTo(config *ServerDeploymentConfiguration) {
	config.ImageID = image.ID
	config.CPU = image.CPU
	config.MemoryGB = image.MemoryGB
	config.SCSIControllers = make(VirtualMachineSCSIControllers, len(image.SCSIControllers))
	for index, scsiController := range image.SCSIControllers {
		config.SCSIControllers[index] = scsiController
	}
}

// ApplyToUncustomized applies the CustomerImage to the specified UncustomizedServerDeploymentConfiguration.
func (image *CustomerImage) ApplyToUncustomized(config *UncustomizedServerDeploymentConfiguration) {
	config.ImageID = image.ID
	config.CPU = image.CPU
	config.MemoryGB = image.MemoryGB
	if len(image.SCSIControllers) == 0 {
		return
	}
	config.Disks = make(VirtualMachineDisks, len(image.SCSIControllers[0].Disks))
	for index, disk := range image.SCSIControllers[0].Disks {
		config.Disks[index] = disk
	}
}

var _ Image = &CustomerImage{}

// Request body when exporting a customer image to an OVF package.
type exportCustomerImage struct {
	ImageID          string `json:"imageId"`
	OVFPackagePrefix string `json:"ovfPackagePrefix"`
}

// Request body when importing a customer image from an OVF package.
type importCustomerImage struct {
	OVFPackageManifest   string `json:"ovfPackage"`
	ImageName            string `json:"name"`
	ImageDescription     string `json:"description"`
	DatacenterID         string `json:"datacenterId"`
	GuestOSCustomization bool   `json:"guestOsCustomization"`
}

// GetCustomerImage retrieves a specific customer image by Id.
func (client *Client) GetCustomerImage(id string) (image *CustomerImage, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/image/customerImage/%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(id),
	)
	request, err := client.newRequestV25(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		var apiResponse *APIResponseV2

		apiResponse, err = readAPIResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return nil, err
		}

		if apiResponse.ResponseCode == ResponseCodeResourceNotFound {
			return nil, nil // Not an error, but was not found.
		}

		return nil, apiResponse.ToError("Request to retrieve customer image '%s' failed with status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	image = &CustomerImage{}
	err = json.Unmarshal(responseBody, image)
	if err != nil {
		return nil, err
	}

	return image, nil
}

// FindCustomerImage finds a customer image by name in a given data centre.
func (client *Client) FindCustomerImage(name string, dataCenterID string) (image *CustomerImage, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/image/customerImage?name=%s&datacenterId=%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(name),
		url.QueryEscape(dataCenterID),
	)
	request, err := client.newRequestV25(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		var apiResponse *APIResponseV2

		apiResponse, err = readAPIResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("Request to find customer image '%s' in data centre '%s' failed with status code %d (%s): %s", name, dataCenterID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	images := &CustomerImages{}
	err = json.Unmarshal(responseBody, images)
	if err != nil {
		return nil, err
	}

	if images.PageCount == 0 {
		return nil, nil
	}

	if images.PageCount != 1 {
		return nil, fmt.Errorf("found multiple images (%d) matching '%s' in data centre '%s'", images.TotalCount, name, dataCenterID)
	}

	return &images.Images[0], err
}

// ListCustomerImagesInDatacenter lists all customer images in a given data centre.
func (client *Client) ListCustomerImagesInDatacenter(dataCenterID string, paging *Paging) (images *CustomerImages, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/image/customerImage?datacenterId=%s&%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(dataCenterID),
		paging.EnsurePaging().toQueryParameters(),
	)
	request, err := client.newRequestV25(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		var apiResponse *APIResponseV2

		apiResponse, err = readAPIResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("Request to list customer images in data centre '%s' failed with status code %d (%s): %s", dataCenterID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	images = &CustomerImages{}
	err = json.Unmarshal(responseBody, images)

	return
}

// ImportCustomerImage imports the specified customer image from an OVF package.
//
// The OVF package can be uploaded via FTPS (call GetDatacenter to determine the FTPS end-point for the target datacenter).
//
// The image's status will be ResourceStatusPendingAdd while the import is in progress, then ResourceStatusNormal once the export is complete.
func (client *Client) ImportCustomerImage(imageName string, imageDescription string, preventGuestOSCustomization bool, ovfPackagePrefix string, datacenterID string) (importID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/image/importImage",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV24(requestURI, http.MethodPost, &importCustomerImage{
		ImageName:            imageName,
		ImageDescription:     imageDescription,
		DatacenterID:         datacenterID,
		OVFPackageManifest:   ovfPackagePrefix + ".mf",
		GuestOSCustomization: !preventGuestOSCustomization,
	})
	if err != nil {
		return "", err
	}
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return "", err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return "", err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return "", fmt.Errorf("Request to import customer image '%s' in datacenter '%s' from OFV package '%s' failed with status code %d (%s): %s",
			imageName,
			datacenterID,
			ovfPackagePrefix,
			statusCode,
			apiResponse.ResponseCode,
			apiResponse.Message,
		)
	}

	// Expected: "info" { "name": "imageId", "value": "the-Id-of-new-customer-image" }
	imageIDMessage := apiResponse.GetFieldMessage("imageId")
	if imageIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'imageId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *imageIDMessage, nil
}

// ExportCustomerImage exports the specified customer image to an OVF package.
//
// The OVF package can then be downloaded via FTPS.
//
// The image's status will be ResourceStatusPendingChange while the export is in progress, then ResourceStatusNormal once the export is complete.
func (client *Client) ExportCustomerImage(imageID string, ovfPackagePrefix string) (exportID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/image/exportImage",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV24(requestURI, http.MethodPost, &exportCustomerImage{
		ImageID:          imageID,
		OVFPackagePrefix: ovfPackagePrefix,
	})
	if err != nil {
		return "", err
	}
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return "", err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return "", err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return "", fmt.Errorf("Request to export customer image '%s' with OFV package prefix '%s' failed with status code %d (%s): %s",
			imageID,
			ovfPackagePrefix,
			statusCode,
			apiResponse.ResponseCode,
			apiResponse.Message,
		)
	}

	// Expected: "info" { "name": "imageExportId", "value": "the-Id-of-the-export-operation" }
	imageExportIDMessage := apiResponse.GetFieldMessage("imageExportId")
	if imageExportIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'imageExportId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *imageExportIDMessage, nil
}
