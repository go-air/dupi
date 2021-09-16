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
	"sync"

	"github.com/go-air/dupi/blotter"
	"github.com/go-air/dupi/post"
	"github.com/go-air/dupi/token"
)

type shatterReq struct {
	docid    uint32
	offset   uint32
	d        []byte
	shutdown bool
}

func startShatter(ns, n, s int,
	tf token.TokenizerFunc, blotcfg *blotter.Config,
	chns []chan []post.T) (chan *shatterReq, error) {
	rch := make(chan *shatterReq)
	mono := newMono()
	for i := 0; i < ns; i++ {
		bler, err := blotter.FromConfig(blotcfg)
		if err != nil {
			return nil, err
		}
		sh := newShatter(n, s, tf, bler, mono)
		copy(sh.shardChns, chns)
		go func(sh *shatter) {
			for {
				req, ok := <-rch
				if !ok {
					return
				}
				if req.shutdown {
					return
				}
				sh.do(req.docid, req.offset, req.d)
			}
		}(sh)
	}
	return rch, nil
}

type mono struct {
	docid uint32
	cond  *sync.Cond
}

func newMono() *mono {
	var mu sync.Mutex
	return &mono{cond: sync.NewCond(&mu)}
}

type shatter struct {
	tokfn     token.TokenizerFunc
	tokb      []token.T
	bler      blotter.T
	seqlen    int
	d         [][]post.T
	shardChns []chan []post.T
	mono      *mono
}

func newShatter(n, s int, tf token.TokenizerFunc, bler blotter.T, mono *mono) *shatter {
	res := &shatter{
		tokfn:     tf,
		bler:      bler,
		seqlen:    s,
		shardChns: make([]chan []post.T, n),
		d:         make([][]post.T, n),
		mono:      mono}
	for i := range res.shardChns {
		res.shardChns[i] = make(chan []post.T)
	}
	return res
}

func (s *shatter) do(did, offset uint32, msg []byte) {
	s.tokb = s.tokfn(s.tokb[:0], msg, offset)
	var (
		words = 0
		b     = uint32(0)
	)
	for i := range s.tokb {
		tok := &s.tokb[i]
		switch tok.Tag {
		case token.Word:
			b = s.bler.Blot(tok.Lit)
			words++
			if words > s.seqlen {
				s.blot(did, b)
			}
		default:
		}
	}
	s.send(did)
}

func (s *shatter) send(did uint32) {
	s.mono.cond.L.Lock()
	for s.mono.docid != did-1 {
		s.mono.cond.Wait()
	}

	var wg sync.WaitGroup
	for i, ps := range s.d {
		wg.Add(1)
		go func(i int, ps []post.T) {
			defer wg.Done()
			s.shardChns[i] <- ps
			<-s.shardChns[i]
			s.d[i] = nil //ps[:0] (was racy)

		}(i, ps)
	}
	wg.Wait()
	s.mono.docid = did
	s.mono.cond.Broadcast()
	s.mono.cond.L.Unlock()
}

func (s *shatter) blot(docid, b uint32) {
	n := uint32(len(s.d))

	i := b % n
	s.d[i] = append(s.d[i], post.Make(docid, b/n))
}
