#!/bin/bash

# Generate TLS certificate
SSL_CN="auth"
SSL_O="Auth"
SSL_C="US"

mkdir -p certs
openssl req -nodes -new -x509 -keyout certs/server.key -out certs/server.crt -subj "/CN=${SSL_CN}/O=${SSL_O}/C=${SSL_C}"
