package compute

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

// Server represents a virtual machine.
type Server struct {
	ID              string                `json:"id"`
	Name            string                `json:"name"`
	Description     string                `json:"description"`
	OperatingSystem OperatingSystem       `json:"operatingSystem"`
	CPU             VirtualMachineCPU     `json:"cpu"`
	MemoryGB        int                   `json:"memoryGb"`
	Disks           []VirtualMachineDisk  `json:"disk"`
	Network         VirtualMachineNetwork `json:"networkInfo"`
	SourceImageID   string                `json:"sourceImageId"`
	State           string                `json:"state"`
	Deployed        bool                  `json:"deployed"`
	Started         bool                  `json:"started"`
}

// GetID returns the server's Id.
func (server *Server) GetID() string {
	return server.ID
}

// GetResourceType returns the network domain's resource type.
func (server *Server) GetResourceType() ResourceType {
	return ResourceTypeServer
}

// GetName returns the server's name.
func (server *Server) GetName() string {
	return server.Name
}

// GetState returns the server's current state.
func (server *Server) GetState() string {
	return server.State
}

// IsDeleted determines whether the server has been deleted (is nil).
func (server *Server) IsDeleted() bool {
	return server == nil
}

var _ Resource = &Server{}

// ToEntityReference creates an EntityReference representing the Server.
func (server *Server) ToEntityReference() EntityReference {
	return EntityReference{
		ID:   server.ID,
		Name: server.Name,
	}
}

var _ NamedEntity = &Server{}

// Servers represents a page of Server results.
type Servers struct {
	Items []Server `json:"server"`

	PagedResult
}

// ServerSummary respresents summary information for a server.
type ServerSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// TODO: Consider adding the "networkingDetails" field.
}

// ToEntityReference creates an EntityReference representing the ServerSummary.
func (serverSummary *ServerSummary) ToEntityReference() EntityReference {
	return EntityReference{
		ID:   serverSummary.ID,
		Name: serverSummary.Name,
	}
}

// ServerDeploymentConfiguration represents the configuration for deploying a virtual machine.
type ServerDeploymentConfiguration struct {
	Name                  string                `json:"name"`
	Description           string                `json:"description"`
	ImageID               string                `json:"imageId"`
	AdministratorPassword string                `json:"administratorPassword"`
	CPU                   VirtualMachineCPU     `json:"cpu"`
	MemoryGB              int                   `json:"memoryGb,omitempty"`
	Disks                 []VirtualMachineDisk  `json:"disk"`
	Network               VirtualMachineNetwork `json:"networkInfo"`
	PrimaryDNS            string                `json:"primaryDns"`
	SecondaryDNS          string                `json:"secondaryDns"`
	Start                 bool                  `json:"start"`
}

// NotifyServerIPAddressChange represents the request body when notifying the system that the IP address for a server's network adapter has changed.
// Exactly at least 1 of IPv4Address or IPv6Address must be specified.
type notifyServerIPAddressChange struct {
	// The server's network adapter Id.
	AdapterID string `json:"nicId"`

	// The server's new private IPv4 address.
	IPv4Address *string `json:"privateIpv4,omitempty"`

	// The server's new private IPv6 address.
	IPv6Address *string `json:"ipv6,omitempty"`
}

// reconfigureServer represents the request body when updating a server's configuration (e.g. memory, CPU count).
type reconfigureServer struct {
	ServerID          string  `json:"id"`
	MemoryGB          *int    `json:"memoryGb,omitempty"`
	CPUCount          *int    `json:"cpuCount,omitempty"`
	CPUCoresPerSocket *int    `json:"coresPerSocket,omitempty"`
	CPUSpeed          *string `json:"cpuSpeed,omitempty"`
}

// addDiskToServer represents the request body when adding a new disk to a server.
type addDiskToServer struct {
	ServerID   string `json:"id"`
	SizeGB     int    `json:"sizeGb"`
	Speed      string `json:"speed"`
	SCSIUnitID int    `json:"scsiId"`
}

type serverNic struct {
	VlanID      string `json:"vlanId,omitempty"`
	PrivateIPv4 string `json:"privateIpv4,omitempty"`
}

// addNicConfiguration represents the request body when adding the new nic.
type addNicConfiguration struct {
	ServerID string    `json:"serverId"`
	Nic      serverNic `json:"nic"`
}

