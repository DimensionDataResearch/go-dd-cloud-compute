package compute

import "encoding/xml"

// NetworkDomain represents a compute network domain.
type NetworkDomain struct {
	XMLName         xml.Name `xml:"networkDomain"`
	Name            string   `xml:"name"`
	Description     string   `xml:"description"`
	Type            string   `xml:"type"`
	SnatIpv4Address string   `xml:"snatIpv4Address"`
	CreateTime      string   `xml:"createTime"`
	State           string   `xml:"state"`
	Progress        string   `xml:"progress"`
	Id              string   `xml:"id"`
	DatacenterId    string   `xml:"datacenterId"`
}

// NetworkDomains represents a page of network domains.
type NetworkDomains struct {
	XMLName             xml.Name        `xml:"networkDomains"`
	NetworkDomain       []NetworkDomain `xml:"networkDomain"`
	PageNumber          int             `xml:"pageNumber"`
	PageNumberSpecified bool            `xml:"pageNumberSpecified"`
	PageCount           int             `xml:"pageCount"`
	PageCountSpecified  bool            `xml:"pageCountSpecified"`
	TotalCount          int             `xml:"totalCount"`
	TotalCountSpecified bool            `xml:"totalCountSpecified"`
	PageSize            int             `xml:"pageSize"`
	PageSizeSpecified   bool            `xml:"pageSizeSpecified"`
}
