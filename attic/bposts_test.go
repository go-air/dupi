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
	"fmt"
	"math/rand"
	"testing"
	"unsafe"

	"github.com/scott-cotton/dupi/attic/ibloom"
	"github.com/scott-cotton/dupi/attic/trigram"
)

func TestTridat(t *testing.T) {
	var tds bposts
	td := &tds
	td.Bloomer = ibloom.NewFnv32(10000, 3)
	td.Set = td.Bloomer.Set()
	N := 100
	added := make([]uint32, N)
	for i := 0; i < N; i++ {
		u := rand.Uint32()
		td.AddPost(u)
		added[i] = u
	}
	for _, u := range added {
		if !td.HasPost(u) {
			t.Errorf("not added: %d\n", u)
		}
	}
	m := make(map[uint32]bool, N)
	for _, u := range added {
		m[u] = true
	}
	nErrs := 0
	for i := 0; i < 8192; i++ {
		u := rand.Uint32()
		if m[u] {
			continue
		}
		if td.HasPost(u) {
			nErrs++
		}
	}
	if nErrs > 200 {
		t.Errorf("error rate %d/8192\n", nErrs)
	}

	var tda [trigram.NumTrigrams]bposts
	fmt.Printf("allocated %d (%d bytes)\n", len(tda), len(tda)*int(unsafe.Sizeof(tda[0])))
}
