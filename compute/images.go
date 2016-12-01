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
	NamedEntity

	// GetType determines the image type.
	GetType() ImageType

	// GetID retrieves the image ID.
	GetID() string

	// GetName retrieves the image name.
	GetName() string

	// ApplyTo applies the Image to the specified ServerDeploymentConfiguration.
	ApplyTo(config *ServerDeploymentConfiguration)
}
