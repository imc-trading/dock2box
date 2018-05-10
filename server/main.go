package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mickep76/kvstore"
	_ "github.com/mickep76/kvstore/etcdv3"
	"github.com/mickep76/qry"

	"github.com/mickep76/kvstore/example/handler"
	"github.com/mickep76/kvstore/example/model"
)

var clientHandler = kvstore.WatchHandler(func(kv kvstore.KeyValue) {
	log.Printf("client event: %s key: %s", kv.Event().Type, kv.Key())

	c := &model.Client{}
	if err := kv.Decode(c); err != nil {
		log.Print(err)
		return
	}

	log.Printf("client value: created: %s updated: %s uuid: %s hostname: %s", c.Created, c.Updated, c.UUID, c.Hostname)

	if kv.PrevValue() != nil {
		c := &model.Client{}
		if err := kv.PrevDecode(c); err != nil {
			log.Print(err)
			return
		}

		log.Printf("client prev. value: created: %s updated: %s uuid: %s hostname: %s", c.Created, c.Updated, c.UUID, c.Hostname)
	}
})

func main() {
	// Parse arguments.
	backend := flag.String("backend", "etcdv3", "Key/value store backend.")
	prefix := flag.String("prefix", "/dock2box", "Key/value store prefix.")
	endpoints := flag.String("endpoints", "127.0.0.1:2379", "Comma-delimited list of hosts in the key/value store cluster.")
	timeout := flag.Int("timeout", 5, "Connection timeout for key/value cluster in seconds.")
	keepalive := flag.Int("keepalive", 5, "Connection keepalive for key/value cluster in seconds.")
	bind := flag.String("bind", "127.0.0.1:8080", "Bind to address and port.")
	flag.Parse()

	// Connect to etcd.
	log.Printf("connect to etcd")
	ds, err := model.NewDatastore(*backend, strings.Split(*endpoints, ","), *keepalive, kvstore.WithTimeout(*timeout), kvstore.WithEncoding("json"), kvstore.WithPrefix(*prefix))
	if err != nil {
		log.Fatal(err)
	}

	// Find existing server in datastore.
	log.Printf("find existing server in datastore")
	hostname, _ := os.Hostname()
	servers, err := ds.QueryServers(qry.Eq("Hostname", hostname))
	if err != nil {
		log.Fatal(err)
	}

	var s *model.Server
	if len(servers) > 0 {
		// Update server in datastore.
		log.Printf("update server in datastore")
		s = servers[0]
		if err := ds.UpdateServer(s); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("create new server")
		s = model.NewServer(hostname, *bind)

		// Create server in datastore.
		log.Printf("create server in datastore")
		if err := ds.CreateServer(s); err != nil {
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

	// Create client watch.
	log.Printf("create client watch")
	go func() {
		if err := ds.Watch("clients").AddHandler(clientHandler).Start(); err != nil {
			log.Fatal(err)
		}
	}()

	// Create new router.
	log.Printf("create http router")
	router := mux.NewRouter()
	h := handler.NewHandler(ds)

	// Client handlers.
	log.Printf("add route /api/clients")
	router.HandleFunc("/api/clients", h.AllClients).Methods("GET")

	// Server handlers.
	log.Printf("add route /api/servers")
	router.HandleFunc("/api/servers", h.AllServers).Methods("GET")

	// Start https listener.
	log.Printf("start http listener")
	logr := handlers.LoggingHandler(os.Stdout, router)
	if err := http.ListenAndServe(*bind, logr); err != nil {
		log.Fatal("http listener:", err)
	}
}
