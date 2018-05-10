package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mickep76/kvstore"
	_ "github.com/mickep76/kvstore/etcdv3"
	"github.com/mickep76/qry"

	"github.com/mickep76/kvstore/example/model"
)

func main() {
	// Parse arguments.
	backend := flag.String("backend", "etcdv3", "Key/value store backend.")
	prefix := flag.String("prefix", "/dock2box", "Key/value store prefix.")
	endpoints := flag.String("endpoints", "127.0.0.1:2379", "Comma-delimited list of hosts in the key/value store cluster.")
	timeout := flag.Int("timeout", 5, "Connection timeout for key/value cluster in seconds.")
	keepalive := flag.Int("keepalive", 5, "Connection keepalive for key/value cluster in seconds.")
	flag.Parse()

	// Connect to etcd.
	log.Printf("connect to etcd")
	ds, err := model.NewDatastore(*backend, strings.Split(*endpoints, ","), *keepalive, kvstore.WithTimeout(*timeout), kvstore.WithEncoding("json"), kvstore.WithPrefix(*prefix))
	if err != nil {
		log.Fatal(err)
	}

	// Find existing client in datastore.
	log.Printf("find existing client in datastore")
	hostname, _ := os.Hostname()
	clients, err := ds.QueryClients(qry.Eq("Hostname", hostname))
	if err != nil {
		log.Fatal(err)
	}

	var c *model.Client
	if len(clients) > 0 {
		// Update client in datastore.
		log.Printf("update client in datastore")
		c = clients[0]
		if err := ds.UpdateClient(c); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("create new client")
		c = model.NewClient(hostname)

		// Create client in datastore.
		log.Printf("create client in datastore")
		if err := ds.CreateClient(c); err != nil {
			log.Fatal(err)
		}
	}

	// Update client in etcd after 10 seconds.
	timer := time.NewTimer(10 * time.Second)
	go func() {
		<-timer.C

		log.Printf("update client in etcd")
		if err := ds.UpdateClient(c); err != nil {
			log.Fatal(err)
		}
	}()

	// Create lease keepalive.
	log.Printf("create lease keepalive")
	ch, err := ds.Lease().KeepAlive()
	if err != nil {
		log.Fatal(err)
	}

	for {
		l := <-ch
		log.Print("send keepalive for lease")
		if l.Error != nil {
			log.Print(l.Error)
		}
	}
}
