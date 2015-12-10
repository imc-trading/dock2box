# Build

```bash
export GOPATH=~/go
export PATH=$PATH:$HOME/go/bin

mkdir -p $GOPATH/src
go get github.com/imc-trading/dock2box/d2bsrv
go get github.com/imc-trading/dock2box/d2bcli
```

# Test server

**Start server:**

```bash
cd $GOPATH/github.com/imc-trading/dock2box/d2bsrv
./d2bsrv
./test.sh
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
