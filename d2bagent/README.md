# d2bagent

Expose hardware info for Dock2Box and provide a system HTML Front-End.

**FrontEnd:**

```
http://myserver.example.com:5050
```

**JSON endpoints:**

```
/json
/cpu/json
/disks/json
/lvm/json
/lvm/log_vols/json
/lvm/phys_vols/json
/lvm/vol_grps/json
/memory/json
/mounts/json
/network/interfaces/json
/network/json
/network/json
/network/routes/json
/opsys/json
/pci/json
/sysctl/json
```

**Example:**

```bash
curl http://myserver.example.com:5050/cpu/json
```

# Usage

```bash
Usage:
  peekaboo [OPTIONS]

Application Options:
  -v, --verbose       Verbose
      --version       Version
  -b, --bind-addr=    Bind to address (0.0.0.0)
  -p, --port=         Port (5050)
  -s, --static-dir=   Static content (static)
  -t, --template-dir= Templates (templates)
  -K, --kafka         Enable Kafka message bus
      --kafka-topic=  Kafka topic (peekaboo)
      --kafka-peers=  Comma-delimited list of Kafka brokers
      --kafka-cert=   Certificate file for client authentication
      --kafka-key=    Key file for client client authentication
      --kafka-ca=     CA file for TLS client authentication
      --kafka-verify  Verify SSL certificate

Help Options:
  -h, --help          Show this help message
```

# Setup Go on Linux

```bash
sudo yum install -y golang
mkdir ~/go
export GOPATH=~/go
export PATH=$GOPATH/bin:$PATH
go get github.com/constabulary/gb/...
```

## Build

```bash
gb build
sudo bin/peekaboo \
  --static-dir src/github.com/mickep76/peekaboo/static \
  --template-dir src/github.com/mickep76/peekaboo/templates
```

## Build RPM

Make sure you have Docker configured.

```bash
make rpm
```

## Change configuration

```bash
systemctl stop peekaboo
vi /etc/sysconfig/peekaboo
```

Add "--bind-addr" bind address, defaults to "0.0.0.0". For port add "--port", defaults to 5050.

**Example:**

```
OPTIONS="--bind-addr 192.168.1.153 --port 8080"
```

Reload SystemD and then restart Peekaboo.

```bash
systemctl start peekaboo
```

# Setup Go on Mac OS X

```bash
brew install go
mkdir ~/go
export GOPATH=~/go
export PATH=$GOPATH/bin:$PATH
go get github.com/constabulary/gb/...
```

## Build

```bash
gb build
sudo bin/peekaboo \
  --static-dir src/github.com/mickep76/peekaboo/static \
  --template-dir src/github.com/mickep76/peekaboo/templates
```

## Install using Brew

```bash
brew tap mickep76/funk-gnarge
brew install peekaboo
ln -sfv /usr/local/opt/peekaboo/*.plist ~/Library/LaunchAgents
launchctl load ~/Library/LaunchAgents/homebrew.mxcl.peekaboo.plist
```

## Change configuration

```bash
launchctl unload ~/Library/LaunchAgents/homebrew.mxcl.peekaboo.plist
vi ~/Library/LaunchAgents/homebrew.mxcl.peekaboo.plist
```

Add "--bind-addr" to change bind address, defaults to "0.0.0.0". To change port add "--port", defaults to 5050.

**Example:**

```
...
    <string>/usr/local/Cellar/peekaboo/0.2.1/bin/peekaboo</string>
    <string>--static-dir</string>
    <string>/usr/local/Cellar/peekaboo/0.2.1/peekaboo/static</string>
    <string>--template-dir</string>
    <string>/usr/local/Cellar/peekaboo/0.2.1/peekaboo/templates</string>
    <string>--bind-addr</string>
    <string>192.168.1.153</string>
    <string>--port</string>
    <string>8080</string>
..
```

Restart Peekaboo.

```bash
launchctl load ~/Library/LaunchAgents/homebrew.mxcl.peekaboo.plist
```
