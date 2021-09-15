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

package ibloom

type fnv32 struct {
	k      int
	mbytes int
}

// NewFnv32 creates an ibloom filter with
// mbytes bytes of data and k hash functions.
// The k hash functions are fnv hash based
// and cost just a multiply and xor on top
// of the cost of hashing the input value
// to achieve k variations.
func NewFnv32(mbytes, k int) Filter {
	b := &fnv32{k: k, mbytes: mbytes}
	return b
}

func (f *fnv32) Config() *Config {
	return &Config{Name: "fnv32", K: f.k, M: f.mbytes * 8}
}

const fnv32Prime = 0x01000193

func fnv32Hash(hop uint32) uint32 {
	hh := uint32(1)
	hh ^= hop & 0xff
	hh *= fnv32Prime
	hh ^= (hop >> 8) & 0xff
	hh *= fnv32Prime
	hh ^= (hop >> 16) & 0xff
	hh *= fnv32Prime
	hh ^= (hop >> 24) & 0xff
	hh *= fnv32Prime
	return hh
}

func (b *fnv32) K() int {
	return b.k
}

func (b *fnv32) M() int {
	return b.mbytes * 8
}

func (b *fnv32) Set() []byte {
	return make([]byte, b.mbytes)
}

func (b *fnv32) Ones(post uint32, fn func(int) bool) bool {
	hh := fnv32Hash(post)
	k := uint32(b.k)
	for i := uint32(0); i < k; i++ {
		hb := hh ^ (i + 'a')
		hb *= fnv32Prime
		hb = clamp(hb, uint32(b.mbytes))
		if !fn(int(hb)) {
			return false
		}
	}
	return true
}

func (f *fnv32) PutOnes(dst []uint32, v uint32) {
	hh := fnv32Hash(v)
	k := uint32(f.k)
	for i := uint32(0); i < k; i++ {
		hb := hh ^ (i + 'a')
		hb *= fnv32Prime
		hb = clamp(hb, uint32(f.mbytes))
		dst[i] = hb
	}
}

func clamp(h, n uint32) uint32 {
	bits := n * 8
	return h % bits
}

func (b *fnv32) Has(d []byte, v uint32) bool {
	return b.Ones(v, func(i int) bool {
		di, dm := i/8, i%8
		return d[di]&(1<<dm) != 0
	})
}

func (b *fnv32) Add(d []byte, v uint32) {
	b.Ones(v, func(i int) bool {
		di, dm := i/8, i%8
		d[di] |= 1 << dm
		return true
	})
}
