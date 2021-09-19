package dmd

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestAdder(t *testing.T) {
	tmp, err := ioutil.TempDir(".", "dmd.test.")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(tmp)
	}()
	adder, err := NewAdder(tmp, 16384)
	if err != nil {
		t.Fatal(err)
	}
	did0, err := adder.Add(1, 2, 3)
	did1, err := adder.Add(2, 3, 4)
	did2, err := adder.Add(3, 0, 11)
	did3, err := adder.Add(1, 3, 222)
	err = adder.Close()
	if err != nil {
		t.Fatal(err)
	}
	adder, err = NewAdder(tmp, 16384)
	did4, err := adder.Add(2, 4, 7)
	err = adder.Close()
	if err != nil {
		t.Fatal(err)
	}
	dmd, err := New(tmp)
	defer dmd.Close()
	if err != nil {
		t.Fatal(err)
	}
	a, b, c, err := dmd.Lookup(did0)
	if err != nil {
		t.Fatal(err)
	}
	if a != 1 || b != 2 || c != 3 {
		t.Error(err)
	}
	a, b, c, err = dmd.Lookup(did1)
	if err != nil {
		t.Fatal(err)
	}
	if a != 2 || b != 3 || c != 4 {
		t.Error(err)
	}
	_ = did2
	_ = did3
	a, b, c, err = dmd.Lookup(did4)
	if err != nil {
		t.Fatal(err)
	}
	if a != 2 || b != 4 || c != 7 {
		t.Error(err)
	}

}
