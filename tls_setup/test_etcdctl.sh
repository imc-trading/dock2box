#!/bin/bash

ETCDCTL="etcdctl --cacert certs/ca.pem --cert certs/example_server.pem --key certs/example_server.key --endpoints https://localhost:2379"

export ETCDCTL_API=3

echo "put"
${ETCDCTL} put /example/test "my value"
echo "get"
${ETCDCTL} get --prefix /example

ETCDCTL="etcdctl --cacert certs/ca.pem --cert certs/example_client.pem --key certs/example_client.key --endpoints https://localhost:2379"

echo "put"
${ETCDCTL} put /example/test "my value"
echo "get"
${ETCDCTL} get --prefix /example
