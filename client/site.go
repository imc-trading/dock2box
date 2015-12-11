package client

import "encoding/json"

// SiteResource structure.
type SiteResource struct {
	Client *Client
}

// Site structure.
type Site struct {
	ID                 string   `json:"id"`
	Site               string   `json:"site"`
	Domain             string   `json:"domain"`
	DNS                []string `json:"dns"`
	DockerRegistry     string   `json:"dockerRegistry"`
	ArtifactRepository string   `json:"artifactRepository"`
	NamingScheme       string   `json:"namingScheme"`
}

// JSON output for a site.
func (s *Site) JSON() []byte {
	b, _ := json.MarshalIndent(s, "", "  ")
	return b
}

// Create site.
func (r *SiteResource) Create(s *Site) (*Site, error) {
	c := *r.Client
	j, err := c.Create("/sites", s)
	if err != nil {
		return nil, err
	}

	site := &Site{}
	if err := json.Unmarshal(j, site); err != nil {
		return nil, err
	}

	return site, nil
}

// Delete site.
func (r *SiteResource) Delete(name string) (*Site, error) {
	c := *r.Client
	j, err := c.Delete("/sites", name)
	if err != nil {
		return nil, err
	}

	site := &Site{}
	if err := json.Unmarshal(j, site); err != nil {
		return nil, err
	}

	return site, nil
}

// Get site.
func (r *SiteResource) Get(name string) (*Site, error) {
	c := *r.Client
	j, err := c.Get("/sites", name)
	if err != nil {
		return nil, err
	}

	site := &Site{}
	if err := json.Unmarshal(j, site); err != nil {
		return nil, err
	}

	return site, nil
}

// All sites.
func (r *SiteResource) All() (*[]Site, error) {
	c := *r.Client
	j, err := c.All("/sites")
	if err != nil {
		return nil, err
	}

	sites := &[]Site{}
	if err := json.Unmarshal(j, sites); err != nil {
		return nil, err
	}

	return sites, nil
}
