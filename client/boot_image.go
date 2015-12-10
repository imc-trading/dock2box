package client

import "encoding/json"

type BootImageResource struct {
	Client *Client
}

type BootImage struct {
	ID       string             `json:"id"`
	Image    string             `json:"image"`
	KOpts    string             `json:"kOpts"`
	Versions []BootImageVersion `json:"versions"`
}

type BootImageVersion struct {
	Version string `json:"version"`
	Created string `json:"created"`
}

// Create boot_image
func (r *BootImageResource) Create(s *BootImage) (*BootImage, error) {
	c := *r.Client
	j, err := c.Create("boot-images", s)
	if err != nil {
		return nil, err
	}

	boot_image := &BootImage{}
	if err := json.Unmarshal(j, boot_image); err != nil {
		return nil, err
	}

	return boot_image, nil
}

// Get boot_image
func (r *BootImageResource) Get(name string) (*BootImage, error) {
	c := *r.Client
	j, err := c.Get("boot-images", name)
	if err != nil {
		return nil, err
	}

	boot_image := &BootImage{}
	if err := json.Unmarshal(j, boot_image); err != nil {
		return nil, err
	}

	return boot_image, nil
}

// All boot_images
func (r *BootImageResource) All() (*[]BootImage, error) {
	c := *r.Client
	j, err := c.All("boot-images")
	if err != nil {
		return nil, err
	}

	boot_images := &[]BootImage{}
	if err := json.Unmarshal(j, boot_images); err != nil {
		return nil, err
	}

	return boot_images, nil
}
