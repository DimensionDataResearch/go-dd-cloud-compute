package compute

// ImageType represents a type of Image.
type ImageType int

const (
	// ImageTypeUnknown represents an unknown image type.
	ImageTypeUnknown ImageType = iota

	// ImageTypeOS represents an OS (built-in) image.
	ImageTypeOS

	// ImageTypeCustomer represents a customer image.
	ImageTypeCustomer
)

// ImageTypeName gets the name of the specified image type.
func ImageTypeName(imageType ImageType) string {
	switch imageType {
	case ImageTypeOS:
		return "OS"
	case ImageTypeCustomer:
		return "Customer"
	default:
		return "Unknown"
	}
}

// Image represents an image used to create servers.
type Image interface {
	Resource

	// GetType determines the image type.
	GetType() ImageType

	// GetDatacenterID retrieves Id of the datacenter where the image is located.
	GetDatacenterID() string

	// GetOS retrieves information about the image's operating system.
	GetOS() OperatingSystem

	// ApplyTo applies the Image to the specified ServerDeploymentConfiguration.
	ApplyTo(config *ServerDeploymentConfiguration)
}
