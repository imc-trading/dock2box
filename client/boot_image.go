package client

import "encoding/json"

// BootImageResource structure.
type BootImageResource struct {
	Client *Client
}

// BootImage structure.
type BootImage struct {
	ID       string             `json:"id"`
	Image    string             `json:"image"`
	KOpts    string             `json:"kOpts"`
	Versions []BootImageVersion `json:"versions"`
}

// BootImageVersion structure.
type BootImageVersion struct {
	Version string `json:"version"`
	Created string `json:"created"`
}

// JSON output for a boot image.
func (i *BootImage) JSON() []byte {
	b, _ := json.MarshalIndent(i, "", "  ")
	return b
}

// Create boot image.
func (r *BootImageResource) Create(s *BootImage) (*BootImage, error) {
	c := *r.Client
	j, err := c.Create("boot-images", s)
	if err != nil {
		return nil, err
	}

	bootImage := &BootImage{}
	if err := json.Unmarshal(j, bootImage); err != nil {
		return nil, err
	}

	return bootImage, nil
}

// Delete boot image.
func (r *BootImageResource) Delete(name string) (*BootImage, error) {
	c := *r.Client
	j, err := c.Delete("/boot-images", name)
	if err != nil {
		return nil, err
	}

	bootImage := &BootImage{}
	if err := json.Unmarshal(j, bootImage); err != nil {
		return nil, err
	}

	return bootImage, nil
}

// Get boot image.
func (r *BootImageResource) Get(name string) (*BootImage, error) {
	c := *r.Client
	j, err := c.Get("/boot-images", name)
	if err != nil {
		return nil, err
	}

	bootImage := &BootImage{}
	if err := json.Unmarshal(j, bootImage); err != nil {
		return nil, err
	}

	return bootImage, nil
}

// All boot images.
func (r *BootImageResource) All() (*[]BootImage, error) {
	c := *r.Client
	j, err := c.All("/boot-images")
	if err != nil {
		return nil, err
	}

	bootImages := &[]BootImage{}
	if err := json.Unmarshal(j, bootImages); err != nil {
		return nil, err
	}

	return bootImages, nil
}
