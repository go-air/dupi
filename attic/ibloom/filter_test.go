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

import (
	"math/rand"
	"testing"
)

func TestEstimate(t *testing.T) {
	n := 1000.0
	p := 0.15
	m, k := Estimate(n, p)
	t.Logf("p=%.2f n=%.2f m=%d k=%d\n", p, n, m, k)
	bloom := NewFnv32(m/8, k)
	set := bloom.Set()
	added := make(map[uint32]bool, 1000)
	for i := 0; i < 100; i++ {
		v := rand.Uint32()
		added[v] = true
		bloom.Add(set, v)
	}
	for v := range added {
		if !bloom.Has(set, v) {
			t.Errorf("added %v but doesn't have it.\n", v)
		}
	}
	errs := 0
	for i := 0; i < 10000; i++ {
		v := rand.Uint32()
		if added[v] {
			i--
			continue
		}
		if bloom.Has(set, v) {
			errs++
		}
	}
	t.Logf("%d/10000 errors.\n", errs)
}
