# Setup TLS encryption

## Create certs

```bash
make clean
make preq
make
```

## Run etcd

```bash
./run_etcd.sh
```

## Setup auth.

```bash
./setup_auth.sh
```

## Test auth.

```bash
./test_etcdctl.sh
```

## Run example server/client

```bash
./run_server.sh
```

And in a separate terminal.

```bash
./run_client.sh
```