// resizeServerDisk represents the request body when resizing a server disk.
type resizeServerDisk struct {
	// The XML name for the "resizeServerDisk" data contract
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/server ChangeDiskSize"`

	// The new disk size, in gigabytes.
	NewSizeGB int `xml:"newSizeGb"`
}

// ApplyOSImage applies the specified OS image (and its default values for CPU, memory, and disks) to the ServerDeploymentConfiguration.
func (config *ServerDeploymentConfiguration) ApplyOSImage(image *OSImage) error {
	if image == nil {
		return fmt.Errorf("Cannot apply image defaults (no image was supplied).")
	}

	config.ImageID = image.ID
	config.CPU = image.CPU
	config.MemoryGB = image.MemoryGB
	config.Disks = make([]VirtualMachineDisk, len(image.Disks))
	for index, disk := range image.Disks {
		config.Disks[index] = disk
	}

	return nil
}

// ApplyCustomerImage applies the specified OS image (and its default values for CPU, memory, and disks) to the ServerDeploymentConfiguration.
func (config *ServerDeploymentConfiguration) ApplyCustomerImage(image *CustomerImage) error {
	if image == nil {
		return fmt.Errorf("Cannot apply image defaults (no image was supplied).")
	}

	config.ImageID = image.ID
	config.CPU = image.CPU
	config.MemoryGB = image.MemoryGB
	config.Disks = make([]VirtualMachineDisk, len(image.Disks))
	for index, disk := range image.Disks {
		config.Disks[index] = disk
	}

	return nil
}

// Request body when deleting a server.
type deleteServer struct {
	// The server Id.
	ID string `json:"id"`
}

// Request body when starting a server.
type startServer struct {
	// The server Id.
	ID string `json:"id"`
}

// Request body when stopping or powering off a server.
type stopServer struct {
	// The server Id.
	ID string `json:"id"`
}

// Request body when deleting a server.
type deleteNic struct {
	// The server Id.
	ID string `json:"id"`
}

