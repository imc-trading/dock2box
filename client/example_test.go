package client_test

import (
	"fmt"
	"log"

	"github.com/imc-trading/dock2box/client"
)

func Example_GetHost() {
	server := "http://localhost:8080/v1"
	hostname := "test1.example.com"

	clnt := client.New(server)
	h, err := clnt.Host.Get(hostname)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%v\n", string(h.JSON()))
}
