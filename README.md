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
  - [x] Pass all tests
- [x] Make demo app fully functional
- [x] Add Future Proofing and Maintenance discussion to the README
  - [x] Add information on abstracting the searching mechanism behind an interface to allow for looking up `Book`s by any field
  - [x] Add information on adding the ability to create arbitrary search indexes
  - [x] Add information on using the builder pattern for constructing `Library` instances
  - [x] Add information on using the visitor pattern to handle persisting `Book`s as efficiently as possible (e.g. in-memory LRU caching and background persistance of `Book`s that have changed)
  - [x] Add thoughts on fuzzing and coverage testing

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

The main thing that a developer of a library like this needs to pay attention to is the different "policies" that may be needed in the future. Since you never know what those policies are, it's best to have well placed interfaces that abstract away key behavior. For instance, this implementation currently abstracts away the Library impl itself so that concrete impls can have different behaviors. The one impl I provide assumes that the data persisitence is to a JSON file on disk and that there is potentially multiple concurrent readers and writers but that the concurrency pattern is "mostly read".

I could see an impl that instead uses a local/remote database as the persistence later. In that case the concurrency is handle by the database layer most likely so the impl in the library can be fairly simple.

If the number of concrete impls begins to grow substantially, it might make sense to create interfaces for persistence and synchronization behavior so that there's only one concrete `Library` impl that relies on concrete impls of `Persisitence` and `Synchronization` to handle the CRUD of each `Book` and the concurrency respectively. That would then allow a combinatorial explosion of combinations of impls of those policies. One interesting approach to handling persistence is to use the visitor pattern combined with an LRU read cache. If the access behavior is mostly read, then an LRU cache hides any latency of accessing the underlying persistent storage. When an Book gets flushed from the cache, it must be checked to see if it has been updated and then stored before removal. When the `Library.Close()` gets called, the visitor pattern can handle the persistence of just the cached Books that have writes pending. This makes the logic of persistence extensible and ready to adapt to new data types other than `Book`.

One thing that is common in libraries like this is an abstracted indexing and searching mechanism. That again would be another policy abstracted away behind interfaces so that there's only one concrete `Library` that aggregates together all of the impls for the different policies. I would still keep the `Library` interface just to keep impl details politely behind the scenes.

As the number of interfaces and concrete impls grows, the builder pattern becomes the best way--IMHO--to clean up the initialization of Library with the desired policies:

```
l, err := LibraryBuilder
	.WithSynchronization(MostlyRead)
	.WithPersistence(MySQL)
	.WithQueryEngine(DynamicIndexing)
	.Build()
```

Finally, the one thing I wanted to add but didn't was a set of tests that use the `go fuzz` tooling to try to break the library by feeding it random ID's and Titles and doing random operations. I think fuzzing is a vastly underutilized tool. It is usually remarkably good at finding flaws in utility libraries like this one.

One thing I want to note is that this repo is organized to be integrated into any CI/CD pipeline. The Makefile uses `-ldflags` to add versioning information into the compiled library for automated tooling that looks at binary outputs from the build process that is going into the running environment. This is a best practice in the SBOM world we live in now. Unfortunately none of this is truly cryptographic; it's just hashes of the git tree. An ideal solution would use code signing and key provenance with corroborated identities to provide true traceability of the binary artifacts. When combined with reproducible builds this utility library would fit nicely even in the most rigorously change-managed and secure software projects.

Cheers!

## Build, Run, and Test

This project uses a simple Makefile to control building the demo app and library, running the demo app, and testing the library. The simplest thing to do is to execute the following commands in the root directory of the cloned repo:

```
$ make
$ make test
$ make run
```

To clean up build artifacts, just run:

```
$ make clean
```

If you have your .gdbinit set up correctly to allow auto-loading of local .gdbinit files, you can also debug the demo app by running:

```
$ make debug
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
