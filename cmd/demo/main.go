package main

import (
	"fmt"
	"library/pkg/library"
)

func main() {
	b := library.Book{ ID: "1234567890", Title: "Foo" }
	fmt.Println(b)
}
