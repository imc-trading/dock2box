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

## Return Code

Method | Code | Description
--- | --- | ---
GET | 200 | OK
GET | 404 | Not Found
GET | 500 | Internal Server Error, something failed in the database get request

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

Field | Type | Description
--- | --- | ---
build | boolean | If host should be provisioned when PXE booting
dhcp | dir | Directory with embeded data
legacynet | boolean |
debug | boolean | Debug output during provisioning and doesn't reboot after provisioning is done
gpt | boolean | Use GUID Partition Table
volmgt | enum | Which volume manager to use (lvm or btrfs)
image | string | Name of host image
version | string | Version of host image
interface | Directory with embeded data
kexec | string | KExec into kernel without a reboot, this is not as fool-proof as a reboot but faster
kopts | string | Kernel options
site | string | Name of site like a datacenter or location
