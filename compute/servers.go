package compute

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// NetworkAdapterTypeE1000 represents the E1000 network adapter type.
	NetworkAdapterTypeE1000 = "E1000"

	// NetworkAdapterTypeVMXNET3 represents the VMXNET3 network adapter type.
	NetworkAdapterTypeVMXNET3 = "VMXNET3"

	// NetworkAdapterTypeE1000E represents the E1000e network adapter type.
	NetworkAdapterTypeE1000E = "E1000E"

	// NetworkAdapterTypeEnhancedVMXNET2 represents the VMXNET2/Enhanced network adapter type.
	NetworkAdapterTypeEnhancedVMXNET2 = "ENHANCED_VMXNET2"

	// NetworkAdapterTypeFlexiblePCNET32 represents the PCNET32/Flexible network adapter type.
	NetworkAdapterTypeFlexiblePCNET32 = "FLEXIBLE_PCNET32"
)

const (
	// ServerDiskSpeedStandard represents the standard speed for server disks.
	ServerDiskSpeedStandard = "STANDARD"

	// ServerDiskSpeedHighPerformance represents the high-performance speed for server disks.
	ServerDiskSpeedHighPerformance = "HIGHPERFORMANCE"

	// ServerDiskSpeedEconomy represents the economy speed for server disks.
	ServerDiskSpeedEconomy = "ECONOMY"
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
	PrimaryDNS            string                `json:"primaryDns,omitempty"`
	SecondaryDNS          string                `json:"secondaryDns,omitempty"`
	Start                 bool                  `json:"start"`
}

// editServerMetadata represents the request body when modifying server metadata.
type editServerMetadata struct {
	ID          string  `json:"id"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
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

// removeDiskFromServer represents the request body when removing an existing disk from a server.
type removeDiskFromServer struct {
	DiskID string `json:"id"`
}

type serverNic struct {
	VlanID      string  `json:"vlanId,omitempty"`
	PrivateIPv4 string  `json:"privateIpv4,omitempty"`
	AdapterType *string `json:"networkAdapter,omitempty"`
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

// changeServerDiskSpeed represents the request body when changing a server disk's speed.
type changeServerDiskSpeed struct {
	// The XML name for the "resizeServerDisk" data contract
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/server ChangeDiskSpeed"`

	// The new disk speed.
	Speed string `xml:"speed"`
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

// Request body when deleting a network adapter.
type deleteNic struct {
	// The network adapter Id.
	ID string `json:"id"`
}

// Request body when changing network adapter type.
type changeNicType struct {
	// The network adapter Id.
	ID string `json:"nicId"`

	// The network adapter type.
	Type string `json:"networkAdapter"`
}

// GetServer retrieves the server with the specified Id.
// id is the Id of the server to retrieve.
// Returns nil if no server is found with the specified Id.
func (client *Client) GetServer(id string) (server *Server, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/server/server/%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(id),
	)
	request, err := client.newRequestV24(requestURI, http.MethodGet, nil)
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
		url.QueryEscape(organizationID),
		url.QueryEscape(networkDomainID),
		paging.PageNumber,
		paging.PageSize,
	)

	var request *http.Request
	request, err = client.newRequestV23(requestURI, http.MethodGet, nil)
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

	requestURI := fmt.Sprintf("%s/server/deployServer",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV23(requestURI, http.MethodPost, &serverConfiguration)
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

// EditServerMetadata modifies a server's name and / or description.
//
// Pass nil for values you don't want to modify.
func (client *Client) EditServerMetadata(serverID string, name *string, description *string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/editServerMetadata",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV23(requestURI, http.MethodPost, &editServerMetadata{
		ID:          serverID,
		Name:        name,
		Description: description,
	})
	if err != nil {
		return err
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return apiResponse.ToError("Request to modify server '%s' failed with status code %d (%s): %s", serverID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// AddDiskToServer adds a disk to an existing server.
func (client *Client) AddDiskToServer(serverID string, scsiUnitID int, sizeGB int, speed string) (diskID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/server/addDisk",
		url.QueryEscape(organizationID),
	)
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

	requestURI := fmt.Sprintf("%s/server/%s/disk/%s/changeSize",
		url.QueryEscape(organizationID),
		url.QueryEscape(serverID),
		url.QueryEscape(diskID),
	)
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

// ChangeServerDiskSpeed requests changing of a server disk's speed.
func (client *Client) ChangeServerDiskSpeed(serverID string, diskID string, newSpeed string) (response *APIResponseV1, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return
	}

	requestURI := fmt.Sprintf("%s/server/%s/disk/%s/changeSpeed",
		url.QueryEscape(organizationID),
		url.QueryEscape(serverID),
		url.QueryEscape(diskID),
	)
	request, err := client.newRequestV1(requestURI, http.MethodPost, &changeServerDiskSpeed{
		Speed: newSpeed,
	})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return
	}

	response, err = readAPIResponseV1(responseBody, statusCode)

	return
}

