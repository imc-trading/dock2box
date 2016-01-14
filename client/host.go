package client

import "encoding/json"

// HostResource structure.
type HostResource struct {
	Client *Client
}

// Host structure.
type Host struct {
	ID       string   `json:"id,omitempty"`
	Host     string   `json:"host,omitempty"`
	Build    bool     `json:"build"`
	Debug    bool     `json:"debug"`
	GPT      bool     `json:"gpt"`
	TagID    string   `json:"tagId,omitempty"`
	KOpts    string   `json:"kOpts,omitempty"`
	TenantID string   `json:"tenantId,omitempty"`
	Labels   []string `json:"labels,omitempty"`
	SiteID   string   `json:"siteId,omitempty"`
}

// JSON output for a host.
func (h *Host) JSON() []byte {
	b, _ := json.MarshalIndent(h, "", "  ")
	return b
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

// Query for hosts.
func (r *HostResource) Query(cond map[string]string) (*[]Host, error) {
	c := *r.Client
	j, err := c.Query("/hosts", cond)
	if err != nil {
		return nil, err
	}

	hosts := &[]Host{}
	if err := json.Unmarshal(j, hosts); err != nil {
		return nil, err
	}

	return hosts, nil
}

// Get host.
func (r *HostResource) Get(id string) (*Host, error) {
	c := *r.Client
	j, err := c.Get("/hosts", id)
	if err != nil {
		return nil, err
	}

	host := &Host{}
	if err := json.Unmarshal(j, host); err != nil {
		return nil, err
	}

	return host, nil
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

// Update host.
func (r *HostResource) Update(id string, h *Host) (*Host, error) {
	c := *r.Client
	j, err := c.Update("/hosts/", id, h)
	if err != nil {
		return nil, err
	}

	host := &Host{}
	if err := json.Unmarshal(j, host); err != nil {
		return nil, err
	}

	return host, nil
}

// Delete host.
func (r *HostResource) Delete(id string) (*Host, error) {
	c := *r.Client
	j, err := c.Delete("/hosts", id)
	if err != nil {
		return nil, err
	}

	host := &Host{}
	if err := json.Unmarshal(j, host); err != nil {
		return nil, err
	}

	return host, nil
}
