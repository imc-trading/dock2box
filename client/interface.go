package client

import "encoding/json"

// InterfaceResource structure.
type InterfaceResource struct {
	Client *Client
}

// Interface structure.
type Interface struct {
	ID        string `json:"id,omitempty"`
	Interface string `json:"interface,omitempty"`
	DHCP      bool   `json:"dhcp"`
	IPv4      string `json:"ipv4,omitempty"`
	HwAddr    string `json:"hwAddr,omitempty"`
	SubnetID  string `json:"subnetId,omitempty"`
	HostID    string `json:"hostId,omitempty"`
}

// JSON output for a interface.
func (i *Interface) JSON() []byte {
	b, _ := json.MarshalIndent(i, "", "  ")
	return b
}

// All interfaces.
func (r *InterfaceResource) All() (*[]Interface, error) {
	c := *r.Client
	j, err := c.All("/interfaces")
	if err != nil {
		return nil, err
	}

	intfs := &[]Interface{}
	if err := json.Unmarshal(j, intfs); err != nil {
		return nil, err
	}

	return intfs, nil
}

// Query for interfaces.
func (r *InterfaceResource) Query(cond map[string]string) (*[]Interface, error) {
	c := *r.Client
	j, err := c.Query("/interfaces", cond)
	if err != nil {
		return nil, err
	}

	intfs := &[]Interface{}
	if err := json.Unmarshal(j, intfs); err != nil {
		return nil, err
	}

	return intfs, nil
}

// Get interface.
func (r *InterfaceResource) Get(id string) (*Interface, error) {
	c := *r.Client
	j, err := c.Get("/interfaces", id)
	if err != nil {
		return nil, err
	}

	intf := &Interface{}
	if err := json.Unmarshal(j, intf); err != nil {
		return nil, err
	}

	return intf, nil
}

// Create interface.
func (r *InterfaceResource) Create(h *Interface) (*Interface, error) {
	c := *r.Client
	j, err := c.Create("/interfaces", h)
	if err != nil {
		return nil, err
	}

	intf := &Interface{}
	if err := json.Unmarshal(j, intf); err != nil {
		return nil, err
	}

	return intf, nil
}

// Update interface.
func (r *InterfaceResource) Update(id string, h *Interface) (*Interface, error) {
	c := *r.Client
	j, err := c.Update("/interfaces/", id, h)
	if err != nil {
		return nil, err
	}

	intf := &Interface{}
	if err := json.Unmarshal(j, intf); err != nil {
		return nil, err
	}

	return intf, nil
}

// Delete interface.
func (r *InterfaceResource) Delete(id string) (*Interface, error) {
	c := *r.Client
	j, err := c.Delete("/interfaces", id)
	if err != nil {
		return nil, err
	}

	intf := &Interface{}
	if err := json.Unmarshal(j, intf); err != nil {
		return nil, err
	}

	return intf, nil
}
