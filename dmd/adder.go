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
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
)

type Adder struct {
	path      string
	flushed   uint32
	flushRate uint32
	buf       []fields
}

func NewAdder(root string, fr int) (*Adder, error) {
	ufr := uint32(fr)
	path := filepath.Join(root, "dmd")
	adder := &Adder{path: path, flushRate: ufr}
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		adder.buf = make([]fields, 1, fr)
		adder.flushed = 0
		return adder, nil
	}
	if err != nil {
		return nil, err
	}
	adder.flushed = uint32(fi.Size() / rcdSize)
	if adder.flushed == 0 {
		adder.buf = make([]fields, 1, fr)
	} else {
		adder.buf = make([]fields, 0, fr)
	}
	return adder, nil
}

func (t *Adder) Add(fid, start, end uint32) (uint32, error) {
	n := uint32(len(t.buf))
	if n == t.flushRate {
		err := t.flush()
		if err != nil {
			return 0, err
		}
	}
	n = uint32(len(t.buf))
	t.buf = append(t.buf, fields{fid, start, end})
	return n + t.flushed, nil
}

func (t *Adder) Last() uint32 {
	return t.flushed + uint32(len(t.buf)) - 1
}

func (t *Adder) Close() error {
	if len(t.buf) == 0 {
		return nil
	}
	return t.flush()
}

func (t *Adder) flush() error {
	f, err := os.OpenFile(t.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	for i := range t.buf {
		fields := &t.buf[i]
		err := binary.Write(f, binary.BigEndian, fields)
		if err != nil {
			return err
		}
	}
	t.flushed += uint32(len(t.buf))
	t.buf = t.buf[:0]
	return nil
}