// RemoveDiskFromServer removes an existing disk from a server.
func (client *Client) RemoveDiskFromServer(diskID string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/removeDisk",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &removeDiskFromServer{
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

// DeleteServer deletes an existing Server.
// Returns an error if the operation was not successful.
func (client *Client) DeleteServer(id string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/deleteServer",
		url.QueryEscape(organizationID),
	)
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

	requestURI := fmt.Sprintf("%s/server/startServer",
		url.QueryEscape(organizationID),
	)
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

	requestURI := fmt.Sprintf("%s/server/shutdownServer",
		url.QueryEscape(organizationID),
	)
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

	requestURI := fmt.Sprintf("%s/server/powerOffServer",
		url.QueryEscape(organizationID),
	)
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

	requestURI := fmt.Sprintf("%s/server/notifyNicIpChange",
		url.QueryEscape(organizationID),
	)
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

	requestURI := fmt.Sprintf("%s/server/reconfigureServer",
		url.QueryEscape(organizationID),
	)
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

// AddNicToServer adds a network adapter to a server
func (client *Client) AddNicToServer(serverID string, ipv4Address string, vlanID string) (nicID string, err error) {
	return client.addNicToServer(serverID, &serverNic{
		PrivateIPv4: ipv4Address,
		VlanID:      vlanID,
	})
}

// AddNicWithTypeToServer adds a network adapter of a specific type to a server
func (client *Client) AddNicWithTypeToServer(serverID string, ipv4Address string, vlanID string, adapterType string) (nicID string, err error) {
	return client.addNicToServer(serverID, &serverNic{
		PrivateIPv4: ipv4Address,
		VlanID:      vlanID,
		AdapterType: &adapterType,
	})
}

// AddNicToServer adds a network adapter to a server
func (client *Client) addNicToServer(serverID string, nicConfiguration *serverNic) (nicID string, err error) {
	if nicConfiguration == nil {
		return "", fmt.Errorf("Must supply a valid server NIC configuration")
	}

	// Don't submit VLAN ID to CloudControl when IPv4 address has been supplied (one implies the other)
	if nicConfiguration.PrivateIPv4 != "" {
		nicConfiguration.VlanID = ""
	}

	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/server/addNic",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV23(requestURI, http.MethodPost, &addNicConfiguration{
		ServerID: serverID,
		Nic:      *nicConfiguration,
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

	requestURI := fmt.Sprintf("%s/server/removeNic",
		url.QueryEscape(organizationID),
	)
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

// ChangeNetworkAdapterType changes the type of a server's network adapter.
func (client *Client) ChangeNetworkAdapterType(networkAdapterID string, networkAdapterType string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/changeNetworkAdapter",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV24(requestURI, http.MethodPost, &changeNicType{
		ID:   networkAdapterID,
		Type: networkAdapterType,
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
		return apiResponse.ToError("Request to notify change NIC type failed with unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
