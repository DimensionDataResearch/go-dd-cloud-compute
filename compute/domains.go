package compute

// DeployNetworkDomain represents a request to deploy a compute network domain.
type DeployNetworkDomain struct {
	// The network domain name.
	Name string `json:"name"`

	// The network domain description.
	Description string `json:"description"`

	// The network domain type.
	Type string `json:"type"`

	// The Id of the data centre in which the network domain is located.
	DatacenterID string `json:"datacenterId"`
}

// NetworkDomain represents a compute network domain.
type NetworkDomain struct {
	// The network domain Id.
	ID string `json:"id"`

	// The network domain name.
	Name string `json:"name"`

	// The network domain description.
	Description string `json:"description"`

	// The network domain type.
	Type string `json:"type"`

	// Network domain's NAT IPv4 address.
	NatIPv4Address string `json:"snatIpv4Address"`

	// The network domain creation timestamp.
	CreateTime string `json:"createTime"`

	// The network domain's current state.
	State string `json:"state"`

	// The network domain's current progress (if any).
	Progress string `json:"progress"`

	// The Id of the data centre in which the network domain is located.
	DatacenterID string `json:"datacenter"`
}

// NetworkDomains represents the response to a "List Network Domains" API call.
// It also contains fields common to all API responses (see ApiResponse for a list of all common fields).
type NetworkDomains struct {
	// The current page of network domains.
	Domains []NetworkDomain `json:"networkDomain"`

	// The current page number.
	PageNumber int `json:"pageNumber"`

	// The number of network domains in the current page of results.
	PageCount int `json:"pageCount"`

	// The total number of network domains that match the requested filter criteria (if any).
	TotalCount int `json:"totalCount"`

	// The maximum number of network domains per page.
	PageSize int `json:"pageSize"`
}
