// Copyright David Huseby. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package library

import (
	"fmt"
)

// All published books are identified by an ISBN number. I will use this as
// the identifier for looking up books for now.
type ISBN = string

// For now, a book only has a title and an ISBN number
type Book struct {
	ID	ISBN
	Title	string
}

// Stringify the book for pretty printing
func (b Book) String() string {
	return fmt.Sprintf("%s, ISBN:%s", b.Title, b.ID)
}
