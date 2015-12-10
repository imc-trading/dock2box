package client

import "encoding/json"

// HostResource structure.
type HostResource struct {
	Client *Client
}

// Host structure.
type Host struct {
	Host       string          `json:"host"`
	Build      bool            `json:"build"`
	Debug      bool            `json:"debug"`
	GPT        bool            `json:"gpt"`
	ImageID    string          `json:"imageId"`
	Version    string          `json:"version"`
	KOpts      string          `json:"kOpts"`
	TenantID   string          `json:"tenantId"`
	Labels     []string        `json:"labels"`
	SiteID     string          `json:"siteId"`
	Interfaces []HostInterface `json:"interfaces,omitempty"`
}

// HostInterface structure.
type HostInterface struct {
	Interface string `json:"interface"`
	DHCP      bool   `json:"dhcp"`
	IPv4      string `json:"ipv4,omitempty"`
	HwAddr    string `json:"hwAddr"`
	SubnetID  string `json:"subnetId,omitempty"`
}

// JSON output for a host.
func (h *Host) JSON() []byte {
	b, _ := json.MarshalIndent(h, "", "  ")
	return b
}

// JSON output for a host interface.
func (i *HostInterface) JSON() []byte {
	b, _ := json.MarshalIndent(i, "", "  ")
	return b
}

// Create host.
func (r *HostResource) Create(h *Host) (*Host, error) {
	c := *r.Client
	j, err := c.Create("/hosts", h)
	if err != nil {
		return nil, err
	}

	host := &Host{}
	if err := json.Unmarshal(j, host); err != nil {
		return nil, err
	}

	return host, nil
}

// Exist host.
func (r *HostResource) Exist(name string) (bool, error) {
	c := *r.Client
	s, err := c.Exist("/hosts", name)
	if err != nil {
		return false, err
	}

	return s, nil
}

// Get host.
func (r *HostResource) Get(name string) (*Host, error) {
	c := *r.Client
	j, err := c.Get("/hosts", name)
	if err != nil {
		return nil, err
	}

	host := &Host{}
	if err := json.Unmarshal(j, host); err != nil {
		return nil, err
	}

	return host, nil
}

// GetByID host.
func (r *HostResource) GetByID(id string) (*Host, error) {
	c := *r.Client
	j, err := c.Get("/hosts/id", id)
	if err != nil {
		return nil, err
	}

	host := &Host{}
	if err := json.Unmarshal(j, host); err != nil {
		return nil, err
	}

	return host, nil
}

// All hosts.
func (r *HostResource) All() (*[]Host, error) {
	c := *r.Client
	j, err := c.All("/hosts")
	if err != nil {
		return nil, err
	}

	hosts := &[]Host{}
	if err := json.Unmarshal(j, hosts); err != nil {
		return nil, err
	}

	return hosts, nil
}
