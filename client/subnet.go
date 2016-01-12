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

// JSON output for a subnet.
func (s *Subnet) JSON() []byte {
	b, _ := json.MarshalIndent(s, "", "  ")
	return b
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

// Query for hosts.
func (r *SubnetResource) Query(cond map[string]string) (*[]Subnet, error) {
	c := *r.Client
	j, err := c.Query("/subnets", cond)
	if err != nil {
		return nil, err
	}

	subnets := &[]Subnet{}
	if err := json.Unmarshal(j, subnets); err != nil {
		return nil, err
	}

	return subnets, nil
}

// Get subnet.
func (r *SubnetResource) Get(id string) (*Subnet, error) {
	c := *r.Client
	j, err := c.Get("/subnets", id)
	if err != nil {
		return nil, err
	}

	subnet := &Subnet{}
	if err := json.Unmarshal(j, subnet); err != nil {
		return nil, err
	}

	return subnet, nil
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

// Update subnet.
func (r *SubnetResource) Update(id string, s *Subnet) (*Subnet, error) {
	c := *r.Client
	j, err := c.Update("/subnets/", id, s)
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
func (r *SubnetResource) Delete(id string) (*Subnet, error) {
	c := *r.Client
	j, err := c.Delete("/subnets", id)
	if err != nil {
		return nil, err
	}

	subnet := &Subnet{}
	if err := json.Unmarshal(j, subnet); err != nil {
		return nil, err
	}

	return subnet, nil
}
