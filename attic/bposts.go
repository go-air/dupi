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

package attic

import (
	"github.com/go-air/dupi/attic/ibloom"
)

type bposts struct {
	Root    string
	Total   uint32
	Post    uint32
	Bloomer ibloom.Filter
	Set     []byte
}

func bpostOnes(b ibloom.Filter) *bposts {
	td := &bposts{Bloomer: b}
	td.Set = b.Set()
	for i := range td.Set {
		td.Set[i] = 0xff
	}
	return td
}

func (t *bposts) init(root string, bloom ibloom.Filter) {
	t.Root = root
	t.Bloomer = bloom
	t.Set = bloom.Set()
}

func (t *bposts) Intersect(o *bposts) *bposts {
	for i := range t.Set {
		t.Set[i] &= o.Set[i]
	}
	return t
}

func (t *bposts) Union(o *bposts) *bposts {
	for i := range t.Set {
		t.Set[i] |= o.Set[i]
	}
	return t
}

func (t *bposts) IsZero() bool {
	for _, b := range t.Set {
		if b != 0 {
			return false
		}
	}
	return true
}

func (t *bposts) HasPost(post uint32) bool {
	return t.Bloomer.Has(t.Set, post)
}

func (t *bposts) Ones(dst []int) []int {
	for i, b := range t.Set {
		if b == 0 {
			continue
		}
		for j := uint(0); j < 8; j++ {
			if b&(1<<j) != 0 {
				dst = append(dst, i*8+int(j))
			}
		}
	}
	return dst
}

func (t *bposts) AddPost(post uint32) error {
	if t.Post == post && t.Total != 0 {
		return nil
	}
	t.Post = post
	t.Total++
	t.Bloomer.Add(t.Set, post)
	return nil
}
