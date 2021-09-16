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

package main

import (
	"flag"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-air/dupi"
)

type indexCmd struct {
	subCmd
	shards  *int
	seqlen  *int
	add     *bool
	verbose *bool
	nshat   *int
	indexer *dupi.Indexer
}

var index = &indexCmd{
	subCmd: subCmd{
		name:  "index",
		flags: flag.NewFlagSet("index", flag.ExitOnError)}}

func init() {
	index.shards = index.flags.Int("n", 4, "shards")
	index.seqlen = index.flags.Int("t", 10, "similarity based seq len")
	index.add = index.flags.Bool("a", false, "add to a given existing index")
	index.verbose = index.flags.Bool("v", false, "verbose")
	index.nshat = index.flags.Int("s", 4, "num shatterers")
}

var _ SubCmd = index

func (x *indexCmd) Usage() string {
	return "paths"
}

func (x *indexCmd) getIndexer() (*dupi.Indexer, error) {
	if *x.add {
		return dupi.OpenIndexer(getIndexRoot())
	}
	cfg, err := dupi.NewConfig(getIndexRoot(), *x.shards, *x.seqlen)
	cfg.NumShatters = *x.nshat
	if err != nil {
		return nil, err
	}
	return dupi.IndexerFromConfig(cfg)
}

func (x *indexCmd) Run(args []string) error {
	x.flags.Parse(args)
	idx, err := x.getIndexer()
	if err != nil {
		return err
	}
	defer func() {
		err := idx.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	x.indexer = idx
	var reterr error
	for _, path := range x.flags.Args() {
		err := x.doPath(path)
		if err != nil {
			log.Printf("error processing '%s': %s", path, err)
			reterr = err
		}
	}
	return reterr
}

func (x *indexCmd) mkWalkFn(perr *error) func(path string, entry fs.DirEntry, err error) error {
	return func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("error %s", err)
			*perr = err
			return fs.SkipDir
		}
		if strings.HasPrefix(entry.Name(), ".") && len(entry.Name()) > 1 {
			if entry.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			return nil
		}

		if *x.verbose {
			log.Printf("doPath %s\n", path)
		}

		f, e := os.Open(path)
		if e != nil {
			return e
		}
		defer f.Close()
		dat, err := ioutil.ReadAll(f)
		if err != nil {
			*perr = err
			return err
		}
		doc := &dupi.Doc{Path: path, Dat: dat, End: uint32(len(dat))}
		return x.indexer.Add(doc)
	}
}

func (x *indexCmd) doPath(fpath string) error {
	var err error
	filepath.WalkDir(fpath, x.mkWalkFn(&err))
	return err
}
