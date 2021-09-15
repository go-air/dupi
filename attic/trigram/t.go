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

// Package trigram supports a trigram alphabet for dupy.
package trigram

import (
	"fmt"

	"github.com/go-air/dupi/token"
)

type T uint16

type C byte

const (
	NumLetters    = 37 // 0-9, a-z, _
	NumTrigrams   = NumLetters * NumLetters * NumLetters
	COther      C = 36
)

func Unletter(c byte) rune {
	if c <= 9 {
		return rune('0' + c)
	}
	if c <= 35 {
		return rune('a' + c - 10)
	}
	return '_'
}

func Letter(r rune) C {
	if r >= '0' && r <= '9' {
		return C(r - '0')
	}
	if r >= 'a' && r <= 'z' {
		return C(r - 'a' + 10)
	}
	if r >= 'A' && r <= 'Z' {
		return C(r - 'A' + 10)
	}
	return COther
}

func (t T) String() string {
	var buf [3]rune
	an := T(NumLetters)
	buf[2] = Unletter(byte(t % an))
	t /= an
	buf[1] = Unletter(byte(t % an))
	t /= an
	buf[0] = Unletter(byte(t % an))
	return string(buf[:])
}

func (t T) Shift(c C) T {
	an := T(NumLetters)
	t = t % (an * an)
	t *= an
	t += T(c)
	return t
}

func Zero() T {
	var t T
	return t.Shift(COther).Shift(COther).Shift(COther)
}

func VisitToken(start T, tok token.T, fn func(T)) T {
	var tri T
	switch tok.Tag {
	case token.Other:
		tri = start.Shift(COther)
		fn(tri)
	case token.Word:
		tri = start
		for _, r := range tok.Lit {
			tri = tri.Shift(Letter(rune(r)))
			fn(tri)
		}
	case token.Eod:
	default:
		panic(fmt.Sprintf("vt: %#v", tok))
	}
	return tri
}

func VisitTokens(toks []token.T, fn func(T)) {
	tri := Zero()
	for _, tok := range toks {
		tri = VisitToken(tri, tok, fn)
	}
}

func FromTokens(dst []T, toks []token.T) []T {
	tri := Zero()
	for _, tok := range toks {
		tri, dst = FromToken(dst, tri, tok)
	}
	return dst[2:]
}

func FromToken(dst []T, start T, tok token.T) (T, []T) {
	var tri T
	switch tok.Tag {
	case token.Other:
		tri = start.Shift(COther)
		dst = append(dst, tri)
	case token.Word:
		tri = start
		for _, r := range tok.Lit {
			tri = tri.Shift(Letter(rune(r)))
			dst = append(dst, tri)
		}
	default:
		panic("unexpected token tag")
	}
	return tri, dst
}
