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

package dupi

// Blot represents a piece of a query or extraction.
// The field Blot gives the blot which was witnessed
// in the docs specified in the field Docs.
//
// The caller of Query.Next supplies a slice of
// Blots, indicating to the index/query implementation
// for how many blots we would like results.
//
// For each sub Blot, the field docs can either
// be nil, indicating to show all documents, or
// non-nil, in which case up to len(Docs) - cap(Docs)
// doc records are returned, each associated with
// Blot.
type Blot struct {
	Blot uint32
	Docs []Doc
}

func (b *Blot) Doc(i int) *Doc {
	return &b.Docs[i]
}

func (b *Blot) Cap() int {
	return cap(b.Docs)
}

func (b *Blot) Len() int {
	return len(b.Docs)
}

func (b *Blot) Next() *Doc {
	n := len(b.Docs)
	b.Docs = b.Docs[:n+1]
	return &b.Docs[n]
}
