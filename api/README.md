# Start API

First start Docker and then run:

```bash
docker-compose up
```

## Update API or configuration

```bash
docker-compose stop
docker-compose pull
docker-compose up
```

## Build and push configuration

```bash
make push
```

## Export data from etcd

First install [etcdtool](https://github.com/mickep76/etcdtool).

```bash
etcdtool -p http://<docker host>:5001 export / >~/d2b.json
```

## Import data into etcd

```bash
etcdtool -p http://<docker host>:5001 import / ~/d2b.json
```

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

The current data model is somewhat messy and grew as a result of features being added. A v2 of the API is in development that should make this it more consistent.

Some values appear in multiple resources and are substituted in order, to derive the result.

### Host

Field | Required | Type | Description
--- | --- | --- | ---
build | | bool | If host should be provisioned when PXE booting
dhcp | | dir | Directory with embeded data
legacynet | | bool |
debug | :heavy_check_mark: | bool | Debug, no reboot after provisioning is done
gpt | :heavy_check_mark: | bool | Use GUID Partition Table
volmgt | | enum | Volume manager (lvm or btrfs)
image | :heavy_check_mark: | string | Name of host image
version | :heavy_check_mark: | string | Version of host image
interface | :heavy_check_mark: | dir | Directory with embeded data
kexec | :heavy_check_mark: | string | KExec into kernel without a reboot
kopts | :heavy_check_mark: | string | Kernel options
site | :heavy_check_mark: | string | Name of site like a datacenter or location

#### Host/DHCP

Field | Required | Type | Description
--- | --- | --- | ---
hwaddr | :heavy_check_mark: | string | Hardware address of primary interface
ipv4 | :heavy_check_mark: | string | DHCP IPv4 address of primary interface

#### Host/Interface

Field | Required | Type | Description
--- | --- | --- | ---
gw | :heavy_check_mark: | string | Default gateway
hwaddr | :heavy_check_mark: | string | Hardware address
ip | :heavy_check_mark: | string | IP address
netmask | :heavy_check_mark: | string | Netmask

### Image

Field | Required | Type | Description
--- | --- | --- | ---
bootstrap_image | :heavy_check_mark: | string | Bootstrap image
debug | :heavy_check_mark: | bool | Debug, no reboot after provisioning is done
docker_registry | :heavy_check_mark: | string | Docker registry
file_registry | :heavy_check_mark: | string | File registry
gpt | :heavy_check_mark: | bool | Use GUID Partition Table
kexec | :heavy_check_mark: | bool | KExec into kernel without a reboot
kopts | :heavy_check_mark: | string | Kernel options
volmgt | | enum | Volume manager (lvm or btrfs)
uefi | | bool | UEFI bios
main_script | :heavy_check_mark: | string | Main provisioning script
post_script | | string | Post provisioning script
naming_scheme | :heavy_check_mark: | string | Naming scheme
type | :heavy_check_mark: | string | Type of image
version | :heavy_check_mark: | string | Version
versions | :heavy_check_mark: | dir | Directory with embeded data

#### Image/Versions

Field | Required | Type | Description
--- | --- | --- | ---
main_script | :heavy_check_mark: | string | Main provisioning script
timestamp | :heavy_check_mark: | string | Timestamp

