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
	"fmt"
	"os"
	"sort"
)

type Index struct {
	id       uint32
	path     string
	heads    [1 << 16]int64
	counts   [1 << 16]uint32
	perm     [1 << 16]uint16
	postFile *os.File
}

func (x *Index) Init(path string) error {
	x.path = path
	if err := x.readIix(); err != nil {
		return err
	}
	f, err := os.Open(x.path)
	if err != nil {
		return err
	}
	x.postFile = f
	return nil
}

func (x *Index) Close() error {
	return x.postFile.Close()
}

func (x *Index) ReadStateFor(blot uint16) *ReadState {
	/*
		        broken still
			refCount := x.counts[blot]
			i := sort.Search(1<<16, func(i int) bool {
				return x.counts[x.perm[i]] > refCount
			})
			if i == len(x.perm) {
				panic("internal error in perm/count")
			}
			if x.perm[i] != blot {
				fmt.Printf("blot %x i %d perm[%d] %x count %d ref %d %#v\n", blot, i, i, x.perm[i], x.counts[x.perm[i]], refCount, x.perm[i-1:i+2])
				panic("internal error2 in perm/count")
			}
	*/
	return x.ReadStateForBlotAt(blot, 0) //uint16(i))
}

func (x *Index) ReadStateForBlotAt(blot, at uint16) *ReadState {
	res := &ReadState{}
	res.Posts = newPosts(x.heads[blot])
	res.Shard = x.id
	res.Blot = blot
	res.At = at
	res.Total = x.counts[at]
	res.rdr = x.postFile
	return res
}

func (x *Index) BlotAt(i uint16) (blot uint16) {
	blot = x.perm[i]
	return
}

func (x *Index) ReadStateAt(i uint16) *ReadState {
	return x.ReadStateForBlotAt(x.BlotAt(i), i)
}

func (x *Index) Count(blot uint32) uint32 {
	return x.counts[blot]
}

func (x *Index) NumPosts() uint64 {
	var ttl uint64
	for _, ct := range x.counts {
		ttl += uint64(ct)
	}
	return ttl
}

func (x *Index) NumBlots() uint64 {
	return 1 << 16
}

func (x *Index) SosDiffs(avg float64) float64 {
	var ttl float64
	for _, ct := range x.counts {
		d := avg - float64(ct)
		ttl += d * d
	}
	return ttl
}

func (x *Index) readIix() error {
	f, err := os.Open(fmt.Sprintf("%s.iix", x.path))
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for i := range x.heads {
		hd, err := binary.ReadVarint(r)
		if err != nil {
			return err
		}
		x.heads[i] = hd
		ct, err := binary.ReadVarint(r)
		if ct&0xffffffff != ct {
			return fmt.Errorf("not a 32 bit count")
		}
		if err != nil {
			return err
		}
		x.counts[i] = uint32(ct)
		x.perm[i] = uint16(i)
	}
	sort.Slice(x.perm[:], func(i, j int) bool {
		return x.counts[x.perm[i]] > x.counts[x.perm[j]]
	})
	return nil
}
