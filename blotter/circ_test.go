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

import "testing"

func TestCirc(t *testing.T) {
	circ := NewCirc(3)
	circ.Blot([]byte("a"))
	circ.Blot([]byte("b"))
	h := circ.Blot([]byte("c"))
	j := circ.Blot([]byte("xyzabc"))
	circ.Blot([]byte("a"))
	circ.Blot([]byte("b"))
	g := circ.Blot([]byte("c"))
	if g != h {
		t.Errorf("circular shifting didn't cycle")
	}
	if g == j {
		t.Errorf("impossible hash collision")
	}
}
