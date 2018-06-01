#!/bin/bash

[ ! -x ../example/server/server ] && (cd ../example/server && go build)

../example/server/server -ca certs/ca.pem -cert certs/example_server.pem -key certs/example_server.key
