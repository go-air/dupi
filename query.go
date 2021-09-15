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

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/go-air/dupi/internal/shard"
)

type qstate struct {
	shardStates []*shard.ReadState
	i           uint32
	n           uint32
	blot        uint32
	nilCount    uint32
}

var ErrInvalidQueryState = errors.New("query state invalid")

type QueryStrategy int

const (
	QueryMaxBlot QueryStrategy = iota
	QueryMaxDoc
	QueryRandom
)

type Query struct {
	index    *Index
	state    *qstate
	strategy QueryStrategy
}

func (q *Query) Next(dst []Blot) (n int, err error) {
	q.state.blot = ^uint32(0)
	state := q.state
	for n < len(dst) {
		dstBlot := &dst[n]
		shardState := state.shardStates[state.i]
		if shardState == nil {
			state.nilCount++
			if state.nilCount == state.n {
				if n == 0 {
					err = io.EOF
					return
				}
				return
			}
			state.i++
			if state.i == state.n {
				state.i = 0
			}
			continue
		}
		_, err = q.fillBlot(dstBlot, shardState, state.i)
		if err != nil {
			return
		}
		if len(dstBlot.Docs) <= 1 {
			q.advance(shardState, state.i)
			dstBlot.Docs = dstBlot.Docs[:0]
			continue
		}
		n++
	}
	return n, nil
}

func (q *Query) fillBlot(dst *Blot, src *shard.ReadState, srcPos uint32) (int, error) {
	var (
		docid uint32
		err   error
		n     int
	)
	q.state.blot = uint32(src.Blot)
	dst.Blot = uint32(src.Blot)*q.state.n + q.state.i
	for dst.Len() < dst.Cap() {
		docid, err = src.Next()
		if err == io.EOF {
			q.advance(src, srcPos)
			return n, nil
		} else if err != nil {
			return 0, err
		} else if docid == 0 {
			continue
		}
		n++
		q.index.docid2Doc(docid, dst.Next())
	}
	return n, err
}

func (q *Query) advance(src *shard.ReadState, pos uint32) *shard.ReadState {
	var rs *shard.ReadState
	if src.At == math.MaxUint16 {

	} else if src.Total <= 1 {

	} else {
		rs = q.index.shards[pos].ReadStateAt(src.At + 1)
	}
	q.state.shardStates[pos] = rs
	return rs
}
