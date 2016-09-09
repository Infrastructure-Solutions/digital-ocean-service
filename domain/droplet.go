package domain

// DropletRequest is the request used to create a droplet in Digital Ocean
type DropletRequest struct {
	Name              string `json:"name"`
	Region            string `json:"region"`
	Size              string `json:"size"`
	Image             string `json:"image"`
	Backups           bool   `json:"backups"`
	IPv6              bool   `json:"ipv6"`
	PrivateNetworking bool   `json:"private_networking"`
	UserData          string `json:"user_data,omitempty"`
	SSHKeys           []Key  `json:"ssh_keys"`
}

// Droplet represents a droplet inside Digital Ocean
type Droplet struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Region            string   `json:"region"`
	InstanceName      string   `json:"size"`
	OperatingSystem   string   `json:"operating_system"`
	IPV6              string   `json:"ipv6,omitempty"`
	PrivateNetworking bool     `json:"private_networking"`
	Networks          Networks `json:"networks"`
	SSHKeys           []Key    `json:"ssh_keys"`
}

// Networks the networks a droplet has
type Networks struct {
	V4 []NetworkV4 `json:"v4"`
	V6 []NetworkV6 `json:"v6"`
}

// NetworkV6 The representation of a V6 network
type NetworkV6 struct {
	IPAddress string `json:"ip_address,omitempty"`
	Netmask   int    `json:"netmask,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Type      string `json:"type,omitempty"`
}

// NetworkV4 represents a V4 network
type NetworkV4 struct {
	IPAddress string `json:"ip_address,omitempty"`
	Netmask   string `json:"netmask,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Type      string `json:"type,omitempty"`
}
