#!/bin/bash

etcd --name etcd1 \
	--data-dir data \
	--trusted-ca-file certs/ca.pem \
	--cert-file certs/etcd.pem \
	--key-file certs/etcd.key \
	--client-cert-auth \
	--advertise-client-urls https://127.0.0.1:2379 \
	--listen-client-urls https://127.0.0.1:2379
