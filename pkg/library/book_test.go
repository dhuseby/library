// Copyright David Huseby. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package library

import "testing"

func TestString(t *testing.T) {
	b := Book { "1234567890", "Foo Bar" }
	got := b.String()
	want := "Foo Bar, 1234567890"
	if got != want {
		t.Log("should be", want, "but got", got)
		t.Fail()
	}
}
