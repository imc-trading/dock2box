package client

import "encoding/json"

// ImageVersionResource structure.
type ImageVersionResource struct {
	Client *Client
}

// ImageVersion structure.
type ImageVersion struct {
	Version string `json:"version,omitempty"`
	Created string `json:"created,omitempty"`
}

// All versions.
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

// AllByID versions.
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
