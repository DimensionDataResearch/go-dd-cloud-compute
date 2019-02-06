package compute

const (
	// ResourceStatusNormal indicates that a resource is active.
	ResourceStatusNormal = "NORMAL"

	// ResourceStatusPendingAdd indicates that an add operation is pending for the resource.
	ResourceStatusPendingAdd = "PENDING_ADD"

	// ResourceStatusPendingChange indicates that a change operation is pending for the resource.
	ResourceStatusPendingChange = "PENDING_CHANGE"

	// ResourceStatusPendingDelete indicates that a delete operation is pending for the resource.
	ResourceStatusPendingDelete = "PENDING_DELETE"

	// ResourceStatusDeleted is a pseudo-status indicates that a resource has been deleted.
	ResourceStatusDeleted = ""
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
	// StorageControllerAdapterTypeBusLogicParallel represents the BusLogic Parallel storage controller adapter type.
	StorageControllerAdapterTypeBusLogicParallel = "BUSLOGIC_PARALLEL"

	// StorageControllerAdapterTypeLSILogicParallel represents the LSI Logic Parallel storage controller adapter type.
	StorageControllerAdapterTypeLSILogicParallel = "LSI_LOGIC_PARALLEL"

	// StorageControllerAdapterTypeLSILogicSAS represents the LSI Logic SAS storage controller adapter type.
	StorageControllerAdapterTypeLSILogicSAS = "LSI_LOGIC_SAS"

	// StorageControllerAdapterTypeEnhancedVMWareParavirtual represents the VMWare Paravirtual storage controller adapter type.
	StorageControllerAdapterTypeEnhancedVMWareParavirtual = "VMWARE_PARAVIRTUAL"
)

const (
	// ServerDiskSpeedEconomy represents the economy speed for server disks.
	ServerDiskSpeedEconomy = "ECONOMY"

	// ServerDiskSpeedStandard represents the standard speed for server disks.
	ServerDiskSpeedStandard = "STANDARD"

	// ServerDiskSpeedHighPerformance represents the high-performance speed for server disks.
	ServerDiskSpeedHighPerformance = "HIGHPERFORMANCE"
)
