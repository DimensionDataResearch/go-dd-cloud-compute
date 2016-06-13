package compute

import "fmt"

// Resources are an abstraction over the various types of entities in the DD compute API

const (
	// ResourceTypeNetworkDomain represents a network domain.
	ResourceTypeNetworkDomain = "network_domain"

	// ResourceTypeVLAN represents a VLAN.
	ResourceTypeVLAN = "vlan"

	// ResourceTypeServer represents a virtual machine.
	ResourceTypeServer = "server"
)

// GetResourceDescription retrieves a textual description of the specified resource type.
func GetResourceDescription(resourceType string) (string, error) {
	switch resourceType {
	case ResourceTypeNetworkDomain:
		return "Network domain", nil

	case ResourceTypeVLAN:
		return "VLAN", nil

	case ResourceTypeServer:
		return "Server", nil

	default:
		return "", fmt.Errorf("Unrecognised resource type '%s'.", resourceType)
	}
}

// GetResource retrieves a compute resource of the specified type by Id.
// id is the resource Id.
// resourceType is the resource type (e.g. ResourceTypeNetworkDomain, ResourceTypeVLAN, etc).
func (client *Client) GetResource(id string, resourceType string) (Resource, error) {
	var resourceLoader func(client *Client, id string) (resource Resource, err error)

	switch resourceType {
	case ResourceTypeNetworkDomain:
		resourceLoader = getNetworkDomainByID

	case ResourceTypeVLAN:
		resourceLoader = getVLANByID

	case ResourceTypeServer:
		resourceLoader = getServerByID

	default:
		return nil, fmt.Errorf("Unrecognised resource type '%s'.", resourceType)
	}

	return resourceLoader(client, id)
}

func getNetworkDomainByID(client *Client, id string) (networkDomain Resource, err error) {
	return client.GetNetworkDomain(id)
}

func getVLANByID(client *Client, id string) (Resource, error) {
	return client.GetVLAN(id)
}

func getServerByID(client *Client, id string) (Resource, error) {
	return client.GetServer(id)
}
