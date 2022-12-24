// Copyright David Huseby. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package library

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"golang.org/x/exp/maps"
)

// This impl of Library persists Books to disk (D) as JSON encoded data and stores
// them in memory using a map (M) with concurrent access (C) protected using
// atomic.Value assuming that the access pattern is mostly reads.
type DMC struct {

	// the data is stored in memory as map[ISBN]Book
	books	atomic.Value

	// writer lock
	mtx	sync.Mutex

	// file to store the data in
	f	string
}

func DMCBuild(arg any) (Library, error) {
	switch arg.(type) {
	case string:
		str, ok := arg.(string)
		if ok {
			return DMCFromFile(str)
		}
		return nil, fmt.Errorf("arg type of string failed to convert")
	case []byte:
		data, ok := arg.([]byte)
		if ok {
			return DMCFromData(data)
		}
		return nil, fmt.Errorf("arg type of data failed to convert")
	default:
		return nil, fmt.Errorf("unknown arg type for constructing DMC: %T", arg)
	}
}

func DMCFromFile(f string) (Library, error) {
	// construct our library
	var dmc = new(DMC)

	// read the books JSON from disk
	data, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	// unmarshal the data into memory
	books := make([]Book, 0)
	err = json.Unmarshal(data, &books)
	if err != nil {
		return nil, err
	}

	// rearrange the book slice into a map indexed by ISBN
	l := make(map[ISBN]Book)
	for _, book := range books {
		l[book.ID] = book
	}

	// store the data in the struct
	dmc.books.Store(l)

	// remember where the data came from
	dmc.f = f

	return dmc, nil
}

func DMCFromData(d []byte) (Library, error) {
	// construct our library
	var dmc = new(DMC)

	books := make([]Book, 0)

	if len(d) > 0 {
		// unmarshal the data into memory
		err := json.Unmarshal(d, &books)
		if err != nil {
			return nil, err
		}
	}

	// rearrange the book slice into a map indexed by ISBN
	l := make(map[ISBN]Book)
	for _, book := range books {
		l[book.ID] = book
	}

	// store the data in the struct
	dmc.books.Store(l)

	return dmc, nil
}

func (l *DMC) ListAll() ([]Book, error) {
	// do an atomic read of the map
	books := l.books.Load().(map[ISBN]Book)
	return maps.Values(books), nil
}

func (l *DMC) Create(b Book) error {

	// first try to read to see if the book is already in the library
	_, err := l.Read(b.ID)
	if err == nil {
		return fmt.Errorf("Book is already in the library: %s", b)
	}

	// sync with other writers
	l.mtx.Lock()
	defer l.mtx.Unlock()

	// load the current value of the data
	books := l.books.Load().(map[ISBN]Book)

	// create a copy
	updated := make(map[ISBN]Book)
	for k, v := range books {
		updated[k] = v
	}

	// add in the new book
	updated[b.ID] = b

	// store the new updated data
	l.books.Store(updated)

	return nil
}

func (l *DMC) Read(id ISBN) (Book, error) {
	// load the current value of the data
	books := l.books.Load().(map[ISBN]Book)

	// try to get the book
	book, ok := books[id]
	if ok {
		return book, nil
	}

	return Book{}, fmt.Errorf("failed to find the book for %s", id)
}

func (l *DMC) Update(b Book) error {
	// this function will create a new book if it doesn't exist already.
	// if it doesn't exist, it will overwrite it with new data

	// sync with other writers
	l.mtx.Lock()
	defer l.mtx.Unlock()

	// load the current value of the data
	books := l.books.Load().(map[ISBN]Book)

	// create a copy
	updated := make(map[ISBN]Book)
	for k, v := range books {
		updated[k] = v
	}

	// add in the new book
	updated[b.ID] = b

	// store the new updated data
	l.books.Store(updated)

	return nil
}

func (l *DMC) Delete(id ISBN) (Book, error) {
	// first try to read to see if the book is in the library
	book, err := l.Read(id)
	if err != nil {
		return Book{}, fmt.Errorf("Delete cannot find book in library: %s", id)
	}

	// sync with other writers
	l.mtx.Lock()
	defer l.mtx.Unlock()

	// load the current value of the data
	books := l.books.Load().(map[ISBN]Book)

	// create a copy of every Book but the one we're deleting
	updated := make(map[ISBN]Book)
	for k, v := range books {
		if k != id {
			updated[k] = v
		}
	}

	// store the new updated data
	l.books.Store(updated)

	// return the book to the caller
	return book, nil
}

func (l *DMC) Close() error {
	// if the data didn't come from a file...return silently
	if len(l.f) == 0 {
		return nil
	}

	// do an atomic read of the map
	books := l.books.Load().(map[ISBN]Book)

	// marshal the data to JSON
	data, err := json.MarshalIndent(maps.Values(books), "", "  ")
	if err != nil {
		return err
	}

	// write it to a file
	err = os.WriteFile(l.f, data, 0600)
	if err != nil {
		return err
	}

	return nil
}
