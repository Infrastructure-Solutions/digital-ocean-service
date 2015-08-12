package domain

type DropletRequest struct {
	Name              string                `json:"name"`
	Region            string                `json:"region"`
	Size              string                `json:"size"`
	Image             string                `json:"image"`
	Backups           bool                  `json:"backups"`
	IPv6              bool                  `json:"ipv6"`
	PrivateNetworking bool                  `json:"private_networking"`
	UserData          string                `json:"user_data,omitempty"`
	SSHKeys           []DropletCreateSSHKey `json:"ssh_keys"`
}

type DropletCreateSSHKey struct {
	ID          int    `json:"id"`
	Fingerprint string `json:"fingerprint"`
}

type Droplet struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Region   string     `json:"region"`
	Size     string     `json:"size"`
	UserData string     `json:"user_data,omitempty"`
	Networks []Networks `json:"networks"`
}

type Networks struct {
	V4 []NetworkV4 `json:"v4,omitempty"`
	V6 []NetworkV6 `json:"v6,omitempty"`
}

type NetworkV6 struct {
	IPAddress string `json:"ip_address,omitempty"`
	Netmask   int    `json:"netmask,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Type      string `json:"type,omitempty"`
}
type NetworkV4 struct {
	IPAddress string `json:"ip_address,omitempty"`
	Netmask   string `json:"netmask,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Type      string `json:"type,omitempty"`
}
