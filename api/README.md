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
Create | PUT | Create new entry
Update | PUT | Update existing entry
Patch | PATCH | JSON Schema
Delete | DELETE | 

## Return Code

Method | Code | Description
--- | --- | ---

## Options

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
