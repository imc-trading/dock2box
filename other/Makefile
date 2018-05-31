.PHONY: cfssl ca req clean
  
all: ca req

clean:
	rm -rf certs

preq:
	go get -u -tags nopkcs11 github.com/cloudflare/cfssl/cmd/cfssl
	go get -u github.com/cloudflare/cfssl/cmd/cfssljson

ca:
	mkdir -p certs
	cfssl gencert -initca ca_csr.json | cfssljson -bare certs/ca

req:
	cfssl gencert \
	  -ca certs/ca.pem \
	  -ca-key certs/ca-key.pem \
	  -config ca_cfg.json req_csr.json | cfssljson -bare certs/etcd1

	cfssl gencert \
	  -ca certs/ca.pem \
	  -ca-key certs/ca-key.pem \
	  -config ca_cfg.json req_csr.json | cfssljson -bare certs/etcd2

	cfssl gencert \
	  -ca certs/ca.pem \
	  -ca-key certs/ca-key.pem \
	  -config ca_cfg.json req_csr.json | cfssljson -bare certs/etcd3

.PHONY: clean preq ca req
