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

	"github.com/imc-trading/dock2box/handler"
	"github.com/imc-trading/dock2box/model"
)

var clientHandler = kvstore.WatchHandler(func(kv kvstore.KeyValue) {
	log.Printf("client event: %s key: %s", kv.Event().Type, kv.Key())

	c := &model.Client{}
	if err := kv.Decode(c); err != nil {
		log.Print(err)
		return
	}

	log.Printf("client value: created: %s updated: %s uuid: %s hostname: %s", c.Created, c.Updated, c.UUID, c.Name)

	if kv.PrevValue() != nil {
		c := &model.Client{}
		if err := kv.PrevDecode(c); err != nil {
			log.Print(err)
			return
		}

		log.Printf("client prev. value: created: %s updated: %s uuid: %s hostname: %s", c.Created, c.Updated, c.UUID, c.Name)
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
	servers, err := ds.QueryServers(qry.New().Eq("Name", hostname))
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
		s = model.NewServer(hostname)

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

	/* client
	 * host
	 * image
	 * pool
	 * rack
	 * role
	 * server
	 * site
	 * subnet
	 * tenant
	 */

	// Client handlers.
	log.Printf("add route /api/clients")
	router.HandleFunc("/api/clients", h.AllClients).Methods("GET")
	router.HandleFunc("/api/clients/{uuid}", h.OneClient).Methods("GET")

	// Host handlers.
	log.Printf("add route /api/hosts")
	router.HandleFunc("/api/hosts", h.AllHosts).Methods("GET")
	router.HandleFunc("/api/hosts", h.CreateHost).Methods("POST")
	router.HandleFunc("/api/hosts/{uuid}", h.OneHost).Methods("GET")
	router.HandleFunc("/api/hosts/{uuid}", h.UpdateHost).Methods("PUT")
	router.HandleFunc("/api/hosts/{uuid}", h.DeleteHost).Methods("DELETE")

	// Image handlers.
	log.Printf("add route /api/images")
	router.HandleFunc("/api/images", h.AllImages).Methods("GET")
	router.HandleFunc("/api/images", h.CreateImage).Methods("POST")
	router.HandleFunc("/api/images/{uuid}", h.OneImage).Methods("GET")
	router.HandleFunc("/api/images/{uuid}", h.UpdateImage).Methods("PUT")
	router.HandleFunc("/api/images/{uuid}", h.DeleteImage).Methods("DELETE")

	// Pool handlers.
	log.Printf("add route /api/pools")
	router.HandleFunc("/api/pools", h.AllPools).Methods("GET")
	router.HandleFunc("/api/pools", h.CreatePool).Methods("POST")
	router.HandleFunc("/api/pools/{uuid}", h.OnePool).Methods("GET")
	router.HandleFunc("/api/pools/{uuid}", h.UpdatePool).Methods("PUT")
	router.HandleFunc("/api/pools/{uuid}", h.DeletePool).Methods("DELETE")

	// Rack handlers.
	log.Printf("add route /api/racks")
	router.HandleFunc("/api/racks", h.AllRacks).Methods("GET")
	router.HandleFunc("/api/racks", h.CreateRack).Methods("POST")
	router.HandleFunc("/api/racks/{uuid}", h.OneRack).Methods("GET")
	router.HandleFunc("/api/racks/{uuid}", h.UpdateRack).Methods("PUT")
	router.HandleFunc("/api/racks/{uuid}", h.DeleteRack).Methods("DELETE")

	// Role handlers.
	log.Printf("add route /api/roles")
	router.HandleFunc("/api/roles", h.AllRoles).Methods("GET")
	router.HandleFunc("/api/roles", h.CreateRole).Methods("POST")
	router.HandleFunc("/api/roles/{uuid}", h.OneRole).Methods("GET")
	router.HandleFunc("/api/roles/{uuid}", h.UpdateRole).Methods("PUT")
	router.HandleFunc("/api/roles/{uuid}", h.DeleteRole).Methods("DELETE")

	// Server handlers.
	log.Printf("add route /api/servers")
	router.HandleFunc("/api/servers", h.AllServers).Methods("GET")
	router.HandleFunc("/api/servers/{uuid}", h.OneServer).Methods("GET")

	// Site handlers.
	log.Printf("add route /api/sites")
	router.HandleFunc("/api/sites", h.AllSites).Methods("GET")
	router.HandleFunc("/api/sites", h.CreateSite).Methods("POST")
	router.HandleFunc("/api/sites/{uuid}", h.OneSite).Methods("GET")
	router.HandleFunc("/api/sites/{uuid}", h.UpdateSite).Methods("PUT")
	router.HandleFunc("/api/sites/{uuid}", h.DeleteSite).Methods("DELETE")

	// Subnet handlers.
	log.Printf("add route /api/subnets")
	router.HandleFunc("/api/subnets", h.AllSubnets).Methods("GET")
	router.HandleFunc("/api/subnets", h.CreateSubnet).Methods("POST")
	router.HandleFunc("/api/subnets/{uuid}", h.OneSubnet).Methods("GET")
	router.HandleFunc("/api/subnets/{uuid}", h.UpdateSubnet).Methods("PUT")
	router.HandleFunc("/api/subnets/{uuid}", h.DeleteSubnet).Methods("DELETE")

	// Tenant handlers.
	log.Printf("add route /api/tenants")
	router.HandleFunc("/api/tenants", h.AllTenants).Methods("GET")
	router.HandleFunc("/api/tenants", h.CreateTenant).Methods("POST")
	router.HandleFunc("/api/tenants/{uuid}", h.OneTenant).Methods("GET")
	router.HandleFunc("/api/tenants/{uuid}", h.UpdateTenant).Methods("PUT")
	router.HandleFunc("/api/tenants/{uuid}", h.DeleteTenant).Methods("DELETE")

	// Start https listener.
	log.Printf("start http listener")
	logr := handlers.LoggingHandler(os.Stdout, router)
	if err := http.ListenAndServe(*bind, logr); err != nil {
		log.Fatal("http listener:", err)
	}
}
