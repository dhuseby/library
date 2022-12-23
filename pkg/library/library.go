// Copyright David Huseby. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package library

// The library needs is an interface so there can be multiple different
// concrete impls with different policies for persistence and access. This
// interface contains the complete set of operations defined in the problem
// statement.
type Library interface {

	// lists all books
	ListAll() ([]*Book, error)

	// create a new book in the library
	Create(b *Book) error

	// looks up a book by ISBN
	Read(id ISBN) (*Book, error)

	// update a book, acts like Create if the book isn't in the library
	Update(b *Book) error

	// deletes a book from the library and returns it
	Delete(id ISBN) (*Book, error)

	// clean up the references to external resources
	Close() error
}
