package prompt

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
