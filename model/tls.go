package model

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

func TLSConfig(caFile, certFile, keyFile, serverName string, insecure bool) (*tls.Config, error) {
	c := &tls.Config{
		InsecureSkipVerify: insecure,
		ServerName:         serverName,
	}

	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}

		c.Certificates = []tls.Certificate{cert}
	}

	if caFile != "" {
		c.RootCAs = x509.NewCertPool()

		b, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		c.RootCAs.AppendCertsFromPEM(b)
	}

	return c, nil
}
