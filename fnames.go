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
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type fname struct {
	parent   uint32
	name     string
	children map[string]uint32
}

type fnames struct {
	d []fname
}

func newFnames() *fnames {
	res := &fnames{
		d: make([]fname, 1, 128)}
	return res
}

func (s *fnames) String() string {
	w := bytes.NewBuffer(nil)
	for i := range s.d {
		abs := s.abs(uint32(i))
		fmt.Fprintf(w, "%s\n", abs)
	}
	return w.String()
}

func (s *fnames) abs(v uint32) string {
	var parts []string
	n := &s.d[v]
	for n != &s.d[0] {
		parts = append(parts, n.name)
		n = &s.d[n.parent]
	}
	N := len(parts)
	H := N / 2
	N--
	for i := 0; i < H; i++ {
		parts[i], parts[N-i] = parts[N-i], parts[i]
	}
	psep := string(os.PathSeparator)
	return psep + strings.Join(parts, psep)
}

func (s *fnames) addPath(path string) (uint32, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("abs\n")
		return 0, err
	}
	var (
		parts []string
		p     = abs
	)
	parts = strings.Split(p, string(os.PathSeparator))
	if parts[0] == "" {
		parts = parts[1:]
	}
	var (
		parent = uint32(0)
		m      map[string]uint32
		child  uint32
	)
	for i := range parts {
		p = parts[i]
		m = s.d[parent].children
		if m == nil {
			m = make(map[string]uint32)
			s.d[parent].children = m
		}
		child = m[p]
		if child == 0 {
			child = uint32(len(s.d))
			m[p] = child
			s.d = append(s.d, fname{name: p, parent: parent})
		}
		parent = child
	}
	return parent, nil
}

func readUvarint32(r *bufio.Reader) (uint32, error) {
	v64, err := binary.ReadUvarint(r)
	if err != nil {
		return 0, err
	}
	if v64&0xffffffff != v64 {
		return 0, fmt.Errorf("overflow 32bit varint: %d", v64)
	}
	return uint32(v64), nil
}

func readFnames(r io.Reader) (*fnames, error) {
	return breadFnames(bufio.NewReader(r))
}

func breadFnames(r *bufio.Reader) (*fnames, error) {
	v, err := readUvarint32(r)
	if err != nil {
		return nil, err
	}
	s := &fnames{}
	s.d = make([]fname, v)
	var buf []byte
	for i := range s.d {
		v, err = readUvarint32(r)
		if err != nil {
			return nil, err
		}
		if v > uint32(len(buf)) {
			buf = make([]byte, v)
		}
		_, err = io.ReadFull(r, buf[:v])
		if err != nil {
			return nil, err
		}
		fname := &s.d[i]
		fname.name = string(buf[:v])
		v, err = readUvarint32(r)
		if err != nil {
			return nil, err
		}
		fname.parent = v
		pname := &s.d[fname.parent]
		if pname.children == nil {
			pname.children = make(map[string]uint32)
		}
		pname.children[fname.name] = uint32(i)
	}
	return s, nil
}

func writeUvarint32(w io.Writer, v uint32) error {
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], uint64(v))
	_, err := w.Write(buf[:n])
	return err
}

func (s *fnames) write(w io.Writer) error {
	err := writeUvarint32(w, uint32(len(s.d)))
	if err != nil {
		return err
	}
	for i := range s.d {
		fname := &s.d[i]
		err := writeUvarint32(w, uint32(len(fname.name)))
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(fname.name))
		if err != nil {
			return err
		}
		err = writeUvarint32(w, uint32(fname.parent))
		if err != nil {
			return err
		}
	}
	return nil
}
