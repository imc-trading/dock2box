package prompt

import (
	"fmt"
	"strconv"
	"strings"
)

func Choice(msg string, list []string) int {
	fmt.Printf("\n")
	for i, v := range list {
		fmt.Printf("%d) %s\n", i, v)
	}
	fmt.Printf("\n")

	for {
		fmt.Printf(msg + ": ")

		var inp string
		fmt.Scanln(&inp)

		i, err := strconv.Atoi(inp)
		if err != nil {
			fmt.Println("Input needs to be a number")
		} else if i >= 0 && i < len(list) {
			return i
		}
	}
}

func Bool(msg string, def bool) bool {
	for {
		fmt.Printf("%s: [true/false]? (%v) ", msg, def)

		var inp string
		fmt.Scanln(&inp)

		switch {
		case inp == "":
			return def
		case strings.ToLower(inp) == "true":
			return true
		case strings.ToLower(inp) == "false":
			return false
		}
	}
}

func String(msg string) string {
	fmt.Printf("%s: ", msg)

	var inp string
	fmt.Scanln(&inp)
	return inp
}
