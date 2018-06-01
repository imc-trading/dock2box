#!/bin/bash

[ ! -x ../server/server ] && (cd ../server && go build)

../server/server -ca certs/ca.pem -cert certs/d2b_server.pem -key certs/d2b_server.key