// GetServer retrieves the server with the specified Id.
// id is the Id of the server to retrieve.
// Returns nil if no server is found with the specified Id.
func (client *Client) GetServer(id string) (server *Server, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/server/server/%s", organizationID, id)
	request, err := client.newRequestV22(requestURI, http.MethodGet, nil)
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

		return nil, apiResponse.ToError("Request to retrieve Server failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	server = &Server{}
	err = json.Unmarshal(responseBody, server)

	return server, err
}

// ListServersInNetworkDomain retrieves a page of servers in the specified network domain.
func (client *Client) ListServersInNetworkDomain(networkDomainID string, paging *Paging) (servers Servers, err error) {
	if paging == nil {
		paging = &Paging{
			PageNumber: 1,
		}
	}
	paging.ensureValidPageSize()

	var organizationID string
	organizationID, err = client.getOrganizationID()
	if err != nil {
		return
	}

	requestURI := fmt.Sprintf("%s/server/server?networkDomainId=%s&pageNumber=%d&pageSize=%d",
		organizationID,
		networkDomainID,
		paging.PageNumber,
		paging.PageSize,
	)

	var request *http.Request
	request, err = client.newRequestV22(requestURI, http.MethodGet, nil)
	if err != nil {
		return
	}

	var (
		responseBody []byte
		statusCode   int
	)
	responseBody, statusCode, err = client.executeRequest(request)
	if err != nil {
		return
	}

	if statusCode != http.StatusOK {
		var apiResponse *APIResponseV2
		apiResponse, err = readAPIResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return
		}

		err = apiResponse.ToError("Request to retrieve Server failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)

		return
	}

	servers = Servers{}
	err = json.Unmarshal(responseBody, &servers)

	return
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
		return "", apiResponse.ToError("Request to deploy server '%s' failed with status code %d (%s): %s", serverConfiguration.Name, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "serverId", "value": "the-Id-of-the-new-server" }
	serverIDMessage := apiResponse.GetFieldMessage("serverId")
	if serverIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'serverId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *serverIDMessage, nil
}

// AddDiskToServer adds a disk to an existing server.
func (client *Client) AddDiskToServer(serverID string, scsiUnitID int, sizeGB int, speed string) (diskID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/server/addDisk", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &addDiskToServer{
		ServerID:   serverID,
		SizeGB:     sizeGB,
		SCSIUnitID: scsiUnitID,
		Speed:      speed,
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
		return "", apiResponse.ToError("Request to add disk with SCSI Unit ID %d to server '%s' failed with status code %d (%s): %s", scsiUnitID, serverID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "diskId", "value": "the-Id-of-the-new-disk" }
	if len(apiResponse.FieldMessages) < 1 || apiResponse.FieldMessages[0].FieldName != "diskId" {
		return "", apiResponse.ToError("Received an unexpected response (missing 'diskId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return apiResponse.FieldMessages[0].Message, nil
}

// ResizeServerDisk requests resizing of a server disk.
func (client *Client) ResizeServerDisk(serverID string, diskID string, newSizeGB int) (response *APIResponseV1, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return
	}

	requestURI := fmt.Sprintf("%s/server/%s/disk/%s/changeSize", organizationID, serverID, diskID)
	request, err := client.newRequestV1(requestURI, http.MethodPost, &resizeServerDisk{
		NewSizeGB: newSizeGB,
	})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return
	}

	response, err = readAPIResponseV1(responseBody, statusCode)

	return
}

// DeleteServer deletes an existing Server.
// Returns an error if the operation was not successful.
func (client *Client) DeleteServer(id string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/deleteServer", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &deleteServer{id})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return apiResponse.ToError("Request to delete server failed with unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// StartServer requests that the specified server be started.
func (client *Client) StartServer(id string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/startServer", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &startServer{id})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return apiResponse.ToError("Request to delete server failed with unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// ShutdownServer requests that the specified server be shut down (gracefully, if possible).
func (client *Client) ShutdownServer(id string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/shutdownServer", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &stopServer{id})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return apiResponse.ToError("Request to shut down server failed with unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// PowerOffServer requests that the specified server be powered off (hard shut-down).
func (client *Client) PowerOffServer(id string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/powerOffServer", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &stopServer{id})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return apiResponse.ToError("Request to power off server failed with unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// NotifyServerIPAddressChange notifies the system that the IP address for a server's network adapter has changed.
// serverNetworkAdapterID is the Id of the server's network adapter.
// Must specify at least one of newIPv4Address / newIPv6Address.
func (client *Client) NotifyServerIPAddressChange(networkAdapterID string, newIPv4Address *string, newIPv6Address *string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/notifyNicIpChange", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &notifyServerIPAddressChange{
		AdapterID:   networkAdapterID,
		IPv4Address: newIPv4Address,
		IPv6Address: newIPv6Address,
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
		return apiResponse.ToError("Request to notify change of server IP address failed with unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// ReconfigureServer updates the configuration for a server.
// serverID is the Id of the server.
func (client *Client) ReconfigureServer(serverID string, memoryGB *int, cpuCount *int, cpuCoresPerSocket *int, cpuSpeed *string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/reconfigureServer", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &reconfigureServer{
		ServerID:          serverID,
		MemoryGB:          memoryGB,
		CPUCount:          cpuCount,
		CPUCoresPerSocket: cpuCoresPerSocket,
		CPUSpeed:          cpuSpeed,
	})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeOK && apiResponse.ResponseCode != ResponseCodeInProgress {
		return apiResponse.ToError("Request to reconfigure server failed with unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

//AddNicToServer adds the nic to the server
func (client *Client) AddNicToServer(serverID string, ipv4Address string, vlanID string) (nicID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	var serverNicConfiguration = serverNic{
		PrivateIPv4: ipv4Address,
		VlanID:      vlanID,
	}
	requestURI := fmt.Sprintf("%s/server/addNic", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &addNicConfiguration{
		ServerID: serverID,
		Nic:      serverNicConfiguration,
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
		return "", apiResponse.ToError("Request to notify add a nic failed with unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}
	nicIDMessage := apiResponse.GetFieldMessage("nicId")
	if nicIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'nicId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}
	return *nicIDMessage, nil
}

//RemoveNicFromServer removes the Nic from the server
func (client *Client) RemoveNicFromServer(networkAdapterID string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/removeNic", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &deleteNic{ID: networkAdapterID})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return apiResponse.ToError("Request to notify remove a nic failed with unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
