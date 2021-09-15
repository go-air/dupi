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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestIndexQueryTrivial(t *testing.T) {
}

func TestIndexDocid2Doc(t *testing.T) {
	tmp, err := ioutil.TempDir(".", "dupi")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	root := filepath.Join(tmp, "dupi")

	idxr, err := CreateIndexer(root, 4, 10)
	if err != nil {
		t.Fatal(err)
	}
	msg := "We need atleast 10 tokens for this to work sensibly."
	msg2 := "We need atleast 10 tokens for this to work sensibly, really"
	doc := &Doc{
		Path:  "/test/msg",
		Start: 0,
		End:   uint32(len(msg)),
		Dat:   []byte(msg)}
	if err := idxr.Add(doc); err != nil {
		t.Fatal(err)
	}
	doc = &Doc{
		Path:  "/test/msg2",
		Start: 0,
		End:   uint32(len(msg2)),
		Dat:   []byte(msg2)}
	if err := idxr.Add(doc); err != nil {
		t.Fatal(err)
	}
	if err := idxr.Close(); err != nil {
		t.Fatal(err)
	}
	idx, err := OpenIndex(root)
	if err != nil {
		t.Fatal(err)
	}
	defer idx.Close()
	var rdoc Doc
	err = idx.docid2Doc(1, &rdoc)
	if err != nil {
		t.Fatal(err)
	}
	// nb Add(doc) clobbers path.
	if rdoc.Path != "/test/msg" {
		t.Errorf("path got %s want %s", rdoc.Path, doc.Path)
	}
	if rdoc.Start != doc.Start {
		t.Errorf("start got %d want %d", rdoc.Start, doc.Start)
	}
	if rdoc.End != uint32(len(msg)) {
		t.Errorf("end got %d want %d", rdoc.End, len(msg))
	}
}
