.PHONY: cfssl ca req clean
  
all: ca req

clean:
	rm -rf certs data

preq:
	go get -u -tags nopkcs11 github.com/cloudflare/cfssl/cmd/cfssl
	go get -u github.com/cloudflare/cfssl/cmd/cfssljson

ca:
	mkdir -p certs
	cfssl gencert -initca config/ca_csr.json | cfssljson -bare certs/ca
	mv certs/ca-key.pem certs/ca.key

req:
	cfssl gencert \
	  -ca certs/ca.pem \
	  -ca-key certs/ca.key \
	  -config config/ca_cfg.json config/etcd_server.json | cfssljson -bare certs/etcd
	mv certs/etcd-key.pem certs/etcd.key

	cfssl gencert \
	  -ca certs/ca.pem \
	  -ca-key certs/ca.key \
	  -config config/ca_cfg.json config/example_server.json | cfssljson -bare certs/example_server
	mv certs/example_server-key.pem certs/example_server.key

	cfssl gencert \
	  -ca certs/ca.pem \
	  -ca-key certs/ca.key \
	  -config config/ca_cfg.json config/example_client.json | cfssljson -bare certs/example_client
	mv certs/example_client-key.pem certs/example_client.key

.PHONY: clean preq ca req
