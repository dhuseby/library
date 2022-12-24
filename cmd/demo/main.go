package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"library/pkg/library"
)

func help() {
	fmt.Printf("Usage: %s <command> <json data file>\n", os.Args[0])
	fmt.Printf("Commands:\n")
	fmt.Printf("\tlist -- lists all books\n")
	fmt.Printf("\tcreate -- creates a new book\n")
	fmt.Printf("\tread -- looks up a book by ISBN\n")
	fmt.Printf("\tupdate -- updates a book by ISBN\n")
	fmt.Printf("\tdelete -- deletes a book by ISBN\n")
}

func main() {
	if len(os.Args[1:]) != 2 {
		help()
		os.Exit(1)
	}

	// get the command and filename
	cmd := os.Args[1]
	f := os.Args[2]

	// construct the library from the data file
	l, err := library.Build(library.DMC{}, f)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	switch cmd {
	case "list":
		// list the books
		books, err := l.ListAll()
		if err != nil {
			log.Fatal(err)
		}
		for _, b := range books {
			fmt.Printf("%s\n", b)
		}
	case "create":
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the ISBN: ")
		id, _ := reader.ReadString('\n')
		id = strings.TrimSpace(id)
		_, err := l.Read(id)
		if err == nil {
			log.Fatal("A book with ISBN:", id, " already exists")
		}
		fmt.Print("Enter the Title: ")
		title, _ := reader.ReadString('\n')
		title = strings.TrimSpace(title)
		book := library.Book{ id, title }
		fmt.Printf("Creating book: %s\n", book)
		err = l.Create(book)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Created")
	case "read":
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the ISBN: ")
		id, _ := reader.ReadString('\n')
		id = strings.TrimSpace(id)
		book, err := l.Read(id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", book)
	case "update":
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the ISBN: ")
		id, _ := reader.ReadString('\n')
		id = strings.TrimSpace(id)
		book, err := l.Read(id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", book)
		fmt.Print("Enter new Title: ")
		title, _ := reader.ReadString('\n')
		title = strings.TrimSpace(title)
		book.Title = title
		err = l.Update(book)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Updated")
	case "delete":
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the ISBN: ")
		id, _ := reader.ReadString('\n')
		id = strings.TrimSpace(id)
		book, err := l.Delete(id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Deleted: %s\n", book)
	default:
		help()
		log.Fatal("invalid command: ", cmd)
	}
}
