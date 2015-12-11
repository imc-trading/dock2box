package client

import "encoding/json"

// SubnetResource structure.
type SubnetResource struct {
	Client *Client
}

// Subnet structure.
type Subnet struct {
	ID     string `json:"id"`
	Subnet string `json:"subnet"`
	Mask   string `json:"mask"`
	Gw     string `json:"gw"`
	SiteID string `json:"siteId"`
}

// Create subnet.
func (r *SubnetResource) Create(s *Subnet) (*Subnet, error) {
	c := *r.Client
	j, err := c.Create("/subnets", s)
	if err != nil {
		return nil, err
	}

	subnet := &Subnet{}
	if err := json.Unmarshal(j, subnet); err != nil {
		return nil, err
	}

	return subnet, nil
}

// Delete subnet.
func (r *SubnetResource) Delete(name string) (*Subnet, error) {
	c := *r.Client
	j, err := c.Delete("/subnets", name)
	if err != nil {
		return nil, err
	}

	subnet := &Subnet{}
	if err := json.Unmarshal(j, subnet); err != nil {
		return nil, err
	}

	return subnet, nil
}

// Get subnet.
func (r *SubnetResource) Get(name string) (*Subnet, error) {
	c := *r.Client
	j, err := c.Get("/subnets", name)
	if err != nil {
		return nil, err
	}

	subnet := &Subnet{}
	if err := json.Unmarshal(j, subnet); err != nil {
		return nil, err
	}

	return subnet, nil
}

// All subnets.
func (r *SubnetResource) All() (*[]Subnet, error) {
	c := *r.Client
	j, err := c.All("/subnets")
	if err != nil {
		return nil, err
	}

	subnets := &[]Subnet{}
	if err := json.Unmarshal(j, subnets); err != nil {
		return nil, err
	}

	return subnets, nil
}
