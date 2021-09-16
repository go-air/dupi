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

// Package dupi provides a library for exploring duplicate
// data in large sets of documents.
package dupi

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Doc struct {
	Path  string
	Start uint32
	End   uint32
	Dat   []byte `json:"-"`
}

func NewDoc(path, body string) *Doc {
	return &Doc{
		Path: path,
		Dat:  []byte(body)}
}

func (doc *Doc) Load() error {
	var (
		f   *os.File
		err error
	)

	f, err = os.Open(doc.Path)
	if err != nil {
		return err
	}

	if doc.Start == 0 && doc.End == 0 {
		doc.Dat, err = ioutil.ReadAll(f)
		if err != nil {
			return fmt.Errorf("readall: %w", err)
		}
		f.Close()
		doc.End = uint32(len(doc.Dat))
	} else {
		_, err = f.Seek(int64(doc.Start), io.SeekStart)
		if err != nil {
			return fmt.Errorf("seek: %w", err)
		}
		doc.Dat = make([]byte, doc.End-doc.Start)
		_, err = f.ReadAt(doc.Dat, int64(doc.Start))
		if err != nil {
			return fmt.Errorf("readat len=%d at=%d: %w\n", len(doc.Dat), doc.Start, err)
		}
		f.Close()
	}
	return nil
}
