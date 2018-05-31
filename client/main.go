package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mickep76/kvstore"
	_ "github.com/mickep76/kvstore/etcdv3"
	"github.com/mickep76/qry"

	"github.com/imc-trading/dock2box/model"
)

var taskHandler = kvstore.WatchHandler(func(kv kvstore.KeyValue) {
	log.Printf("task event: %s key: %s", kv.Event().Type, kv.Key())

	t := &model.Task{}
	if err := kv.Decode(t); err != nil {
		log.Print(err)
		return
	}

	log.Printf("task value: created: %s updated: %s uuid: %s taskDef: %s", t.Created, t.Updated, t.UUID, t.TaskDef.Name)

	if kv.PrevValue() != nil {
		t := &model.Task{}
		if err := kv.PrevDecode(t); err != nil {
			log.Print(err)
			return
		}

		log.Printf("task prev. value: created: %s updated: %s uuid: %s taskDef: %s", t.Created, t.Updated, t.UUID, t.TaskDef.Name)
	}
})

func main() {
	// Parse arguments.
	backend := flag.String("backend", "etcdv3", "Key/value store backend.")
	prefix := flag.String("prefix", "/dock2box", "Key/value store prefix.")
	endpoints := flag.String("endpoints", "127.0.0.1:2379", "Comma-delimited list of hosts in the key/value store cluster.")
	timeout := flag.Int("timeout", 5, "Connection timeout for key/value cluster in seconds.")
	keepalive := flag.Int("keepalive", 5, "Connection keepalive for key/value cluster in seconds.")
	//	ca := flag.String("ca", "", "Key/value store TLS CA certificate.")
	//	cert := flag.String("cert", "", "Key/value store TLS certificate.")
	//	key := flag.String("key", "", "Key/value store TLS key.")
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
	clients, err := ds.QueryClients(qry.New().Eq("Name", hostname))
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

	// Create lease keepalive.
	log.Printf("create lease keepalive")
	ch, err := ds.Lease().KeepAlive()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			l := <-ch
			log.Print("send keepalive for lease")
			if l.Error != nil {
				log.Print(l.Error)
			}
		}
	}()

	// Create task watch.
	log.Printf("create tasks watch")
	if err := ds.Watch(fmt.Sprintf("tasks/%s", c.HostUUID)).AddHandler(taskHandler).Start(); err != nil {
		log.Fatal(err)
	}
}
