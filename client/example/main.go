package main

import (
	"fmt"
	"log"

	"github.com/mickep76/dock2box/client"
)

func main() {
	c := client.New("http://localhost:8080/v1")

	hosts, err := c.Host.All()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(hosts.JSON())
}
