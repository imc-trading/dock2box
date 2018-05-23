package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mickep76/auth"
	_ "github.com/mickep76/auth/ldap"
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

/* TODO:
 * - User/pass and TLS for etcd
 * - tasks prefixed by host like: /tasks/<host uuid>/<task uuid> to allow client to have a watcher
 */

func main() {
	// Parse arguments.
	kvsBackend := flag.String("kvs-backend", "etcdv3", "Key/value store backend.")
	kvsPrefix := flag.String("kvs-prefix", "/dock2box", "Key/value store prefix.")
	kvsEndpoints := flag.String("kvs-endpoints", "127.0.0.1:2379", "Comma-delimited list of hosts in the key/value store cluster.")
	kvsTimeout := flag.Int("kvs-timeout", 5, "Connection timeout for key/value cluster in seconds.")
	kvsKeepalive := flag.Int("kvs-keepalive", 5, "Connection keepalive for key/value cluster in seconds.")
	kvsUser := flag.String("kvs-user", "", "Key/avlue store user.")
	kvsPassword := flag.String("kvs-password", "", "Key/value store password.")
	kvsInsecure := flag.Bool("kvs-insecure", false, "Insecure TLS.")

	httpBind := flag.String("http-bind", "127.0.0.1:8080", "Bind to address and port.")
	httpCert := flag.String("http-cert", "server.crt", "TLS HTTPS cert.")
	httpKey := flag.String("http-key", "server.key", "TLS HTTPS key.")

	authBackend := flag.String("auth-backend", "ad", "Auth. backend either ad or ldap.")
	authEndpoint := flag.String("auth-endpoint", "ldap:389", "LDAP server and port.")
	authInsecure := flag.Bool("auth-insecure", false, "Insecure TLS.")
	authDomain := flag.String("auth-domain", "", "AD Domain.")
	authBase := flag.String("auth-base", "", "LDAP Base.")

	jwtPrivKey := flag.String("jwt-priv-key", "private.rsa", "Private RSA key.")
	jwtPubKey := flag.String("jwt-pub-key", "public.rsa", "Public RSA key.")

	flag.Parse()

	// Create auth. TLS config.
	authTLS := &tls.Config{
		InsecureSkipVerify: *authInsecure,
		ServerName:         strings.Split(*authEndpoint, ":")[0], // Send SNI (Server Name Indication) for host that serves multiple aliases.
	}

	// Create new auth. connection.
	c, err := auth.Open(*authBackend, []string{*authEndpoint}, auth.TLS(authTLS), auth.Domain(*authDomain), auth.Base(*authBase))
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Create JWT.
	j := auth.NewJWT(auth.SignRS512, time.Duration(24)*time.Hour, time.Duration(5)*time.Minute)

	// Load RSA private key.
	if j.LoadPrivateKey(*jwtPrivKey); err != nil {
		log.Fatal(err)
	}

	// Load RSA public key.
	if err := j.LoadPublicKey(*jwtPubKey); err != nil {
		log.Fatal(err)
	}

	// Create etcd TLS config.
	etcdTLS := &tls.Config{
		InsecureSkipVerify: *kvsInsecure,
	}

	// Connect to etcd.
	log.Printf("connect to etcd")
	ds, err := model.NewDatastore(*kvsBackend, strings.Split(*kvsEndpoints, ","), *kvsKeepalive, kvstore.WithTimeout(*kvsTimeout), kvstore.WithEncoding("json"), kvstore.WithPrefix(*kvsPrefix), kvstore.WithUser(*kvsUser), kvstore.WithPassword(*kvsPassword), kvstore.WithTLS(etcdTLS))
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
	h := handler.NewHandler(ds, c, j)

	// Auth. handlers.
	log.Printf("add route /login, /renew, /verify")
	router.HandleFunc("/login", h.Login).Methods("POST")
	router.HandleFunc("/renew", h.Renew).Methods("GET")
	router.HandleFunc("/verify", h.Verify).Methods("GET")

	// Client handlers.
	log.Printf("add route /api/clients")
	router.Handle("/api/clients", j.Authorized(http.HandlerFunc(h.AllClients))).Methods("GET")
	router.Handle("/api/clients/{uuid}", j.Authorized(http.HandlerFunc(h.OneClient))).Methods("GET")

	// Host handlers.
	log.Printf("add route /api/hosts")
	router.Handle("/api/hosts", j.Authorized(http.HandlerFunc(h.AllHosts))).Methods("GET")
	router.Handle("/api/hosts", j.Authorized(http.HandlerFunc(h.CreateHost))).Methods("POST")
	router.Handle("/api/hosts/{uuid}", j.Authorized(http.HandlerFunc(h.OneHost))).Methods("GET")
	router.Handle("/api/hosts/{uuid}", j.Authorized(http.HandlerFunc(h.UpdateHost))).Methods("PUT")
	router.Handle("/api/hosts/{uuid}", j.Authorized(http.HandlerFunc(h.DeleteHost))).Methods("DELETE")

	// Image handlers.
	log.Printf("add route /api/images")
	router.Handle("/api/images", j.Authorized(http.HandlerFunc(h.AllImages))).Methods("GET")
	router.Handle("/api/images", j.Authorized(http.HandlerFunc(h.CreateImage))).Methods("POST")
	router.Handle("/api/images/{uuid}", j.Authorized(http.HandlerFunc(h.OneImage))).Methods("GET")
	router.Handle("/api/images/{uuid}", j.Authorized(http.HandlerFunc(h.UpdateImage))).Methods("PUT")
	router.Handle("/api/images/{uuid}", j.Authorized(http.HandlerFunc(h.DeleteImage))).Methods("DELETE")

	// Pool handlers.
	log.Printf("add route /api/pools")
	router.Handle("/api/pools", j.Authorized(http.HandlerFunc(h.AllPools))).Methods("GET")
	router.Handle("/api/pools", j.Authorized(http.HandlerFunc(h.CreatePool))).Methods("POST")
	router.Handle("/api/pools/{uuid}", j.Authorized(http.HandlerFunc(h.OnePool))).Methods("GET")
	router.Handle("/api/pools/{uuid}", j.Authorized(http.HandlerFunc(h.UpdatePool))).Methods("PUT")
	router.Handle("/api/pools/{uuid}", j.Authorized(http.HandlerFunc(h.DeletePool))).Methods("DELETE")

	// Rack handlers.
	log.Printf("add route /api/racks")
	router.Handle("/api/racks", j.Authorized(http.HandlerFunc(h.AllRacks))).Methods("GET")
	router.Handle("/api/racks", j.Authorized(http.HandlerFunc(h.CreateRack))).Methods("POST")
	router.Handle("/api/racks/{uuid}", j.Authorized(http.HandlerFunc(h.OneRack))).Methods("GET")
	router.Handle("/api/racks/{uuid}", j.Authorized(http.HandlerFunc(h.UpdateRack))).Methods("PUT")
	router.Handle("/api/racks/{uuid}", j.Authorized(http.HandlerFunc(h.DeleteRack))).Methods("DELETE")

	// Role handlers.
	log.Printf("add route /api/roles")
	router.Handle("/api/roles", j.Authorized(http.HandlerFunc(h.AllRoles))).Methods("GET")
	router.Handle("/api/roles", j.Authorized(http.HandlerFunc(h.CreateRole))).Methods("POST")
	router.Handle("/api/roles/{uuid}", j.Authorized(http.HandlerFunc(h.OneRole))).Methods("GET")
	router.Handle("/api/roles/{uuid}", j.Authorized(http.HandlerFunc(h.UpdateRole))).Methods("PUT")
	router.Handle("/api/roles/{uuid}", j.Authorized(http.HandlerFunc(h.DeleteRole))).Methods("DELETE")

	// Server handlers.
	log.Printf("add route /api/servers")
	router.Handle("/api/servers", j.Authorized(http.HandlerFunc(h.AllServers))).Methods("GET")
	router.Handle("/api/servers/{uuid}", j.Authorized(http.HandlerFunc(h.OneServer))).Methods("GET")

	// Site handlers.
	log.Printf("add route /api/sites")
	router.Handle("/api/sites", j.Authorized(http.HandlerFunc(h.AllSites))).Methods("GET")
	router.Handle("/api/sites", j.Authorized(http.HandlerFunc(h.CreateSite))).Methods("POST")
	router.Handle("/api/sites/{uuid}", j.Authorized(http.HandlerFunc(h.OneSite))).Methods("GET")
	router.Handle("/api/sites/{uuid}", j.Authorized(http.HandlerFunc(h.UpdateSite))).Methods("PUT")
	router.Handle("/api/sites/{uuid}", j.Authorized(http.HandlerFunc(h.DeleteSite))).Methods("DELETE")

	// Subnet handlers.
	log.Printf("add route /api/subnets")
	router.Handle("/api/subnets", j.Authorized(http.HandlerFunc(h.AllSubnets))).Methods("GET")
	router.Handle("/api/subnets", j.Authorized(http.HandlerFunc(h.CreateSubnet))).Methods("POST")
	router.Handle("/api/subnets/{uuid}", j.Authorized(http.HandlerFunc(h.OneSubnet))).Methods("GET")
	router.Handle("/api/subnets/{uuid}", j.Authorized(http.HandlerFunc(h.UpdateSubnet))).Methods("PUT")
	router.Handle("/api/subnets/{uuid}", j.Authorized(http.HandlerFunc(h.DeleteSubnet))).Methods("DELETE")

	// Tenant handlers.
	log.Printf("add route /api/tenants")
	router.Handle("/api/tenants", j.Authorized(http.HandlerFunc(h.AllTenants))).Methods("GET")
	router.Handle("/api/tenants", j.Authorized(http.HandlerFunc(h.CreateTenant))).Methods("POST")
	router.Handle("/api/tenants/{uuid}", j.Authorized(http.HandlerFunc(h.OneTenant))).Methods("GET")
	router.Handle("/api/tenants/{uuid}", j.Authorized(http.HandlerFunc(h.UpdateTenant))).Methods("PUT")
	router.Handle("/api/tenants/{uuid}", j.Authorized(http.HandlerFunc(h.DeleteTenant))).Methods("DELETE")

	// Start https listener.
	log.Printf("start http listener")
	logr := handlers.LoggingHandler(os.Stdout, router)
	if err := http.ListenAndServeTLS(*httpBind, *httpCert, *httpKey, logr); err != nil {
		log.Fatal("http listener:", err)
	}
}
