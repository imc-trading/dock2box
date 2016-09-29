# Structure

Directory | Description
--------- | -----------
base-images | Base images built using Packer and then imported into Docker
boot-images | Boot images used for PXE boot
build-images | Build image used during builds

# Authentication

For the base images "root" is locked by default. You can use the following default credentials:

**Default user:** dock2box

**Default password:** D0ck2B0x

# Kernel options

Option | Type | Description
------ | ---- | -----------
quiet | generic | Boot kernel in quiet mode and use progress bar
video | generic | Set uvesafb video mode ex. 800x600-32
dma | generic | Disable DMA since it can cause issues
modules=module:module | generic | Load kernel modules
blacklist=module:module | generic | Blacklist kernel modules
install | generic | Install OS, this will wipe disks so do a backup first
sshkey=AAAAB3NzaC1y... | security | Add SSH key to "authorized_keys" for "dock2box" user
distro=distribution | image | Distribution being provisioned centos, fedora or ubuntu
image=repo/image | image | Host image name
tag=tag | image | Host image tag, defaults to "latest"
registry=registry:port | image | Docker registry, defaults to Docker Hub
hostname=hostname | network | Short hostname
ip=ip:netmask:gw | network | Static IP, defaults to DHCP
interface=name | network | Network interface, defaults to first interface
dns=dns1:dns2:search | network | DNS configuration
gpt | disk | Use GPT partitions
root_size=size | disk | size of "/root", ex. "10G"
var_size=size | disk | Size of "/var", ex. "15G"
swap_size=size | disk | Size of "swap", ex. "4G"

# Example iPXE script

```
#!ipxe

kernel http://my-server/dock2box/boot-images/alpine3.4.3/kernel distro=centos image=dock2box/centos7.2.1511 install gpt quiet video=800x600-32
initrd http://my-server/dock2box/boot-images/alpine3.4.3/initrd
boot
```
