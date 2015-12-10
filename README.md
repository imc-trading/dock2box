# Build

```bash
export GOPATH=~/go
export PATH=$PATH:$HOME/go/bin

mkdir -p $GOPATH/src
go get github.com/imc-trading/dock2box/d2bsrv
go get github.com/imc-trading/dock2box/d2bcli
```

# Test server

First you need to install and run MongoDB.

**Install MongoDB on Mac OS X:**

```bash
brew install mongodb
ln -sfv /usr/local/opt/mongodb/*.plist ~/Library/LaunchAgents
launchctl load ~/Library/LaunchAgents/homebrew.mxcl.mongodb.plist
```

**Start server:**

```bash
cd $GOPATH/github.com/imc-trading/dock2box/d2bsrv
./d2bsrv
```

**Test server:**

```bash
cd $GOPATH/github.com/imc-trading/dock2box/d2bsrv
./test.sh
```

# Test CLI

**Get host:**

```bash
d2bcli get host test1.example.com
```

**Create host:**

```bash
d2bcli create host test2.example.com -p
```
