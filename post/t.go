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

// Package post provides a data structure coupling
// dupi blots with dupi internal document ids.
package post

// a post is a tuple of document id, blot
type T uint64

func (p T) Docid() uint32 {
	return uint32(p >> 32)
}

func (p T) Blot() uint32 {
	return uint32(p)
}

func (p T) Split() (docid, blot uint32) {
	return p.Docid(), p.Blot()
}

func Make(docid, blot uint32) T {
	return T(uint64(docid)<<32 | uint64(blot))
}
