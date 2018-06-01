#!/bin/bash

[ ! -x ../example/client/client ] && (cd ../example/client && go build)

../example/client/client -ca certs/ca.pem -cert certs/example_client.pem -key certs/example_client.key
