# Dock2Box Go Client

Go client for Dock2Box.

[![GoDoc](https://godoc.org/github.com/imc-trading/dock2box/client?status.svg)](https://godoc.org/github.com/imc-trading/dock2box/client)

# Documentation


# client
    import "github.com/imc-trading/dock2box/client"







## type Client
``` go
type Client struct {
    URL       string
    Host      HostResource
    Interface InterfaceResource
    Image     ImageResource
    Tag       TagResource
    Site      SiteResource
    Tenant    TenantResource
    Subnet    SubnetResource
    Debug     bool
}
```
Client structure.









### func New
``` go
func New(url string) *Client
```
New client.




### func (\*Client) All
``` go
func (c *Client) All(endp string) ([]byte, error)
```
All resources.



### func (\*Client) Create
``` go
func (c *Client) Create(endp string, s interface{}) ([]byte, error)
```
Create resource.



### func (\*Client) Delete
``` go
func (c *Client) Delete(endp string, id string) ([]byte, error)
```
Delete resource.



### func (\*Client) Fatal
``` go
func (c *Client) Fatal(msg string)
```
Fatal log and exit



### func (\*Client) Fatalf
``` go
func (c *Client) Fatalf(fmt string, args ...interface{})
```
Fatalf log and exit



### func (\*Client) Get
``` go
func (c *Client) Get(endp string, id string) ([]byte, error)
```
Get resource.



### func (\*Client) Info
``` go
func (c *Client) Info(msg string)
```
Info log



### func (\*Client) Infof
``` go
func (c *Client) Infof(fmt string, args ...interface{})
```
Infof log



### func (\*Client) Query
``` go
func (c *Client) Query(endp string, cond map[string]string) ([]byte, error)
```
Query for resources.



### func (\*Client) SetDebug
``` go
func (c *Client) SetDebug()
```
SetDebug enable debug.



### func (\*Client) Update
``` go
func (c *Client) Update(endp string, id string, s interface{}) ([]byte, error)
```
Update resource.



## type Host
``` go
type Host struct {
    ID       string   `json:"id"`
    Host     string   `json:"host"`
    Build    bool     `json:"build"`
    Debug    bool     `json:"debug"`
    GPT      bool     `json:"gpt"`
    TagID    string   `json:"tagId"`
    KOpts    string   `json:"kOpts"`
    TenantID string   `json:"tenantId"`
    Labels   []string `json:"labels"`
    SiteID   string   `json:"siteId"`
}
```
Host structure.











### func (\*Host) JSON
``` go
func (h *Host) JSON() []byte
```
JSON output for a host.



## type HostResource
``` go
type HostResource struct {
    Client *Client
}
```
HostResource structure.











### func (\*HostResource) All
``` go
func (r *HostResource) All() (*[]Host, error)
```
All hosts.



### func (\*HostResource) Create
``` go
func (r *HostResource) Create(h *Host) (*Host, error)
```
Create host.



### func (\*HostResource) Delete
``` go
func (r *HostResource) Delete(id string) (*Host, error)
```
Delete host.



### func (\*HostResource) Get
``` go
func (r *HostResource) Get(id string) (*Host, error)
```
Get host.



### func (\*HostResource) Query
``` go
func (r *HostResource) Query(cond map[string]string) (*[]Host, error)
```
Query for hosts.



### func (\*HostResource) Update
``` go
func (r *HostResource) Update(id string, h *Host) (*Host, error)
```
Update host.



## type Image
``` go
type Image struct {
    ID        string `json:"id"`
    Image     string `json:"image"`
    Type      string `json:"type"`
    KOpts     string `json:"kOpts,omitempty"`
    BootTagID string `json:"bootTagId,omitempty"`
}
```
Image structure.











### func (\*Image) JSON
``` go
func (i *Image) JSON() []byte
```
JSON output for a image.



## type ImageResource
``` go
type ImageResource struct {
    Client *Client
}
```
ImageResource structure.











### func (\*ImageResource) All
``` go
func (r *ImageResource) All() (*[]Image, error)
```
All images.



### func (\*ImageResource) Create
``` go
func (r *ImageResource) Create(i *Image) (*Image, error)
```
Create image.



### func (\*ImageResource) Delete
``` go
func (r *ImageResource) Delete(id string) (*Image, error)
```
Delete image.



### func (\*ImageResource) Get
``` go
func (r *ImageResource) Get(id string) (*Image, error)
```
Get image.



### func (\*ImageResource) Query
``` go
func (r *ImageResource) Query(cond map[string]string) (*[]Image, error)
```
Query for image.



### func (\*ImageResource) Update
``` go
func (r *ImageResource) Update(id string, i *Image) (*Image, error)
```
Update image.



## type Interface
``` go
type Interface struct {
    ID        string `json:"id"`
    Interface string `json:"interface"`
    DHCP      bool   `json:"dhcp"`
    IPv4      string `json:"ipv4,omitempty"`
    HwAddr    string `json:"hwAddr"`
    SubnetID  string `json:"subnetId,omitempty"`
    HostID    string `json:"hostId"`
}
```
Interface structure.











### func (\*Interface) JSON
``` go
func (i *Interface) JSON() []byte
```
JSON output for a interface.



## type InterfaceResource
``` go
type InterfaceResource struct {
    Client *Client
}
```
InterfaceResource structure.











### func (\*InterfaceResource) All
``` go
func (r *InterfaceResource) All() (*[]Interface, error)
```
All interfaces.



### func (\*InterfaceResource) Create
``` go
func (r *InterfaceResource) Create(h *Interface) (*Interface, error)
```
Create interface.



### func (\*InterfaceResource) Delete
``` go
func (r *InterfaceResource) Delete(id string) (*Interface, error)
```
Delete interface.



### func (\*InterfaceResource) Get
``` go
func (r *InterfaceResource) Get(id string) (*Interface, error)
```
Get interface.



### func (\*InterfaceResource) Query
``` go
func (r *InterfaceResource) Query(cond map[string]string) (*[]Interface, error)
```
Query for interfaces.



### func (\*InterfaceResource) Update
``` go
func (r *InterfaceResource) Update(id string, h *Interface) (*Interface, error)
```
Update interface.



## type Site
``` go
type Site struct {
    ID                 string   `json:"id"`
    Site               string   `json:"site"`
    Domain             string   `json:"domain"`
    DNS                []string `json:"dns"`
    DockerRegistry     string   `json:"dockerRegistry"`
    ArtifactRepository string   `json:"artifactRepository"`
    NamingScheme       string   `json:"namingScheme"`
    PXETheme           string   `json:"pxeTheme"`
}
```
Site structure.











### func (\*Site) JSON
``` go
func (s *Site) JSON() []byte
```
JSON output for a site.



## type SiteResource
``` go
type SiteResource struct {
    Client *Client
}
```
SiteResource structure.











### func (\*SiteResource) All
``` go
func (r *SiteResource) All() (*[]Site, error)
```
All sites.



### func (\*SiteResource) Create
``` go
func (r *SiteResource) Create(s *Site) (*Site, error)
```
Create site.



### func (\*SiteResource) Delete
``` go
func (r *SiteResource) Delete(id string) (*Site, error)
```
Delete site.



### func (\*SiteResource) Get
``` go
func (r *SiteResource) Get(id string) (*Site, error)
```
Get site.



### func (\*SiteResource) Query
``` go
func (r *SiteResource) Query(cond map[string]string) (*[]Site, error)
```
Query for sites.



### func (\*SiteResource) Update
``` go
func (r *SiteResource) Update(id string, s *Site) (*Site, error)
```
Update site.



## type Subnet
``` go
type Subnet struct {
    ID     string `json:"id"`
    Subnet string `json:"subnet"`
    Mask   string `json:"mask"`
    Gw     string `json:"gw"`
    SiteID string `json:"siteId"`
}
```
Subnet structure.











### func (\*Subnet) JSON
``` go
func (s *Subnet) JSON() []byte
```
JSON output for a subnet.



## type SubnetResource
``` go
type SubnetResource struct {
    Client *Client
}
```
SubnetResource structure.











### func (\*SubnetResource) All
``` go
func (r *SubnetResource) All() (*[]Subnet, error)
```
All subnets.



### func (\*SubnetResource) Create
``` go
func (r *SubnetResource) Create(s *Subnet) (*Subnet, error)
```
Create subnet.



### func (\*SubnetResource) Delete
``` go
func (r *SubnetResource) Delete(id string) (*Subnet, error)
```
Delete subnet.



### func (\*SubnetResource) Get
``` go
func (r *SubnetResource) Get(id string) (*Subnet, error)
```
Get subnet.



### func (\*SubnetResource) Query
``` go
func (r *SubnetResource) Query(cond map[string]string) (*[]Subnet, error)
```
Query for hosts.



### func (\*SubnetResource) Update
``` go
func (r *SubnetResource) Update(id string, s *Subnet) (*Subnet, error)
```
Update subnet.



## type Tag
``` go
type Tag struct {
    ID      string `json:"id"`
    Tag     string `json:"tag"`
    Created string `json:"created"`
    SHA256  string `json:"sha256"`
    ImageID string `json:"imageId"`
}
```
Tag structure.











### func (\*Tag) JSON
``` go
func (s *Tag) JSON() []byte
```
JSON output for a tag.



## type TagResource
``` go
type TagResource struct {
    Client *Client
}
```
TagResource structure.











### func (\*TagResource) All
``` go
func (r *TagResource) All() (*[]Tag, error)
```
All tags.



### func (\*TagResource) Create
``` go
func (r *TagResource) Create(s *Tag) (*Tag, error)
```
Create tag.



### func (\*TagResource) Delete
``` go
func (r *TagResource) Delete(id string) (*Tag, error)
```
Delete tag.



### func (\*TagResource) Get
``` go
func (r *TagResource) Get(id string) (*Tag, error)
```
Get tag.



### func (\*TagResource) Query
``` go
func (r *TagResource) Query(cond map[string]string) (*[]Tag, error)
```
Query for tags.



### func (\*TagResource) Update
``` go
func (r *TagResource) Update(id string, s *Tag) (*Tag, error)
```
Update tag.



## type Tenant
``` go
type Tenant struct {
    ID     string `json:"id"`
    Tenant string `json:"tenant"`
}
```
Tenant structure.











### func (\*Tenant) JSON
``` go
func (t *Tenant) JSON() []byte
```
JSON output for a tenant.



## type TenantResource
``` go
type TenantResource struct {
    Client *Client
}
```
TenantResource structure.











### func (\*TenantResource) All
``` go
func (r *TenantResource) All() (*[]Tenant, error)
```
All tenants.



### func (\*TenantResource) Create
``` go
func (r *TenantResource) Create(s *Tenant) (*Tenant, error)
```
Create tenant.



### func (\*TenantResource) Delete
``` go
func (r *TenantResource) Delete(id string) (*Tenant, error)
```
Delete tenant.



### func (\*TenantResource) Get
``` go
func (r *TenantResource) Get(id string) (*Tenant, error)
```
Get tenant.



### func (\*TenantResource) Query
``` go
func (r *TenantResource) Query(cond map[string]string) (*[]Tenant, error)
```
Query for hosts.



### func (\*TenantResource) Update
``` go
func (r *TenantResource) Update(id string, s *Tenant) (*Tenant, error)
```
Update tenant.









- - -
