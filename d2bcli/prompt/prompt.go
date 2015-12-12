package prompt

// TODO
// - Validate by using a ref to JSON Schema

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Prompt struct {
	NoDefault bool
	Default   string
	FuncPtr   func(string, string) bool
	FuncInp   string
}

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
		if def == true {
			fmt.Printf("%s: [yes/no]? (yes) ", msg)
		} else {
			fmt.Printf("%s: [yes/no]? (no) ", msg)
		}

		var inp string
		fmt.Scanln(&inp)

		switch {
		case inp == "":
			return def
		case strings.ToLower(inp) == "yes":
			return true
		case strings.ToLower(inp) == "no":
			return false
		}
	}
}

func String(msg string, p Prompt) string {
	for {
		if p.NoDefault {
			fmt.Printf("%s: ", msg)
		} else {
			fmt.Printf("%s: (%s) ", msg, p.Default)
		}

		var inp string
		fmt.Scanln(&inp)

		switch {
		case !p.NoDefault && inp == "":
			return p.Default
		case p.FuncPtr(inp, p.FuncInp):
			return inp
		}
	}
}

func Regex(inp string, regex string) bool {
	rx := regexp.MustCompile(regex)
	if rx.MatchString(inp) {
		return true
	}

	fmt.Printf("Input: %s doesn't match regex: %s\n", inp, regex)
	return false
}

func Enum(inp string, list string) bool {
	for _, v := range strings.Split(list, ",") {
		fmt.Println(inp, v)
		if inp == v {
			return true
		}
	}

	fmt.Printf("Input: %s doesn't match enum: %s\n", inp, list)
	return false
}
