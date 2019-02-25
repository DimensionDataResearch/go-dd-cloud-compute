package compute

import "fmt"

// Entity represents a Cloud Control entity.
type Entity interface {
	// GetID retrieves the entity's ID.
	GetID() string
}

// NamedEntity represents a named Cloud Control entity.
type NamedEntity interface {
	Entity

	// GetName retrieves the entity's name.
	GetName() string

	// ToEntityReference creates an EntityReference representing the entity.
	ToEntityReference() EntityReference
}

// EntityReference is used to group an entity Id and name together for serialisation / deserialisation purposes.
type EntityReference struct {
	// The entity Id.
	ID string `json:"id"`
	// The entity name.
	Name string `json:"name,omitempty"`
}

// IPRange represents an IPvX range.
type IPRange interface {
	// Convert the IPvX range to a display string.
	ToDisplayString() string
}

// IPv4Range represents an IPv4 network (base address and prefix size)
type IPv4Range struct {
	// The network base address.
	BaseAddress string `json:"address"`
	// The network prefix size.
	PrefixSize int `json:"prefixSize"`
}

// ToDisplayString converts the IPv4 range to a display string.
func (network IPv4Range) ToDisplayString() string {
	return fmt.Sprintf("%s/%d", network.BaseAddress, network.PrefixSize)
}

// IPv6Range represents an IPv6 network (base address and prefix size)
type IPv6Range struct {
	// The network base address.
	BaseAddress string `json:"address"`
	// The network prefix size.
	PrefixSize int `json:"prefixSize"`
}

// ToDisplayString converts the IPv6 range to a display string.
func (network IPv6Range) ToDisplayString() string {
	return fmt.Sprintf("%s/%d", network.BaseAddress, network.PrefixSize)
}

// OperatingSystem represents a well-known operating system for virtual machines.
type OperatingSystem struct {
	// The operating system Id.
	ID string `json:"id"`

	// The operating system type.
	Family string `json:"family"`

	// The operating system display-name.
	DisplayName string `json:"displayName"`
}

// ImageGuestInformation represents guest-related information about a virtual machine image.
type ImageGuestInformation struct {
	OperatingSystem OperatingSystem `json:"operatingSystem"`
	OSCustomization bool            `json:"osCustomization"`
}

// VirtualMachineCPU represents the CPU configuration for a virtual machine.
type VirtualMachineCPU struct {
	Count          int    `json:"count,omitempty"`
	Speed          string `json:"speed,omitempty"`
	CoresPerSocket int    `json:"coresPerSocket,omitempty"`
}

// AttachedVlan represents the VLAN's gatewayAddressing configuration
type AttachedVlan struct {
	GatewayAddressing string `json:"gatewayAddressing,omitempty"`
}

// VirtualMachineSCSIController represents the configuration for a SCSI controller in a virtual machine.
type VirtualMachineSCSIController struct {
	ID          string              `json:"id,omitempty"`
	BusNumber   int                 `json:"busNumber"`
	Key         int                 `json:"key"`
	AdapterType string              `json:"adapterType"`
	Disks       VirtualMachineDisks `json:"disk"`
	State       string              `json:"state,omitempty"`
}

// GetDiskByUnitID retrieves the VirtualMachineDisk (if any) attached to the VirtualMachineSCSIController that matches the specified SCSI logical unit ID.
func (controller *VirtualMachineSCSIController) GetDiskByUnitID(unitID int) *VirtualMachineDisk {
	if controller == nil {
		return nil
	}

	return controller.Disks.GetByUnitID(unitID)
}

// GetDiskCount determines the number of VirtualMachineDisk entries contained in the VirtualMachineSCSIControllers.
func (controllers VirtualMachineSCSIControllers) GetDiskCount() (count int) {
	for _, controller := range controllers {
		count += len(controller.Disks)
	}

	return
}

// VirtualMachineSCSIControllers is an array of VirtualMachineSCSIController that adds various convenience methods.
type VirtualMachineSCSIControllers []VirtualMachineSCSIController

// GetByID retrieves the VirtualMachineSCSIController that matches the specified CloudControl identifier.
func (controllers VirtualMachineSCSIControllers) GetByID(controllerID string) *VirtualMachineSCSIController {
	for _, controller := range controllers {
		if controller.ID == controllerID {
			return &controller
		}
	}

	return nil
}

// GetByBusNumber retrieves the VirtualMachineSCSIController that matches the specified SCSI bus number.
func (controllers VirtualMachineSCSIControllers) GetByBusNumber(busNumber int) *VirtualMachineSCSIController {
	for _, controller := range controllers {
		if controller.BusNumber == busNumber {
			return &controller
		}
	}

	return nil
}

