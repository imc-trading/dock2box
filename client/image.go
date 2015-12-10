package client

import "encoding/json"

// ImageResource structure.
type ImageResource struct {
	Client *Client
}

// Image structure.
type Image struct {
	ID           string         `json:"id,omitempty"`
	Image        string         `json:"image,omitempty"`
	Type         string         `json:"type,omitempty"`
	BootImageID  string         `json:"bootImageId,omitempty"`
	BootImageRef string         `json:"bootImageRef,omitempty"`
	BootImage    string         `json:"bootImage,omitempty"`
	Versions     []ImageVersion `json:"versions,omitempty"`
}

// Create image.
func (r *ImageResource) Create(i *Image) (*Image, error) {
	c := *r.Client
	j, err := c.Create("/images", i)
	if err != nil {
		return nil, err
	}

	image := &Image{}
	if err := json.Unmarshal(j, image); err != nil {
		return nil, err
	}

	return image, nil
}

// Get image.
func (r *ImageResource) Get(name string) (*Image, error) {
	c := *r.Client
	j, err := c.Get("/images", name)
	if err != nil {
		return nil, err
	}

	image := &Image{}
	if err := json.Unmarshal(j, image); err != nil {
		return nil, err
	}

	return image, nil
}

// All images.
func (r *ImageResource) All() (*[]Image, error) {
	c := *r.Client
	j, err := c.All("/images")
	if err != nil {
		return nil, err
	}

	images := &[]Image{}
	if err := json.Unmarshal(j, images); err != nil {
		return nil, err
	}

	return images, nil
}
