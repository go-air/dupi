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
	"log"
	"os"

	"github.com/scott-cotton/dupi/dmd"
	"github.com/scott-cotton/dupi/internal/shard"
	"github.com/scott-cotton/dupi/lock"
	"github.com/scott-cotton/dupi/post"
	"github.com/scott-cotton/dupi/token"
)

// Indexer is a struct for duplicate indexing.
type Indexer struct {
	config *Config
	lock   *lock.File

	didoff uint32
	dmds   *dmd.Adder

	// shatter *shatter
	shatter chan *shatterReq
	shards  []shard.Indexer
	fnames  *fnames
}

func IndexerFromConfig(cfg *Config) (*Indexer, error) {
	var err error
	res := &Indexer{config: cfg}
	res.lock, err = lock.New(cfg.LockPath())
	if err != nil {
		return nil, err
	}
	err = res.lock.Lock()
	if err != nil {
		return nil, err
	}
	res.shards = make([]shard.Indexer, cfg.NumBuckets)
	res.fnames = newFnames()
	if err := os.Mkdir(res.Root(), 0755); err != nil {
		return nil, err
	}
	tokfn, err := token.FromConfig(&cfg.TokenConfig)
	if err != nil {
		return nil, err
	}
	res.dmds = dmd.NewAdder(cfg.IndexRoot, cfg.DocFlushRate)
	postChans := make([]chan []post.T, len(res.shards))
	for i := range res.shards {
		shard := &res.shards[i]
		broot := cfg.PostPath(i)
		if err := shard.InitCreate(uint32(i), broot, uint32(cfg.DocFlushRate)); err != nil {
			return nil, err
		}
		postChans[i] = shard.PostChan()
		go shard.Serve()
	}
	res.shatter, err = startShatter(cfg.NumShatters,
		len(res.shards), cfg.SeqLen, tokfn, &cfg.BlotConfig, postChans)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func OpenIndexer(root string) (*Indexer, error) {
	cfg, err := ReadConfig(root)
	if err != nil {
		return nil, err
	}
	idx, err := openFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

// CreateIndexer attempts to creat a new dupy index.
// root is the directory root of the dupy index
// nbuckets states how many buckets
// docCap should be a conservative estimate of
// number of documents
// toksPerDoc should indicate about how many tokens
// are expected per document.
func CreateIndexer(root string, nbuckets, seqLen int) (*Indexer, error) {
	cfg, err := NewConfig(root, nbuckets, seqLen)
	if err != nil {
		return nil, err
	}
	return IndexerFromConfig(cfg)
}

// opens an index from a config.  cfg
// actually is read from the index as a first
// step, and this completes opening.
func openFromConfig(cfg *Config) (*Indexer, error) {
	var err error
	res := &Indexer{config: cfg}
	// lock file
	res.lock, err = lock.New(cfg.LockPath())
	if err != nil {
		return nil, err
	}
	err = res.lock.LockShared()
	if err != nil {
		return nil, err
	}
	// internal setup
	res.shards = make([]shard.Indexer, cfg.NumBuckets)
	res.fnames = newFnames()
	res.dmds = dmd.NewAdder(cfg.IndexRoot, cfg.DocFlushRate)

	if err = res.readfiles(); err != nil {
		return nil, err
	}
	tokenfn, err := token.FromConfig(&cfg.TokenConfig)
	if err != nil {
		return nil, err
	}
	blotcfg := &cfg.BlotConfig

	postChans := make([]chan []post.T, len(res.shards))
	for i := range res.shards {
		shard := &res.shards[i]
		broot := cfg.PostPath(i)
		if err := shard.InitAppend(uint32(i), broot, uint32(cfg.DocFlushRate)); err != nil {
			return nil, err
		}
		postChans[i] = shard.PostChan()
		go shard.Serve()
	}
	res.shatter, err = startShatter(cfg.NumShatters,
		len(res.shards), cfg.SeqLen, tokenfn, blotcfg, postChans)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// write the files that apeared in added documents.
func (x *Indexer) writeFiles() error {
	docPath := x.config.FnamesPath()
	f, e := os.OpenFile(docPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if e != nil {
		return e
	}
	defer f.Close()
	return x.fnames.write(f)
}

// Root returns the path to the root of the index 'x'.
// the returned root is an absolute path.
func (x *Indexer) Root() string {
	return x.config.IndexRoot
}

// readfiles reads the list of files associated
// with added documents.
func (x *Indexer) readfiles() error {
	docPath := x.config.FnamesPath()
	f, e := os.OpenFile(docPath, os.O_RDONLY, 0644)
	if e != nil {
		return e
	}
	defer f.Close()
	//r := bufio.NewReader(f)
	fnames, err := readFnames(f)
	x.fnames = fnames
	return err
}

// Close attempts to flush all data associated
// with the index to disk.
func (x *Indexer) Close() error {
	defer x.lock.Close()
	for i := 0; i < x.config.NumShatters; i++ {
		x.shatter <- &shatterReq{shutdown: true}
		// no more shatters running, each waits
		// for all shards to complete before
		// returning => shards are no longer busy.
	}
	for i := 0; i < x.config.NumBuckets; i++ {
		close(x.shards[i].PostChan())
	}
	if err := x.config.Write(); err != nil {
		return err
	}
	if err := x.writeFiles(); err != nil {
		return err
	}
	if err := x.dmds.Close(); err != nil {
		return err
	}
	errs := make(chan error, len(x.shards))
	for i := range x.shards {
		b := &x.shards[i]
		go func(b *shard.Indexer) {
			errs <- b.Close()
		}(b)
	}
	var err error
	for i := range x.shards {
		ierr := <-errs
		if err != nil && ierr != nil {
			log.Printf("dupy.Index.Close: dropping error %s from bucket %d", ierr, i)
		} else if ierr != nil {
			err = ierr
		}
	}
	return err
}

// Add adds 'doc' to the index.
func (x *Indexer) Add(doc *Doc) error {
	did, err := x.doc2Id(doc)
	if err != nil {
		return err
	}
	doc.Path = ""
	x.shatter <- &shatterReq{docid: did, offset: doc.Start, d: doc.Dat}
	return nil
}

func (x *Indexer) doc2Id(doc *Doc) (uint32, error) {
	n, err := x.fnames.addPath(doc.Path)
	if err != nil {
		return 0, err
	}
	doc.Path = ""
	return x.dmds.Add(n, doc.Start, doc.End)
}
