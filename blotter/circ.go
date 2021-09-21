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

package blotter

import (
	"bytes"
	"hash"
	"hash/fnv"
)

type Circ struct {
	hashes []uint32
	hash   uint32
	fn     hash.Hash32
	i      int
}

func NewCirc(n int) *Circ {
	return &Circ{hashes: make([]uint32, n), hash: 0, fn: fnv.New32()}
}

func (c *Circ) Config() *Config {
	return &Config{SeqLen: len(c.hashes), Interleave: 1}
}

func (c *Circ) Interleaving() int {
	return 1
}

func (c *Circ) Blot(word []byte) uint32 {
	fn := c.fn
	fn.Reset()
	fn.Write(bytes.ToLower(word))
	h := fn.Sum32()
	c.hash ^= c.hashes[c.i]
	c.hashes[c.i] = h
	c.hash ^= h
	c.i++
	if c.i == len(c.hashes) {
		c.i = 0
	}
	return c.hash
}
