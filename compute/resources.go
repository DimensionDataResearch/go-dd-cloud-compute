package compute

import (
	"fmt"
	"strings"
)

// Resources are an abstraction over the various types of entities in the DD compute API

// ResourceType represents a well-known DD compute resource type.
type ResourceType int

const (
	// ResourceTypeNetworkDomain represents a network domain.
	ResourceTypeNetworkDomain ResourceType = iota

	// ResourceTypeVLAN represents a VLAN.
	ResourceTypeVLAN

	// ResourceTypeServer represents a virtual machine.
	ResourceTypeServer

	// ResourceTypeServerAntiAffinityRule represents a server anti-affinity rule.
	ResourceTypeServerAntiAffinityRule

	// ResourceTypeNetworkAdapter represents a network adapter in a virtual machine.
	// Note that when calling methods such as WaitForChange, the Id must be of the form 'serverId/networkAdapterId'.
	ResourceTypeNetworkAdapter

	// ResourceTypePublicIPBlock represents a block of public IP addresses.
	ResourceTypePublicIPBlock

	// ResourceTypeFirewallRule represents a firewall rule.
	ResourceTypeFirewallRule

	// ResourceTypeVIPNode represents a VIP node.
	ResourceTypeVIPNode

	// ResourceTypeVIPPool represents a VIP pool.
	ResourceTypeVIPPool

	// ResourceTypeVirtualListener represents a virtual listener.
	ResourceTypeVirtualListener

	// ResourceTypeOSImage represents an OS image.
	ResourceTypeOSImage

	// ResourceTypeCustomerImage represents a customer image.
	ResourceTypeCustomerImage

	// ResourceTypeSSLDomainCertificate represents an SSL certificate for a domain name.
	ResourceTypeSSLDomainCertificate

	// ResourceTypeSSLCertificateChain represents an SSL certificate chain
	ResourceTypeSSLCertificateChain

	// ResourceTypeSSLOffloadProfile represents an SSL-offload profile
	ResourceTypeSSLOffloadProfile

	// ResourceStaticRoutes represents network domain static routes
	ResourceTypeStaticRoutes
)

// Resource represents a compute resource.
type Resource interface {
	NamedEntity

	// The resource type.
	GetResourceType() ResourceType

	// The resource's current state (e.g. ResourceStatusNormal, etc).
	GetState() string

	// Has the resource been deleted (i.e. the underlying struct is nil)?
	IsDeleted() bool
}

// GetResourceDescription retrieves a textual description of the specified resource type.
func GetResourceDescription(resourceType ResourceType) (string, error) {
	switch resourceType {
	case ResourceTypeNetworkDomain:
		return "network domain", nil

	case ResourceTypeVLAN:
		return "VLAN", nil

	case ResourceTypeServer:
		return "server", nil

	case ResourceTypeServerAntiAffinityRule:
		return "server anti-affinity rule", nil

	case ResourceTypeNetworkAdapter:
		return "network adapter", nil

	case ResourceTypePublicIPBlock:
		return "public IPv4 address block", nil

	case ResourceTypeFirewallRule:
		return "Firewall rule", nil

	case ResourceTypeVIPNode:
		return "VIP node", nil

	case ResourceTypeVIPPool:
		return "VIP pool", nil

	case ResourceTypeVirtualListener:
		return "virtual listener", nil

	case ResourceTypeOSImage:
		return "OS image", nil

	case ResourceTypeCustomerImage:
		return "customer image", nil

	case ResourceTypeSSLDomainCertificate:
		return "SSL domain certificate", nil

	case ResourceTypeSSLCertificateChain:
		return "SSL certificate chain", nil

	case ResourceTypeSSLOffloadProfile:
		return "SSL-offload profile", nil

	case ResourceTypeStaticRoutes:
		return "Static Routes", nil

	default:
		return "", fmt.Errorf("unrecognised resource type (value = %d)", resourceType)
	}
}

// GetResource retrieves a compute resource of the specified type by Id.
// id is the resource Id.
// resourceType is the resource type (e.g. ResourceTypeNetworkDomain, ResourceTypeVLAN, etc).
func (client *Client) GetResource(id string, resourceType ResourceType) (Resource, error) {
	switch resourceType {
	case ResourceTypeNetworkDomain:
		return client.GetNetworkDomain(id)

	case ResourceTypeVLAN:
		return client.GetVLAN(id)

	case ResourceTypeServer:
		return client.GetServer(id)

	case ResourceTypeServerAntiAffinityRule:
		return client.getServerAntiAffinityRuleByQualifiedID(id)

	case ResourceTypeNetworkAdapter:
		return client.getNetworkAdapterByID(id)

	case ResourceTypePublicIPBlock:
		return client.GetPublicIPBlock(id)

	case ResourceTypeFirewallRule:
		return client.GetFirewallRule(id)

	case ResourceTypeVIPNode:
		return client.GetVIPNode(id)

	case ResourceTypeVIPPool:
		return client.GetVIPPool(id)

	case ResourceTypeVirtualListener:
		return client.GetVirtualListener(id)

	case ResourceTypeOSImage:
		return client.GetCustomerImage(id)

	case ResourceTypeCustomerImage:
		return client.GetCustomerImage(id)

	case ResourceTypeSSLDomainCertificate:
		return client.GetSSLDomainCertificate(id)

	case ResourceTypeSSLCertificateChain:
		return client.GetSSLCertificateChain(id)

	case ResourceTypeSSLOffloadProfile:
		return client.GetSSLOffloadProfile(id)

	case ResourceTypeStaticRoutes:
		return client.GetStaticRoute(id)
	}

	return nil, fmt.Errorf("unrecognised resource type (value = %d)", resourceType)
}

func (client *Client) getNetworkAdapterByID(id string) (Resource, error) {
	compositeIDComponents := strings.Split(id, "/")
	if len(compositeIDComponents) != 2 {
		return nil, fmt.Errorf("'%s' is not a valid network adapter Id (when loading as a resource, the Id must be of the form 'serverId/networkAdapterId')", id)
	}

	server, err := client.GetServer(compositeIDComponents[0])
	if err != nil {
		return nil, err
	}
	if server == nil {
		return nil, fmt.Errorf("No server found with Id '%s.'", compositeIDComponents[0])
	}

	var targetAdapterID = compositeIDComponents[1]
	if *server.Network.PrimaryAdapter.ID == targetAdapterID {
		return &server.Network.PrimaryAdapter, nil
	}

	for _, adapter := range server.Network.AdditionalNetworkAdapters {
		if *adapter.ID == targetAdapterID {
			return &adapter, nil
		}
	}

	return nil, nil
}

// Retrieve a server anti-affinity rule by qualified ID ("networkDomainId/ruleId").
func (client *Client) getServerAntiAffinityRuleByQualifiedID(id string) (Resource, error) {
	compositeIDComponents := strings.Split(id, "/")
	if len(compositeIDComponents) != 2 {
		return nil, fmt.Errorf("'%s' is not a valid network adapter Id (when loading as a resource, the Id must be of the form 'serverId/networkAdapterId')", id)
	}

	networkDomainID := compositeIDComponents[0]
	ruleID := compositeIDComponents[1]

	rule, err := client.GetServerAntiAffinityRule(ruleID, networkDomainID)
	if err != nil {
		return nil, err
	}

	return rule, nil
}

func getPublicIPBlockByID(client *Client, id string) (Resource, error) {
	return client.GetPublicIPBlock(id)
}

func getFirewallRuleByID(client *Client, id string) (Resource, error) {
	return client.GetFirewallRule(id)
}
