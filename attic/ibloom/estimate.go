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

import "math"

// Estimate gives estimates of the size in bits
// (m) and the number of hash functions (k) which
// can represent sets of size 'n' with false positive
// rate 'fpr'.  See Wikipedia for Bloom
// filters.
//
// Estimate rounds m up to be byte-aligned.
func Estimate(n, fpr float64) (m int, k int) {
	a := n * math.Abs(math.Log(fpr))
	a /= math.Ln2 * math.Ln2
	m = int(math.Round(a))
	m += 8 - (m % 8)
	a = float64(m)
	k = int(math.Round(a / n * math.Ln2))
	return
}
