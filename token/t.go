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

// Package token tokenizes data for dupi.
package token

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// Tag represents a value in an enumeration of
// values to associate with a token.
type Tag int

const (
	Word Tag = 35 + iota
	Other
	Eod
)

func (t Tag) String() string {
	switch t {
	case Word:
		return "w"
	case Other:
		return "_"
	case Eod:
		return "$"
	default:
		panic("wow")
	}
}

// Type T represents a token.
type T struct {
	Tag Tag
	Lit []byte
	Pos uint32
}

func (t *T) String() string {
	return fmt.Sprintf("<token %s: '%s' @%d>", t.Tag, t.Lit, t.Pos)
}

// Tokenize is a tokenizer function.
func Tokenize(dst []T, d []byte, offset uint32) []T {
	if !utf8.Valid(d) {
		return dst
	}
	inWord := false
	var i, j int
	var r rune
	for i, r = range string(d) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if !inWord && j < i {
				dst = append(dst,
					T{
						Lit: d[j:i],
						Pos: offset + uint32(j),
						Tag: Other})
				j = i
			}
			inWord = true
			continue
		}
		if inWord && j < i {
			dst = append(dst,
				T{
					Lit: d[j:i],
					Pos: offset + uint32(j),
					Tag: Word})
			j = i
		}
		inWord = false
	}
	if j < i {
		if inWord {
			dst = append(dst,
				T{
					Lit: d[j:i],
					Pos: offset + uint32(j),
					Tag: Word})
		} else {
			dst = append(dst,
				T{
					Lit: d[j:i],
					Pos: offset + uint32(j),
					Tag: Other})
		}
	}
	return dst
}
