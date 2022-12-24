# Library Coding Assignment

![Apache 2.0][license-image]

Using Go write a minimal library API that can perform the following functions:

1. List all books in the library
2. CRUD operations on a single book

## Todo

- [x] Set up project
  - [x] Create `library` package
  - [x] Create demo binary
  - [x] Set up the Makefile
  - [x] Create this README
  - [x] Create failing tests
- [x] Fill in tests
- [x] Create concrete impl for the `Library` interface that uses an in-memory map and atomic.Value protection
  - [ ] Pass all tests
- [ ] Add Future Proofing and Maintenance discussion to the README
  - [ ] Add information on abstracting the searching mechanism behind an interface to allow for looking up `Book`s by any field
  - [ ] Add information on adding the ability to create arbitrary search indexes
  - [ ] Add information on using the builder pattern for constructing `Library` instances
  - [ ] Add information on using the visitor pattern to handle persisting `Book`s as efficiently as possible (e.g. in-memory LRU caching and background persistance of `Book`s that have changed)
  - [ ] Add thoughts on fuzzing and coverage testing

## Discussion

The nouns in this assignment are "library" and "book". The verbs are "list", "create", "read", "update", and "delete". The obvious questions are:

1. Will the library data persist over time?
2. If the requirements include persisitence how will it persist?
3. What data is associated with each book?
4. Do we need to build indexes for fast look up over multiple pieces of book data?
5. Does the library require concurrent access?

Based on the limited requirements the package must export at least `Library` and `Book` types. Since the `Library` is likely to have different implementations with different persistence and access protection policies, it makes sense to create it as an `interface` and then provide a basic implementation that demonstrates that functions. The `Book` type is obvious however in any CRUD based data structure, there is usually multiple ways to look up an object. The `Book` is no different. Currently a `Book` is just an ISBN number and a Title and the functions that look up books rely on the ISBN number to find the book.

I'm building this using test-driven development. This first version has tests for both the `Book` and `Library` type that simply fail for now.

### Future Proofing and Maintenance

TODO

## Build, Run, and Test

This project uses a simple Makefile to control building the demo app and library, running the demo app, and testing the library. The simplest thing to do is to execute the following commands in the root directory of the cloned repo:

```
$ make build
$ make test
$ make run
```

To clean up build artifacts, just run:

```
$ make clean
```

## Project layout
```shell
├── bin             //the resulting demo executable is in here
├── cmd             //the demo binary that uses the library is in here
├── Makefile        //the only one make file that can build, test, and run the demo
└── pkg             //the library used in the project
```
## License
Apache 2.0 License

Copyright (c) 2022 David Huseby

[//]: # (badges)

[license-image]: https://img.shields.io/badge/license-Apache2.0-blue.svg
