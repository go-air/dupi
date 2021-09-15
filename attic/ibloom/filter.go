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

// Package ibloom implements a bloom filter on integer (uint32) sets.
package ibloom

type Filter interface {
	K() int
	M() int
	Config() *Config
	Ones(v uint32, fn func(int) bool) bool
	PutOnes(dst []uint32, v uint32)
	Add(d []byte, v uint32)
	Has(d []byte, v uint32) bool
	Set() []byte
}
