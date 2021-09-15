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

// Package dupi provides a library for exploring duplicate
// data in large sets of documents.
package dupi

type Doc struct {
	Path  string
	Start uint32
	End   uint32
	Dat   []byte `json:"-"`
}

func NewDoc(path, body string) *Doc {
	return &Doc{
		Path: path,
		Dat:  []byte(body)}
}
