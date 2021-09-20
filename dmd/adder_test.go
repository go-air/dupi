// Copyright 2021 the Dupi authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
