package client

import "encoding/json"

type ImageVersionResource struct {
	Client *Client
}

type ImageVersion struct {
	Version string `json:"version,omitempty"`
	Created string `json:"created,omitempty"`
}

// Create image version
/*
func (r *ImageVersionResource) Create(i *ImageVersion) (*ImageVersion, error) {
	c := *r.Client
	j, err := c.Create("/images", i)
	if err != nil {
		return nil, err
	}

	image := &ImageVersion{}
	if err := json.Unmarshal(j, image); err != nil {
		return nil, err
	}

	return image, nil
}

// Get version
func (r *ImageVersionResource) Get(name string) (*ImageVersion, error) {
	c := *r.Client
	j, err := c.Get("/images", name)
	if err != nil {
		return nil, err
	}

	image := &ImageVersion{}
	if err := json.Unmarshal(j, image); err != nil {
		return nil, err
	}

	return image, nil
}
*/

// All versions
func (r *ImageVersionResource) All(name string) (*[]ImageVersion, error) {
	c := *r.Client
	j, err := c.All("/images/" + name + "/versions")
	if err != nil {
		return nil, err
	}

	versions := &[]ImageVersion{}
	if err := json.Unmarshal(j, versions); err != nil {
		return nil, err
	}

	return versions, nil
}

// All versions
func (r *ImageVersionResource) AllByID(id string) (*[]ImageVersion, error) {
	c := *r.Client
	j, err := c.All("/images/id/" + id + "/versions")
	if err != nil {
		return nil, err
	}

	versions := &[]ImageVersion{}
	if err := json.Unmarshal(j, versions); err != nil {
		return nil, err
	}

	return versions, nil
}
