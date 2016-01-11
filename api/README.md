# Usage

```
Usage of ./d2bsrv:
  -base-uri string
    	Base URI for server
  -bind string
    	Bind to address and port (default "127.0.0.1:8080")
  -database string
    	Database name (default "d2b")
  -disable-hateoas
    	Disable HATEOAS per default
  -enable-envelope
    	Enable Envelopes per default
  -schema-uri string
    	URI to JSON schemas (default "file://schemas")
  -version
    	Version
```

# Filtering

## Query

To make queries just add the field and value, you can specify one or several fields.

> At the moment it only support's equal to.

**Example:**

```
/v1/tags?tag=latest&imageId=568d2ba85d099040397ae363
```

## fields

This allows you to specify which fields you want in the result.

**Example:**
```
/v1/tags?fields=id,tag
```

## sort

This allows you to sort the result ascending or descending, you can specify one or several fields.

**Example ascending:**
```
/v1/tags?sort=tag,created
```

You can sort descending by adding a minus sign.

**Example descending:**
```
/v1/tags?sort=tag,-created
```

# Options

## envelope

This enabled/disables embedding data in an envelope with additional info that normally is available as HTTP status code.

**Example:**
```
/v1/tags?envelope=true
```

**Example output:**
```
{
  "code": 200,
  "data": [
    ...
  ]
}
```

> This can be enabled globally by using **-enable-envelope** when starting the server.

## hateoas

This enables/disables HATEOAS which includes links to methods and related endpoints.

**Example:**
```
/v1/tags?hateoas=true
```

**Example output:**
```
"links": [
  {
    "href": "http://yggdrasil.trading.imc.intra:8080/v1/images/568d2ba85d099040397ae363",
    "rel": "self",
    "method": "GET"
  }
]
```

> This can be disabled globally by using *-disable-hateoas* when starting the server.

## embed

This enables/disables embedding related data in the result. This will affect perfomance since it has the server has to do additional queries.

**Example:**
```
/v1/tags?embed=true
```

**Example output:**
```
{
  "id": "568d2ba85d099040397ae365",
  "tag": "untested",
  "created": "2015-12-01T13:01:05Z",
  "sha256": "37ff8e2ae04a1570781a63a247fce789352beae2889f1d720b2efbec50ef8e0d",
  "imageId": "568d2ba85d099040397ae363",
  "image": {
    "id": "568d2ba85d099040397ae363",
    "image": "test2",
    "type": "docker",
    "bootTagId": "568d2ba85d099040397ae362"
  }
}
```

# Data model

## Host

Field | Type | Description
--- | --- | ---
id | Unique id | Unique Id for host
host | String | Hostname
build | Boolean | If host should be provisioned when PXE booting
debug | Boolean | Debug info and no reboot after provisioning finishes
gpt | Boolean | Use GUID Partition Table
tagId | Foreign key | Image and tag
kOpts | String | Kernel options
tenantId | Foreign key | Tenant
labels | List of strings | Labels for host
siteId | Foreign key | Site

## Interface

Field | Type | Description
--- | --- | ---
id | Unique id | Unique Id for interface
interface | String | Interface name
dhcp | Boolean | Use DHCP
ipv4 | String | IP address (only for fixed IP address)
hwAddr | String | Hardware address
subnetId | Foreign key | Subnet (only for fixed IP address)
hostId | Foreign key | Host

## Image

Field | Type | Description
--- | --- | ---
id | Unique id | Unique Id for interface
image | String | Image name
type | String | Image type dockera ,file or boot
kOpts | String | Kernel options (only for boot image)
bootTagId | Foreign key | Boot image and tag (only for boot image)

## Tag

Field | Type | Description
--- | --- | ---
id | Unique id | Unique Id for tag
tag | String | Tag name
created | Date-time | When tag was created
sha256 | String | SHA256 checksum of image
imageId | Foreign key | Image

## Site

Field | Type | Description
--- | --- | ---
id | Unique id | Unique Id for interface
site | String | Site name
domain | String Domain name
dns | List of strings | DNS servers
dockerRegistry | String | Docker registry
ArtifactRepository | String | Artifact repository for file images and scripts
pxeTheme | String | Theme for PXE boot menu
namingScheme | String | Dynamic naming scheme for unregistered hosts

## Subnet

Field | Type | Description
--- | --- | ---
id | Unique id | Unique Id for interface
subnet | String | Subnet IP address and Prefix
mask | String | Mask
gw | String | Gateway
siteId | String | Site

## Tenant

Field | Type | Description
--- | --- | ---
id | Unique id | Unique Id for interface
tenant | String | Tenant name

# ROADMAP

- Add checks for foreign keys when creating/updating an entry
- Add more complex queries
- Refactor boiler plate code in each controller and normalize it
- Rename upper-case ID to Id
- Add direct support for SSL and auth. for now rely on a proxy
