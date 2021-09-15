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

package shard

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type Posts struct {
	i       int
	nextpos int64
	current uint32
	docids  []uint32
	buf     []byte
	vbuf    []byte
}

func readVarintAt(r io.ReaderAt, p int64) (int64, int64, error) {
	var (
		buf [binary.MaxVarintLen64]byte
		n   int64
		err error
	)
	for {
		_, err = r.ReadAt(buf[n:n+1], p+n)
		if err != nil {
			return 0, n, err
		}
		n++
		if buf[n-1] < 0x80 {
			break
		}
		if n == binary.MaxVarintLen64 {
			return 0, n, errors.New("read varint overflow")
		}
	}
	v, m := binary.Varint(buf[:n])
	if m <= 0 {
		return 0, 0, fmt.Errorf("%#v binary.Varint gave n=%d v=%d\n", buf[:n], m, v)
	}
	return v, n, nil
}

func newPosts(head int64) *Posts {
	res := &Posts{}
	res.init(head)
	return res
}

func (p *Posts) init(head int64) {
	p.i = 0
	p.nextpos = head
	p.docids = make([]uint32, 0, flushRate)
	p.buf = make([]byte, flushRate+2*binary.MaxVarintLen64+8)
	p.vbuf = p.buf[:binary.MaxVarintLen64]
	p.buf = p.buf[binary.MaxVarintLen64:]
}

func (p *Posts) next(r io.ReaderAt) (uint32, error) {
	if p.i == len(p.docids) {
		if p.nextpos == -1 {
			return 0, io.EOF
		}
		err := p.readNext(r)
		if err != nil {
			return 0, err
		}
	}
	res := p.docids[p.i]
	p.i++
	return res, nil
}

func (p *Posts) readNext(r io.ReaderAt) error {
	v, n, err := readVarintAt(r, p.nextpos)
	if err != nil {
		return err
	}
	_, err = r.ReadAt(p.buf[:v], p.nextpos+n)
	if err != nil {
		return err
	}
	var (
		t, i int
		d    uint64
		buf  = p.buf[:v]
	)
	p.docids = p.docids[:0]

	for {
		d, t = binary.Uvarint(buf[i:])
		if t <= 0 {
			return fmt.Errorf("error decoding uvarint delta: %#v", buf[i:])
		}
		i += t
		if (d & 0xffffffff) != d {
			return fmt.Errorf("32bit uvarint overflow: %x", d)
		}
		p.current += uint32(d)
		p.docids = append(p.docids, p.current)
		if i > len(buf)-8 {
			return fmt.Errorf("misaligned")
		}
		if i == len(buf)-8 {
			break
		}
	}
	p.nextpos = int64(binary.BigEndian.Uint64(buf[i:]))
	p.i = 0
	return nil
}
