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

package dupi

import (
	"bytes"
	"testing"
)

func TestFnames(t *testing.T) {
	s := newFnames()

	a, err := s.addPath("/a")
	if err != nil {
		t.Fatal(err)
	}
	if s.abs(a) != "/a" {
		t.Errorf("bad absa")
	}
	b, err := s.addPath("/b")
	if err != nil {
		t.Fatal(err)
	}
	if s.abs(b) != "/b" {
		t.Errorf("bad absb")
	}
	ab, err := s.addPath("/a/b")
	if err != nil {
		t.Fatal(err)
	}
	if s.abs(ab) != "/a/b" {
		t.Errorf("bad absb")
	}
	cwd, err := s.addPath(".")
	if err != nil {
		t.Fatal(err)
	}
	io := bytes.NewBuffer(nil)
	if err := s.write(io); err != nil {
		t.Fatal(err)
	}
	io = bytes.NewBuffer(io.Bytes())
	ss, err := readFnames(io)
	if err != nil {
		t.Fatal(err)
	}
	if ss.abs(a) != "/a" {
		t.Errorf("bad ss-absa")
	}
	if ss.abs(b) != "/b" {
		t.Errorf("bad ss-absb")
	}
	if ss.abs(ab) != "/a/b" {
		t.Errorf("bad ss-absab")
	}
	bc, err := ss.addPath("/b/c")
	if err != nil {
		t.Fatal(err)
	}
	if ss.abs(bc) != "/b/c" {
		t.Errorf("bad bc path")
	}
	t.Logf("cwd: '%s'\n", ss.abs(cwd))
	t.Logf("%s\n", ss)

}
