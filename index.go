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

	"github.com/go-air/dupi/dmd"
	"github.com/go-air/dupi/internal/shard"
	"github.com/go-air/dupi/lock"
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
	res.shards = make([]shard.Index, cfg.NumBuckets)
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
