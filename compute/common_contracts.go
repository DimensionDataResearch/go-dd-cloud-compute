package compute

import "fmt"

// EntitySummary is used to group an entity Id and name together for serialisation / deserialisation purposes.
type EntitySummary struct {
	// The entity Id.
	ID string `json:"id"`
	// The entity name.
	Name string `json:"name"`
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
