package main

import (
	"fmt"
	"log"
	"os"

	"library/pkg/library"
)

func main() {
	if len(os.Args[1:]) == 0 {
		fmt.Printf("Usage: %s <json data file>", os.Args[0])
		os.Exit(1)
	}

	// construct the library from the data file
	l, err := library.Build(library.DMC{}, os.Args[1])

	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
}
