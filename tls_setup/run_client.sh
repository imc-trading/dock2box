#!/bin/bash

[ ! -x ../client/client ] && (cd ../client && go build)

../client/client -ca certs/ca.pem -cert certs/d2b_client.pem -key certs/d2b_client.key
