package client

import "encoding/json"

// TenantResource structure.
type TenantResource struct {
	Client *Client
}

// Tenant structure.
type Tenant struct {
	ID     string `json:"id"`
	Tenant string `json:"tenant"`
}

// JSON output for a tenant.
func (t *Tenant) JSON() []byte {
	b, _ := json.MarshalIndent(t, "", "  ")
	return b
}

// Create tenant.
func (r *TenantResource) Create(s *Tenant) (*Tenant, error) {
	c := *r.Client
	j, err := c.Create("/tenants", s)
	if err != nil {
		return nil, err
	}

	tenant := &Tenant{}
	if err := json.Unmarshal(j, tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

// Delete tenant.
func (r *TenantResource) Delete(name string) (*Tenant, error) {
	c := *r.Client
	j, err := c.Delete("/tenants", name)
	if err != nil {
		return nil, err
	}

	tenant := &Tenant{}
	if err := json.Unmarshal(j, tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

// Get tenant.
func (r *TenantResource) Get(name string) (*Tenant, error) {
	c := *r.Client
	j, err := c.Get("/tenants", name)
	if err != nil {
		return nil, err
	}

	tenant := &Tenant{}
	if err := json.Unmarshal(j, tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

// All tenants.
func (r *TenantResource) All() (*[]Tenant, error) {
	c := *r.Client
	j, err := c.All("/tenants")
	if err != nil {
		return nil, err
	}

	tenants := &[]Tenant{}
	if err := json.Unmarshal(j, tenants); err != nil {
		return nil, err
	}

	return tenants, nil
}
