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
	"fmt"
	"io"
	"os"
)

type poster struct {
	link int64
	head int64

	current uint32
	total   uint32
	posts   []byte
	buf     []byte
	blot    uint16
}

const (
	flushRate = 512
)

func (p *poster) initCommon(blot uint16) {
	p.link = -1
	p.head = -1
	p.posts = make([]byte, flushRate+2*binary.MaxVarintLen64+8)
	p.buf = p.posts[:binary.MaxVarintLen64]
	p.posts = p.posts[binary.MaxVarintLen64:]
	p.posts = p.posts[:0]
	p.blot = blot
}

func (p *poster) initCreate(w io.Writer, blot uint16) error {
	p.initCommon(blot)
	_, err := p.writeVarint64(w, p.head)
	if err != nil {
		return err
	}
	var buf [4]byte
	_, err = w.Write(buf[:])
	return err
}

func (p *poster) initAppend(r io.ByteReader, blot uint16) error {
	p.initCommon(blot)
	var err error
	p.head, err = binary.ReadVarint(r)
	if err != nil {
		return err
	}
	var ttl int64
	ttl, err = binary.ReadVarint(r)
	if uint64(ttl)&0xffffffff != uint64(ttl) {
		return fmt.Errorf("invalid total: %d", ttl)
	}
	p.total = uint32(ttl)
	var buf [4]byte
	for i := range buf {
		buf[i], err = r.ReadByte()
		if err != nil {
			return err
		}
	}
	p.current = binary.BigEndian.Uint32(buf[:])

	return err
}

func (p *poster) AddPost(v uint32, f *os.File) error {
	if p.current == v && p.total != 0 {
		return nil
	}

	delta := v - p.current
	p.current = v
	p.total++
	n := binary.PutUvarint(p.buf, uint64(delta))
	p.posts = append(p.posts, p.buf[:n]...)
	if len(p.posts) >= flushRate {
		return p.flushTo(f)
	}
	return nil
}

func (p *poster) String() string {
	return fmt.Sprintf("poster.%x hd=%d link=%d count=%d",
		p.blot, p.head, p.link, p.total)

}

func (p *poster) readToLink(f *os.File) error {
	if p.head == -1 {
		return nil
	}
	if p.link != -1 {
		panic(p.link)
	}
	buf := make([]byte, flushRate+binary.MaxVarintLen64+8+10)
	head := p.head
	for {
		v, n, err := readVarintAt(f, head)
		if err != nil {
			return err
		}
		if v > int64(len(buf)) {
			buf = make([]byte, v*2)
		}
		if v <= 8 {
			return fmt.Errorf("%s: invalid length of flushed post buffer: %d", p, v)
		}
		_, err = f.ReadAt(buf[:v], head+n)
		if err != nil {
			return err
		}
		link := int64(binary.BigEndian.Uint64(buf[v-8 : v]))
		if link == -1 {
			p.link = head + n + v - 8
			return nil
		} else {
			head = link
		}
	}
}

func (p *poster) writeVarint64(w io.Writer, v int64) (int, error) {
	n := binary.PutVarint(p.buf, v)
	_, err := w.Write(p.buf[:n])
	return n, err
}

func (p *poster) flushTo(f *os.File) error {
	if len(p.posts) == 0 {
		return nil
	}

	fi, err := f.Stat()
	if err != nil {
		return err
	}
	end := fi.Size()
	if p.link != -1 {
		binary.BigEndian.PutUint64(p.buf[:8], uint64(end))
		_, err = f.WriteAt(p.buf[:8], p.link)
		if err != nil {
			return err
		}
	} else {
		p.head = end
	}
	_, err = f.Seek(end, io.SeekStart)
	if err != nil {
		return err
	}
	var z [8]byte
	var nolink = int64(-1)
	binary.BigEndian.PutUint64(z[:], uint64(nolink))
	p.posts = append(p.posts, z[:]...)

	lenSz, err := p.writeVarint64(f, int64(len(p.posts)))
	if err != nil {
		return err
	}
	zz, err := f.Write(p.posts)
	if err != nil {
		return err
	}
	if zz != len(p.posts) {
		panic(z)
	}
	p.link = end + int64(lenSz) + int64(len(p.posts)) - 8
	p.posts = p.posts[:0]
	return nil
}
