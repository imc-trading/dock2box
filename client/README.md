# Dock2Box Go Client

Go client for Dock2Box.

[![GoDoc](https://godoc.org/github.com/imc-trading/dock2box/client?status.svg)](https://godoc.org/github.com/imc-trading/dock2box/client)

# Documentation


# client
    import "github.com/imc-trading/dock2box/client"







## type BootImage
``` go
type BootImage struct {
    ID       string             `json:"id,omitempty"`
    Image    string             `json:"image"`
    KOpts    string             `json:"kOpts"`
    Versions []BootImageVersion `json:"versions,omitempty"`
}
```
BootImage structure.











### func (\*BootImage) JSON
``` go
func (i *BootImage) JSON() []byte
```
JSON output for a boot image.



## type BootImageResource
``` go
type BootImageResource struct {
    Client *Client
}
```
BootImageResource structure.











### func (\*BootImageResource) All
``` go
func (r *BootImageResource) All() (*[]BootImage, error)
```
All boot images.



### func (\*BootImageResource) Create
``` go
func (r *BootImageResource) Create(s *BootImage) (*BootImage, error)
```
Create boot image.



### func (\*BootImageResource) Delete
``` go
func (r *BootImageResource) Delete(name string) (*BootImage, error)
```
Delete boot image.



### func (\*BootImageResource) Get
``` go
func (r *BootImageResource) Get(name string) (*BootImage, error)
```
Get boot image.



### func (\*BootImageResource) Update
``` go
func (r *BootImageResource) Update(name string, s *BootImage) (*BootImage, error)
```
Update boot image.



## type BootImageVersion
``` go
type BootImageVersion struct {
    Version string `json:"version"`
    Created string `json:"created"`
}
```
BootImageVersion structure.











## type Client
``` go
type Client struct {
    URL          string
    Host         HostResource
    Image        ImageResource
    ImageVersion ImageVersionResource
    Site         SiteResource
    Tenant       TenantResource
    Subnet       SubnetResource
    BootImage    BootImageResource
    Debug        bool
}
```
Client structure.









### func New
``` go
func New(url string) *Client
```
New client.




### func (Client) All
``` go
func (c Client) All(endp string) ([]byte, error)
```
All resources.



### func (Client) Create
``` go
func (c Client) Create(endp string, s interface{}) ([]byte, error)
```
Create resource.



### func (Client) Delete
``` go
func (c Client) Delete(endp string, name string) ([]byte, error)
```
Delete resource.



### func (Client) Exist
``` go
func (c Client) Exist(endp string, name string) (bool, error)
```
Exist resource.



### func (Client) Fatal
``` go
func (c Client) Fatal(msg string)
```
Fatal log and exit



### func (Client) Fatalf
``` go
func (c Client) Fatalf(fmt string, args ...interface{})
```
Fatalf log and exit



### func (Client) Get
``` go
func (c Client) Get(endp string, name string) ([]byte, error)
```
Get resource.



### func (Client) Info
``` go
func (c Client) Info(msg string)
```
Info log



### func (Client) Infof
``` go
func (c Client) Infof(fmt string, args ...interface{})
```
Infof log



### func (Client) SetDebug
``` go
func (c Client) SetDebug()
```
SetDebug enable debug.



### func (Client) Update
``` go
func (c Client) Update(endp string, s interface{}) ([]byte, error)
```
Update resource.



## type Host
``` go
type Host struct {
    Host       string          `json:"host"`
    Build      bool            `json:"build"`
    Debug      bool            `json:"debug"`
    GPT        bool            `json:"gpt"`
    ImageID    string          `json:"imageId"`
    Version    string          `json:"version"`
    KOpts      string          `json:"kOpts"`
    TenantID   string          `json:"tenantId"`
    Labels     []string        `json:"labels"`
    SiteID     string          `json:"siteId"`
    Interfaces []HostInterface `json:"interfaces,omitempty"`
}
```
Host structure.











### func (\*Host) JSON
``` go
func (h *Host) JSON() []byte
```
JSON output for a host.



## type HostInterface
``` go
type HostInterface struct {
    Interface string `json:"interface"`
    DHCP      bool   `json:"dhcp"`
    IPv4      string `json:"ipv4,omitempty"`
    HwAddr    string `json:"hwAddr"`
    SubnetID  string `json:"subnetId,omitempty"`
}
```
HostInterface structure.











### func (\*HostInterface) JSON
``` go
func (i *HostInterface) JSON() []byte
```
JSON output for a host interface.



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
func (r *HostResource) Delete(name string) (*Host, error)
```
Delete host.



### func (\*HostResource) Exist
``` go
func (r *HostResource) Exist(name string) (bool, error)
```
Exist host.



### func (\*HostResource) Get
``` go
func (r *HostResource) Get(name string) (*Host, error)
```
Get host.



### func (\*HostResource) GetByID
``` go
func (r *HostResource) GetByID(id string) (*Host, error)
```
GetByID host.



### func (\*HostResource) Update
``` go
func (r *HostResource) Update(name string, h *Host) (*Host, error)
```
Update host.



## type Image
``` go
type Image struct {
    ID           string         `json:"id,omitempty"`
    Image        string         `json:"image,omitempty"`
    Type         string         `json:"type,omitempty"`
    BootImageID  string         `json:"bootImageId,omitempty"`
    BootImageRef string         `json:"bootImageRef,omitempty"`
    BootImage    string         `json:"bootImage,omitempty"`
    Versions     []ImageVersion `json:"versions,omitempty"`
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
func (r *ImageResource) Delete(name string) (*Image, error)
```
Delete image.



### func (\*ImageResource) Get
``` go
func (r *ImageResource) Get(name string) (*Image, error)
```
Get image.



### func (\*ImageResource) Update
``` go
func (r *ImageResource) Update(name string, i *Image) (*Image, error)
```
Update image.



## type ImageVersion
``` go
type ImageVersion struct {
    Version string `json:"version,omitempty"`
    Created string `json:"created,omitempty"`
}
```
ImageVersion structure.











## type ImageVersionResource
``` go
type ImageVersionResource struct {
    Client *Client
}
```
ImageVersionResource structure.











### func (\*ImageVersionResource) All
``` go
func (r *ImageVersionResource) All(name string) (*[]ImageVersion, error)
```
All versions.



### func (\*ImageVersionResource) AllByID
``` go
func (r *ImageVersionResource) AllByID(id string) (*[]ImageVersion, error)
```
AllByID versions.



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
func (r *SiteResource) Delete(name string) (*Site, error)
```
Delete site.



### func (\*SiteResource) Get
``` go
func (r *SiteResource) Get(name string) (*Site, error)
```
Get site.



### func (\*SiteResource) Update
``` go
func (r *SiteResource) Update(name string, s *Site) (*Site, error)
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
func (r *SubnetResource) Delete(name string) (*Subnet, error)
```
Delete subnet.



### func (\*SubnetResource) Get
``` go
func (r *SubnetResource) Get(name string) (*Subnet, error)
```
Get subnet.



### func (\*SubnetResource) Update
``` go
func (r *SubnetResource) Update(name string, s *Subnet) (*Subnet, error)
```
Update subnet.



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
func (r *TenantResource) Delete(name string) (*Tenant, error)
```
Delete tenant.



### func (\*TenantResource) Get
``` go
func (r *TenantResource) Get(name string) (*Tenant, error)
```
Get tenant.



### func (\*TenantResource) Update
``` go
func (r *TenantResource) Update(name string, s *Tenant) (*Tenant, error)
```
Update tenant.









- - -
