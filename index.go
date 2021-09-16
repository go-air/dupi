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
	"fmt"
	"log"
	"os"

	"github.com/go-air/dupi/blotter"
	"github.com/go-air/dupi/dmd"
	"github.com/go-air/dupi/internal/shard"
	"github.com/go-air/dupi/lock"
	"github.com/go-air/dupi/token"
)

type Index struct {
	config *Config
	lock   *lock.File
	dmd    *dmd.T
	fnames *fnames
	shards []shard.Index
}

func OpenIndex(root string) (*Index, error) {
	var err error
	cfg := &Config{IndexRoot: root}
	res := &Index{config: cfg}
	res.lock, err = lock.New(cfg.LockPath())
	if err != nil {
		return nil, err
	}
	if err := res.lock.LockShared(); err != nil {
		return nil, err
	}
	cfg, err = ReadConfig(root)
	if err != nil {
		return nil, err
	}
	res.config = cfg
	res.dmd, err = dmd.New(cfg.IndexRoot)
	if err != nil {
		return nil, err
	}
	fnf, err := os.Open(cfg.FnamesPath())
	if err != nil {
		return nil, err
	}
	defer fnf.Close()
	res.fnames, err = readFnames(fnf)
	if err != nil {
		return nil, err
	}
	res.shards = make([]shard.Index, cfg.NumShards)
	for i := range res.shards {
		shard := &res.shards[i]
		if err := shard.Init(cfg.PostPath(i)); err != nil {
			return nil, fmt.Errorf("error initializing shard %d: %w", i, err)
		}
	}
	return res, nil
}

func (x *Index) Close() error {
	var err error
	defer func() { err = x.lock.Close() }()
	err = x.dmd.Close()
	for i := range x.shards {
		s := &x.shards[i]
		serr := s.Close()
		if serr != nil {
			if err != nil {
				log.Printf("dropping close error: %s shard %d", serr, i)
			} else {
				err = serr
			}
		}
	}
	return err
}

func (x *Index) Root() string {
	return x.config.IndexRoot
}

func (x *Index) TokenFunc() token.TokenizerFunc {
	tf, err := token.FromConfig(&x.config.TokenConfig)
	if err != nil {
		panic(err) // should be impossible.
	}
	return tf
}

func (x *Index) Blotter() blotter.T {
	sh, err := blotter.FromConfig(&x.config.BlotConfig)
	if err != nil {
		panic(err) // should be impossible.
	}
	return sh
}

func (x *Index) NumShatters() int {
	return x.config.NumShatters
}

func (x *Index) NumShards() int {
	return x.config.NumShards
}

func (x *Index) SeqLen() int {
	return x.config.SeqLen
}

func (x *Index) BlotDoc(dst []uint32, doc *Doc) []uint32 {
	tokfn := x.TokenFunc()
	blotter := x.Blotter()
	toks := tokfn(nil, doc.Dat, doc.Start)
	j := 0
	for _, tok := range toks {
		if tok.Tag != token.Word {
			continue
		}
		toks[j] = tok
		j++
	}
	seqLen := x.SeqLen()
	for i, tok := range toks[:j] {
		blot := blotter.Blot(tok.Lit)
		if i < seqLen {
			continue
		}
		dst = append(dst, blot)
	}
	return dst
}

func (x *Index) SplitBlot(b uint32) (shard uint32, sblot uint16) {
	nsh := uint32(x.NumShards())
	b = b % (nsh * (1 << 16))
	shard = b % nsh
	sblot = uint16(b / nsh)
	return
}

func (x *Index) JoinBlot(shard uint32, sblot uint16) uint32 {
	nsh := uint32(x.NumShards())
	blot := nsh * uint32(sblot)
	blot += shard
	return blot

}

func (x *Index) FindBlot(theBlot uint32, doc *Doc) (start, end uint32, err error) {
	if doc.Dat == nil {
		err = doc.Load()
		if err != nil {
			return
		}
	}
	toks := x.TokenFunc()(nil, doc.Dat, doc.Start)
	j := 0
	for _, tok := range toks {
		if tok.Tag != token.Word {
			continue
		}
		toks[j] = tok
		j++
	}
	blotter := x.Blotter()
	seqLen := x.SeqLen()
	nShard := uint32(x.NumShards())
	for i, tok := range toks[:j] {
		blot := blotter.Blot(tok.Lit)
		if i < seqLen {
			continue
		}
		blot %= nShard * (1 << 16)
		if blot != theBlot {
			continue
		}
		start = toks[i-seqLen].Pos
		end = tok.Pos + uint32(len(tok.Lit))
		return
	}
	err = fmt.Errorf("blot %x not found", theBlot)
	return
}

func (x *Index) StartQuery(s QueryStrategy) *Query {
	q := &Query{
		index:    x,
		strategy: QueryMaxBlot,
		state:    x.qstate(s)}
	return q
}

func (x *Index) qstate(s QueryStrategy) *qstate {
	qstate := &qstate{}
	qstate.i = uint32(0)
	qstate.n = uint32(len(x.shards))
	qstate.shardStates = make([]*shard.ReadState, qstate.n)
	for i := range x.shards {
		shard := &x.shards[i]
		qstate.shardStates[i] = shard.ReadStateAt(0)
	}
	qstate.blot = uint32(qstate.shardStates[qstate.i].Blot)
	qstate.blot *= uint32(qstate.n)
	qstate.blot += qstate.i
	return qstate
}

func (x *Index) docid2Doc(did uint32, doc *Doc) error {
	fid, start, end, err := x.dmd.Lookup(did)
	if err != nil {
		return err
	}
	doc.Path = x.fnames.abs(fid)
	doc.Start = start
	doc.End = end
	return nil
}
