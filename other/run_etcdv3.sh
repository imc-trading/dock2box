#!/bin/bash

docker run --rm -p 2379:2379 -p 2380:2380 --name etcd -v $PWD/certs:/etc/ssl/certs gcr.io/etcd-development/etcd:v3.2 /usr/local/bin/etcd \
	--name my-etcd-1 \
	--data-dir /etcd-data \
	--listen-client-urls https://0.0.0.0:2379 \
	--advertise-client-urls https://0.0.0.0:2379 \
	--listen-peer-urls https://0.0.0.0:2380 \
	--initial-advertise-peer-urls https://0.0.0.0:2380 \
	--initial-cluster my-etcd-1=https://0.0.0.0:2380 \
	--initial-cluster-token my-etcd-token \
	--initial-cluster-state new \
	--cert-file /etc/ssl/certs/server.crt \
	--key-file /etc/ssl/certs/server.key \
	--peer-cert-file /etc/ssl/certs/server.crt \
	--peer-key-file /etc/ssl/certs/server.key
