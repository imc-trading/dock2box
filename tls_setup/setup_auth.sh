#!/bin/bash

ETCDCTL="etcdctl --cacert certs/ca.pem --cert certs/etcd.pem --key certs/etcd.key --endpoints https://localhost:2379"

export ETCDCTL_API=3

# Create root user
${ETCDCTL} user add root:abc123

# Create roles
${ETCDCTL} role add d2b-server-rw
${ETCDCTL} role add d2b-client-rw

# Add permissions to roles
${ETCDCTL} role grant-permission d2b-server-rw --prefix=true readwrite /dock2box
${ETCDCTL} role grant-permission d2b-client-rw --prefix=true readwrite /dock2box/clients

# Create users
${ETCDCTL} user add d2b-server:abc123
${ETCDCTL} user add d2b-client:abc123

# Grant role to user
${ETCDCTL} user grant-role d2b-server d2b-server-rw
${ETCDCTL} user grant-role d2b-client d2b-client-rw

# Enable authentication
${ETCDCTL} auth enable
