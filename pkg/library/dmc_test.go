// Copyright David Huseby. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package library

import (
	"fmt"
	"testing"

	"golang.org/x/exp/slices"
)

func testDMC() (Library, error) {
	data := []byte("[ { \"ID\": \"1234567890\", \"Title\": \"Foo\" }, { \"ID\": \"0987654321\", \"Title\": \"Bar\" } ]")
	return Build(DMC{}, data)
}

func TestListAll(t *testing.T) {
	l, err := testDMC()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	got, err := l.ListAll()
	want := []Book{ Book{ "1234567890", "Foo"}, Book{ "0987654321", "Bar" } }

	for _, g := range got {
		idx := slices.IndexFunc(want, func(w Book) bool { return w.ID == g.ID })
		if idx == -1 {
			t.Log(fmt.Errorf("failed to ListAll. %s isn't in the want list", g))
			t.Fail()
		}
	}
}

func TestCreate(t *testing.T) {
	l, err := Build(DMC{}, []byte{})
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	b := Book{ "1234567890", "Foo" }
	err = l.Create(b)

	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestRead(t *testing.T) {
	l, err := testDMC()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	want := Book{ "1234567890", "Foo"}
	got, err := l.Read(ISBN("1234567890"))

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if got != want {
		t.Log(fmt.Errorf("read unexpected Book value: %s != %s", got, want))
		t.Fail()
	}
}

func TestUpdate(t *testing.T) {
	l, err := testDMC()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	want := Book{ "1234567890", "Blah"}
	err = l.Update(want)
	got, err := l.Read(ISBN("1234567890"))

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if got != want {
		t.Log(fmt.Errorf("failed to update Book value: %s != %s", got, want))
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	l, err := testDMC()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	id := ISBN("1234567890")
	want := Book{ id, "Foo"}
	got, err := l.Delete(id)

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if got != want {
		t.Log(fmt.Errorf("unexpected Book value from Delete: %s != %s", got, want))
		t.Fail()
	}

	_, err = l.Read(id)
	if err == nil {
		t.Log(fmt.Errorf("failed to delete Book for %s", id))
		t.Fail()
	}
}
