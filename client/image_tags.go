package client

import "encoding/json"

// ImageTagResource structure.
type ImageTagResource struct {
	Client *Client
}

// ImageTag structure.
type ImageTag struct {
	ID      string `json:"id"`
	Tag     string `json:"tag"`
	Created string `json:"created"`
	SHA256  string `json:"sha256"`
	ImageID string `json:"imageId"`
}

// JSON output for a image tag.
func (s *ImageTag) JSON() []byte {
	b, _ := json.MarshalIndent(s, "", "  ")
	return b
}

// Create image tag.
func (r *ImageTagResource) Create(s *ImageTag) (*ImageTag, error) {
	c := *r.Client
	j, err := c.Create("/image-tags", s)
	if err != nil {
		return nil, err
	}

	imageTag := &ImageTag{}
	if err := json.Unmarshal(j, imageTag); err != nil {
		return nil, err
	}

	return imageTag, nil
}

// Update image tag.
func (r *ImageTagResource) Update(name string, s *ImageTag) (*ImageTag, error) {
	c := *r.Client
	j, err := c.Update("/image-tags/"+name, s)
	if err != nil {
		return nil, err
	}

	imageTag := &ImageTag{}
	if err := json.Unmarshal(j, imageTag); err != nil {
		return nil, err
	}

	return imageTag, nil
}

// Delete image tag.
func (r *ImageTagResource) Delete(name string) (*ImageTag, error) {
	c := *r.Client
	j, err := c.Delete("/image-tags", name)
	if err != nil {
		return nil, err
	}

	imageTag := &ImageTag{}
	if err := json.Unmarshal(j, imageTag); err != nil {
		return nil, err
	}

	return imageTag, nil
}

// Get image tag.
func (r *ImageTagResource) Get(name string) (*ImageTag, error) {
	c := *r.Client
	j, err := c.Get("/image-tags", name)
	if err != nil {
		return nil, err
	}

	imageTag := &ImageTag{}
	if err := json.Unmarshal(j, imageTag); err != nil {
		return nil, err
	}

	return imageTag, nil
}

// All image tags.
func (r *ImageTagResource) All() (*[]ImageTag, error) {
	c := *r.Client
	j, err := c.All("/image-tags")
	if err != nil {
		return nil, err
	}

	imageTags := &[]ImageTag{}
	if err := json.Unmarshal(j, imageTags); err != nil {
		return nil, err
	}

	return imageTags, nil
}