// GetDiskBySCSIPath retrieves the VirtualMachineDisk (if any) attached to a VirtualMachineSCSIController that matches the specified SCSI device path (bus number and unit ID).
func (controllers VirtualMachineSCSIControllers) GetDiskBySCSIPath(busNumber int, unitID int) *VirtualMachineDisk {
	return controllers.GetByBusNumber(busNumber).GetDiskByUnitID(unitID)
}

// VirtualMachineDisk represents the configuration for disk in a virtual machine.
type VirtualMachineDisk struct {
	ID         string `json:"id,omitempty"`
	SCSIUnitID int    `json:"scsiId"`
	SizeGB     int    `json:"sizeGb"`
	Speed      string `json:"speed"`
	State      string `json:"state,omitempty"`
}

// VirtualMachineDisks is an array of VirtualMachineDisk that adds convenience methods.
type VirtualMachineDisks []VirtualMachineDisk

// GetByID retrieves the disk (if any) with the specified Id.
func (disks VirtualMachineDisks) GetByID(id string) *VirtualMachineDisk {
	for index := range disks {
		if disks[index].ID == id {
			return &disks[index]
		}
	}

	return nil
}

// GetByUnitID retrieves the disk (if any) with the specified SCSI unit Id.
func (disks VirtualMachineDisks) GetByUnitID(unitID int) *VirtualMachineDisk {
	for index := range disks {
		if disks[index].SCSIUnitID == unitID {
			return &disks[index]
		}
	}

	return nil
}

// VirtualMachineNetwork represents the networking configuration for a virtual machine.
type VirtualMachineNetwork struct {
	NetworkDomainID           string                         `json:"networkDomainId,omitempty"`
	PrimaryAdapter            VirtualMachineNetworkAdapter   `json:"primaryNic"`
	AdditionalNetworkAdapters []VirtualMachineNetworkAdapter `json:"additionalNic"`
}

// VirtualMachineNetworkAdapter represents the configuration for a virtual machine's network adapter.
// If deploying a new VM, exactly one of VLANID / PrivateIPv4Address must be specified.
//
// AdapterType (if specified) must be either E1000 or VMXNET3.
type VirtualMachineNetworkAdapter struct {
	ID                 *string `json:"id,omitempty"`
	MACAddress         *string `json:"macAddress,omitempty"` // CloudControl v2.4 and higher
	VLANID             *string `json:"vlanId,omitempty"`
	VLANName           *string `json:"vlanName,omitempty"`
	PrivateIPv4Address *string `json:"privateIpv4,omitempty"`
	PrivateIPv6Address *string `json:"ipv6,omitempty"`
	AdapterType        *string `json:"networkAdapter,omitempty"`
	AdapterKey         *int    `json:"key,omitempty"` // CloudControl v2.4 and higher
	State              *string `json:"state,omitempty"`
}

// GetID returns the network adapter's Id.
func (networkAdapter *VirtualMachineNetworkAdapter) GetID() string {
	if networkAdapter.ID == nil {
		return ""
	}

	return *networkAdapter.ID
}

// GetResourceType returns the network domain's resource type.
func (networkAdapter *VirtualMachineNetworkAdapter) GetResourceType() ResourceType {
	return ResourceTypeNetworkAdapter
}

// GetName returns the network adapter's name (actually Id, since adapters don't have names).
func (networkAdapter *VirtualMachineNetworkAdapter) GetName() string {
	return networkAdapter.GetID()
}

// GetState returns the network adapter's current state.
func (networkAdapter *VirtualMachineNetworkAdapter) GetState() string {
	if networkAdapter.State == nil {
		return ""
	}

	return *networkAdapter.State
}

// IsDeleted determines whether the network adapter has been deleted (is nil).
func (networkAdapter *VirtualMachineNetworkAdapter) IsDeleted() bool {
	return networkAdapter == nil
}

// ToEntityReference creates an EntityReference representing the CustomerImage.
func (networkAdapter *VirtualMachineNetworkAdapter) ToEntityReference() EntityReference {
	id := ""
	if networkAdapter.ID != nil {
		id = *networkAdapter.ID
	}
	name := ""
	if networkAdapter.VLANName != nil {
		name = *networkAdapter.VLANName
	}

	return EntityReference{
		ID:   id,
		Name: name,
	}
}

var _ Resource = &VirtualMachineNetworkAdapter{}
