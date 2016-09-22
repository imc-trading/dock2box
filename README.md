# Kernel options

Option | Type | Description
------ | ---- | -----------
quiet | generic | Boot kernel in quiet mode
debug | generic | Echo all commands in init to the Console
dma | generic | Disable DMA since it can cause issues
modules=module:module | generic | Load kernel modules
blacklist=module:module | generic | Blacklist kernel modules
install | generic | Install OS, this will wipe disks so do a backup first
sshkey= AAAAB3NzaC1y... | security | Add SSH key to "authorized_keys" for "dock2box" user
distro=distribution | image | Distribution being provisioned centos, fedora or ubuntu
image=repo/image | image | Host image name
tag=tag | image | Host image tag, defaults to "latest"
registry=registry:port | image | Docker registry, defaults to Docker Hub
hostname=hostname | network | Short hostname
ip=ip:netmask:gw | network | Static IP, defaults to DHCP
interface=eth0 | network | Network interface, defaults to first interface
dns=dns1:dns2:search | network | DNS configuration
gpt | disk | Use GPT partitions
root_size | disk | size of "/root", ex. "10G"
var_size | disk | Size of "/var", ex. "15G"
swap_size | disk | Size of "swap", ex. "4G"
