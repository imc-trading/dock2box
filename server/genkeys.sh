#!/bin/bash

# Generate RSA private and public key
openssl genrsa -out private.rsa 2048
openssl rsa -in private.rsa -outform PEM -pubout -out public.rsa

# Generate TLS certificate
SSL_CN="auth"
SSL_O="Auth"
SSL_C="US"

openssl req -nodes -new -x509 -keyout server.key -out server.crt -subj "/CN=${SSL_CN}/O=${SSL_O}/C=${SSL_C}"
