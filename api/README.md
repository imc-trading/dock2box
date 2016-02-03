# API

## API Endpoints

- /api/v1/hosts
- /api/v1/hosts/{host}/interfaces
- /api/v1/images
- /api/v1/images/{image}/versions
- /api/v1/hwaddr

## Requests

Action | Method | Description
--- | --- | ---
Create | PUT | Create new resource
Update | PUT | Update existing resource
Patch | PATCH | Patch resource using [JSON Patch](http://jsonpatch.com)
Delete | DELETE | Delete resource

> POST is not supported since we're not using unique ID's but rather each operation is idempotent

## Return Codes

Method | Code | Description
--- | --- | ---
GET | 200 | OK
GET | 404 | Not Found
GET | 500 | Internal Server Error, something failed in the database get request
PUT | 200 | OK
PUT | 400 | Bad Request, something was incorrectly formatted in your request
PUT | 500 | Internal Server Error, something failed getting document or writing it to the database
PATCH | 200 | OK
PATCH | 400 | Bad Request, something was incorrectly formatted in your requ
PATCH | 500 | Internal Server Error, something failed getting document or writing it to the database
DELETE | 200 | OK
DELETE | 404 | Not Found

## Options

### indent

This enables/disables pretty print which is enabled by default.

**Example:**
```
/v1/tags?indent=false
```
**Example output:**
```
{ "code": 200, "data": [ ... ]}
```

### envelope

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

## Data Model

### Host

Field | Require | Type | Description
--- | --- | --- | ---
build | | boolean | If host should be provisioned when PXE booting
dhcp | | dir | Directory with embeded data
legacynet | | boolean |
debug | :heavy_check_mark: | boolean | Debug output during provisioning and doesn't reboot after provisioning is done
gpt | :white_check_mark: | boolean | Use GUID Partition Table
volmgt | | enum | Which volume manager to use (lvm or btrfs)
image | :white_check_mark: | string | Name of host image
version | :white_check_mark: | string | Version of host image
interface | :white_check_mark: | dir | Directory with embeded data
kexec | :white_check_mark: | string | KExec into kernel without a reboot, this is not as fool-proof as a reboot but faster
kopts | :white_check_mark: | string | Kernel options
site | :white_check_mark: | string | Name of site like a datacenter or location

#### Host / DHCP

Field | Type | Description
--- | --- | ---
hwaddr | string | Hardware address of primary interface
ipv4 | string | DHCP IPv4 address of primary interface


