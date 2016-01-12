package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/Shopify/sarama"
)

type EventType int

// Constants for events.
const (
	STARTED EventType = iota
	CHANGED
	STOPPED
	UNKNOWN
)

type Event struct {
	Name      string      `json:"name"`
	EventType EventType   `json:"event_type"`
	Created   string      `json:"created"`
	CreatedBy CreatedBy   `json:created_by"`
	Descr     string      `json:"descr"`
	Data      interface{} `json:"data"`

	encoded []byte
	err     error
}

type CreatedBy struct {
	User    string `json:"user"`
	Service string `json:"service"`
	Host    string `json:"host"`
}

func (e *Event) ensureEncoded() {
	if e.encoded == nil && e.err == nil {
		e.encoded, e.err = json.Marshal(e)
	}
}

func (e *Event) Length() int {
	e.ensureEncoded()
	return len(e.encoded)
}

func (e *Event) Encode() ([]byte, error) {
	e.ensureEncoded()
	return e.encoded, e.err
}

func createTLSConfig(certFile *string, keyFile *string, caFile *string, verify bool) (t *tls.Config) {
	if certFile == nil || keyFile == nil || caFile == nil {
		return t
	}

	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	caCert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	t = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: verify,
	}

	return t
}

func newProducer(brokerList []string, cert *string, key *string, ca *string, verify bool) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	tlsConfig := createTLSConfig(cert, key, ca, verify)
	if tlsConfig != nil {
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Enable = true
	}

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %s", err.Error())
	}

	return producer
}
