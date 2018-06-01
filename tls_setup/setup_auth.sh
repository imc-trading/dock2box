#!/bin/bash

ETCDCTL="etcdctl --cacert certs/ca.pem --cert certs/etcd.pem --key certs/etcd.key --endpoints https://localhost:2379"

export ETCDCTL_API=3

# Create root user
${ETCDCTL} user add root:abc123

# Create roles
${ETCDCTL} role add example-server-rw
${ETCDCTL} role add example-client-rw

# Add permissions to roles
${ETCDCTL} role grant-permission example-server-rw --prefix=true readwrite /example
${ETCDCTL} role grant-permission example-client-rw --prefix=true readwrite /example/clients

# Create users
${ETCDCTL} user add example-server:abc123
${ETCDCTL} user add example-client:abc123

# Grant role to user
${ETCDCTL} user grant-role example-server example-server-rw
${ETCDCTL} user grant-role example-client example-client-rw

# Enable authentication
${ETCDCTL} auth enable
