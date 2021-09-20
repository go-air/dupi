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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/go-air/dupi"
	"github.com/go-air/dupi/token"
)

type likeCmd struct {
	verb
}

func newLikeCmd() *likeCmd {
	sc := &verb{flags: flag.NewFlagSet("like", flag.ExitOnError)}
	lc := &likeCmd{verb: *sc}
	return lc
}

func (lc *likeCmd) Usage() string {
	return "[file path]"
}

func (lc *likeCmd) Run(args []string) error {
	lc.flags.Parse(args)
	root := getIndexRoot()
	idx, err := dupi.OpenIndex(root)
	if err != nil {
		log.Fatalf("couldn't open dupi index at '%s': %s", root, err)
	}
	defer idx.Close()
	for _, fname := range args {
		if err := doFilename(idx, fname); err != nil {
			return err
		}
	}
	return nil
}

type docKey struct {
	Path       string
	start, end uint32
}

func dkey(doc *dupi.Doc) docKey {
	return docKey{
		Path:  doc.Path,
		start: doc.Start,
		end:   doc.End}
}

func doFilename(idx *dupi.Index, fname string) error {
	f, e := os.Open(fname)
	if e != nil {
		return e
	}
	defer f.Close()
	dat, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	doc := &dupi.Doc{Path: fname, Dat: dat, End: uint32(len(dat))}
	blots := idx.BlotDoc(nil, doc)
	var toks []token.T
	toks = idx.TokenFunc()(nil, doc.Dat, 0)
	j := 0
	for i := range toks {
		if toks[i].Tag != token.Word {
			continue
		}
		toks[j] = toks[i]
		j++
	}
	N := uint32(idx.NumShards()) * (1 << 16)
	seqLen := idx.SeqLen()
	bm := make(map[uint32][]byte, len(blots))
	for i, b := range blots {
		beg := toks[i].Pos
		end := toks[i+seqLen].Pos + uint32(len(toks[i+seqLen].Lit))
		bm[b%N] = dat[beg:end]
	}
	query := idx.StartQuery(dupi.QueryMaxBlot)
	var db dupi.Blot
	found := make(map[docKey]int)
	for _, blot := range blots {
		b := blot % N
		db.Blot = b
		db.Docs = nil
		if err := query.Get(&db); err != nil {
			return err
		}
		for i := range db.Docs {
			doc := &db.Docs[i]
			dk := dkey(doc)
			if found[dk] != 0 {
				continue
			}
			rm, err := idx.FindBlots(bm, doc)
			if err != nil {
				return err
			}
			if len(rm) == 0 {
				continue
			}
			found[dk]++
		}
	}
	keys := make([]docKey, 0, len(found))
	for k, _ := range found {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return found[keys[i]] < found[keys[j]]
	})
	fmt.Printf("like %s:\n", fname)
	for _, dk := range keys {
		fmt.Printf("\t%s %d:%d\n", dk.Path, dk.start, dk.end)
	}
	return nil
}
