package main

import (
	"fmt"
	"log"

	"github.com/GoldenStain/GoBeginner/greetings"
)

func main() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	names := []string{"zzf", "zzf2", "zzf3"}
	messages, err := greetings.Hellos(names)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(messages)
}
