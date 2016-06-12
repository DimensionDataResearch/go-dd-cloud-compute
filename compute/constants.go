package compute

const (
	// ResourceStatusNormal indicates that a resource is active.
	ResourceStatusNormal = "NORMAL"

	// ResourceStatusPendingAdd indicates that an add operation is pending for the resource.
	ResourceStatusPendingAdd = "PENDING_ADD"

	// ResourceStatusPendingDelete indicates that a delete operation is pending for the resource.
	ResourceStatusPendingDelete = "PENDING_DELETE"
)
