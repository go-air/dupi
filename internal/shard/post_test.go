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
	"bufio"
	"encoding/binary"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

func TestPosts(t *testing.T) {
	iix, d, err := postFiles()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		iix.Close()
		d.Close()
		os.Remove(iix.Name())
		os.Remove(d.Name())
	}()
	poster := &poster{}
	if err := poster.initCreate(iix, 0x7); err != nil {
		t.Fatal(err)
	}
	ps := gen(1311)
	for _, p := range ps {
		if err := poster.AddPost(p, d); err != nil {
			t.Fatal(err)
		}
	}
	if err := poster.flushTo(d); err != nil {
		t.Fatal(err)
	}
	if err := iix.Truncate(0); err != nil {
		t.Fatal(err)
	}
	if _, err := poster.writeVarint64(iix, poster.head); err != nil {
		t.Fatal(err)
	}
	if _, err := poster.writeVarint64(iix, int64(poster.total)); err != nil {
		t.Fatal(err)
	}

	if _, err := iix.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	if _, err := d.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	rdr := bufio.NewReader(iix)
	hd, err := binary.ReadVarint(rdr)
	if err != nil {
		t.Fatal(err)
	}
	ct, err := binary.ReadVarint(rdr)
	if err != nil {
		t.Fatal(err)
	}
	_ = ct
	pr := newPosts(hd)
	qs := make([]uint32, 0, len(ps))
	for {
		did, err := pr.next(d)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		qs = append(qs, did)
	}
	if len(qs) != len(ps) {
		t.Errorf("got %d docids put %d", len(qs), len(ps))
	}
	for i, rdid := range qs {
		if i == len(ps) {
			break
		}
		pdid := ps[i]
		if rdid != pdid {
			t.Errorf("got doc %d expected %d", rdid, pdid)
		}
	}
}

func postFiles() (*os.File, *os.File, error) {
	iix, err := ioutil.TempFile(".", "post.iix")
	if err != nil {
		return nil, nil, err
	}
	d, err := ioutil.TempFile(".", "post")
	if err != nil {
		iix.Close()
		return nil, nil, err
	}
	return iix, d, nil
}

func gen(n int) []uint32 {
	res := make([]uint32, n)
	cur := rand.Uint32() >> 4
	var delta uint32
	for i := 0; i < n; i++ {
		delta = uint32(rand.Intn(117)) + 1
		cur += delta
		res[i] = cur
	}
	return res
}
