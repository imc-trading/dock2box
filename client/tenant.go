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

// Get tenant.
func (r *TenantResource) Get(id string) (*Tenant, error) {
	c := *r.Client
	j, err := c.Get("/tenants", id)
	if err != nil {
		return nil, err
	}

	tenant := &Tenant{}
	if err := json.Unmarshal(j, tenant); err != nil {
		return nil, err
	}

	return tenant, nil
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

// Update tenant.
func (r *TenantResource) Update(id string, s *Tenant) (*Tenant, error) {
	c := *r.Client
	j, err := c.Update("/tenants/", id, s)
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
func (r *TenantResource) Delete(id string) (*Tenant, error) {
	c := *r.Client
	j, err := c.Delete("/tenants", id)
	if err != nil {
		return nil, err
	}

	tenant := &Tenant{}
	if err := json.Unmarshal(j, tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}
