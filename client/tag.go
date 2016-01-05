package client

import "encoding/json"

// TagResource structure.
type TagResource struct {
	Client *Client
}

// Tag structure.
type Tag struct {
	ID      string `json:"id"`
	Tag     string `json:"tag"`
	Created string `json:"created"`
	SHA256  string `json:"sha256"`
	ImageID string `json:"imageId"`
}

// JSON output for a tag.
func (s *Tag) JSON() []byte {
	b, _ := json.MarshalIndent(s, "", "  ")
	return b
}

// Create tag.
func (r *TagResource) Create(s *Tag) (*Tag, error) {
	c := *r.Client
	j, err := c.Create("/tags", s)
	if err != nil {
		return nil, err
	}

	tag := &Tag{}
	if err := json.Unmarshal(j, tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// Update tag.
func (r *TagResource) Update(name string, s *Tag) (*Tag, error) {
	c := *r.Client
	j, err := c.Update("/tags/"+name, s)
	if err != nil {
		return nil, err
	}

	tag := &Tag{}
	if err := json.Unmarshal(j, tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// Delete tag.
func (r *TagResource) Delete(name string) (*Tag, error) {
	c := *r.Client
	j, err := c.Delete("/tags", name)
	if err != nil {
		return nil, err
	}

	tag := &Tag{}
	if err := json.Unmarshal(j, tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// Get tag.
func (r *TagResource) Get(name string) (*Tag, error) {
	c := *r.Client
	j, err := c.Get("/tags", name)
	if err != nil {
		return nil, err
	}

	tag := &Tag{}
	if err := json.Unmarshal(j, tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// All tags.
func (r *TagResource) All() (*[]Tag, error) {
	c := *r.Client
	j, err := c.All("/tags")
	if err != nil {
		return nil, err
	}

	tags := &[]Tag{}
	if err := json.Unmarshal(j, tags); err != nil {
		return nil, err
	}

	return tags, nil
}
