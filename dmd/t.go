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

// Package dmd maps document, offset pairs to internal document ids.
package dmd

import (
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
)

type T struct {
	path string
	file *os.File
}

const rcdSize = 12

func New(root string) (*T, error) {
	res := &T{path: filepath.Join(root, "dmd")}
	var err error
	res.file, err = os.Open(res.path)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *T) NumDocs() (uint64, error) {
	fi, err := t.file.Stat()
	if err != nil {
		return 0, err
	}
	return uint64(fi.Size()) / rcdSize, nil
}

func (t *T) Lookup(did uint32) (fid, start, end uint32, err error) {
	f := t.file
	_, err = f.Seek(int64(did)*rcdSize, 0)
	if err != nil {
		return
	}
	var buf [rcdSize]byte
	_, err = io.ReadFull(f, buf[:])
	if err != nil {
		return
	}
	fid = binary.BigEndian.Uint32(buf[0:4])
	start = binary.BigEndian.Uint32(buf[4:8])
	end = binary.BigEndian.Uint32(buf[8:rcdSize])
	return
}

func (t *T) Close() error {
	return t.file.Close()
}
