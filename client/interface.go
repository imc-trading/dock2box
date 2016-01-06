package client

import "encoding/json"

// InterfaceResource structure.
type InterfaceResource struct {
	Client *Client
}

// Interface structure.
type Interface struct {
	ID        string `json:"id"`
	Interface string `json:"interface"`
	DHCP      bool   `json:"dhcp"`
	IPv4      string `json:"ipv4,omitempty"`
	HwAddr    string `json:"hwAddr"`
	SubnetID  string `json:"subnetId,omitempty"`
	HostID    string `json:"hostId"`
}

// JSON output for a interface.
func (i *Interface) JSON() []byte {
	b, _ := json.MarshalIndent(i, "", "  ")
	return b
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
func (r *InterfaceResource) Update(name string, h *Interface) (*Interface, error) {
	c := *r.Client
	j, err := c.Update("/interfaces/"+name, h)
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
func (r *InterfaceResource) Delete(name string) (*Interface, error) {
	c := *r.Client
	j, err := c.Delete("/interfaces", name)
	if err != nil {
		return nil, err
	}

	intf := &Interface{}
	if err := json.Unmarshal(j, intf); err != nil {
		return nil, err
	}

	return intf, nil
}

// Get interface.
func (r *InterfaceResource) Get(name string) (*Interface, error) {
	c := *r.Client
	j, err := c.Get("/interfaces", name)
	if err != nil {
		return nil, err
	}

	intf := &Interface{}
	if err := json.Unmarshal(j, intf); err != nil {
		return nil, err
	}

	return intf, nil
}

// GetByID interface.
func (r *InterfaceResource) GetByID(id string) (*Interface, error) {
	c := *r.Client
	j, err := c.Get("/interfaces/id", id)
	if err != nil {
		return nil, err
	}

	intf := &Interface{}
	if err := json.Unmarshal(j, intf); err != nil {
		return nil, err
	}

	return intf, nil
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
