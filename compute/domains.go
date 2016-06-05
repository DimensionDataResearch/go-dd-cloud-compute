package compute

// NetworkDomain represents a compute network domain.
type NetworkDomain struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Type            string   `json:"type"`
	SnatIpv4Address string   `json:"snatIpv4Address"`
	CreateTime      string   `json:"createTime"`
	State           string   `json:"state"`
	Progress        string   `json:"progress"`
	ID              string   `json:"id"`
	DatacenterID    string   `json:"datacenterId"`
}

// NetworkDomains represents a page of network domains.
type NetworkDomains struct {
	Domains       		[]NetworkDomain	`json:"networkDomain"`
	PageNumber          int             `json:"pageNumber"`
	PageNumberSpecified bool            `json:"pageNumberSpecified"`
	PageCount           int             `json:"pageCount"`
	PageCountSpecified  bool            `json:"pageCountSpecified"`
	TotalCount          int             `json:"totalCount"`
	TotalCountSpecified bool            `json:"totalCountSpecified"`
	PageSize            int             `json:"pageSize"`
	PageSizeSpecified   bool            `json:"pageSizeSpecified"`
}
