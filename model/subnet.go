package model

type Subnet struct {
	UUID       string   `json:"uuid,omitempty"`
	Name       string   `json:"name"`
	Network    string   `json:"network"`
	CIDR       int      `json:"cidr"`
	Gateway    string   `json:"gateway,omitempty"`
	Site       *Site    `json:"site,omitempty"`
	DNSServers []string `json:"dnsServers,omitempty"`
	DNSSearch  []string `json:"dnsSearch,omitempty"`
}

type Subnets []*Subnet
