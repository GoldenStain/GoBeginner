package main

import (
	"golang.org/x/text/language"

	"fmt"
	"os"
)

func main() {
	for _, arg := range os.Args[1:] {
		tag, err := language.Parse(arg)
		if err != nil {
			fmt.Printf("%s: error: %v\n", arg, err)
		} else if tag == language.Und {
			fmt.Printf("%s: Undefined\n", arg)
		} else {
			fmt.Printf("%s: %s\n", arg, tag)
		}
	}
}
