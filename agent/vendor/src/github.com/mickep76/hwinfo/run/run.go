package main

import (
	"encoding/json"
	"fmt"

	"github.com/mickep76/hwinfo"
)

func main() {
	d := hwinfo.NewHWInfo()
	d.GetTTL()

	//	fmt.Println(d)

	b, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(b))
}
