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

	"github.com/go-air/dupi"
	"github.com/go-air/dupi/token"
)

type blotCmd struct {
	subCmd
	offsets *bool
}

func newBlotCmd() *blotCmd {
	cmd := &blotCmd{
		subCmd: subCmd{name: "blot", flags: flag.NewFlagSet("blot", flag.ExitOnError)}}
	cmd.offsets = cmd.flags.Bool("offsets", false, "show text position of blots")
	return cmd
}

func (b *blotCmd) Usage() string {
	return "blot [files]"
}

func (b *blotCmd) Run(args []string) error {
	b.flags.Parse(args)
	root := getIndexRoot()
	idx, err := dupi.OpenIndex(root)
	if err != nil {
		log.Fatalf("couldn't open dupi index at '%s': %s", root, err)
	}
	defer idx.Close()
	for _, fname := range b.flags.Args() {
		if err = b.doFilename(fname, idx); err != nil {
			log.Printf("error processing %s: %s", fname, err)
		}
	}
	return nil
}

func (bc *blotCmd) doFilename(fname string, idx *dupi.Index) error {
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
	if *bc.offsets {
		toks = idx.TokenFunc()(nil, doc.Dat, 0)
	}
	N := uint32(idx.NumShards()) * (1 << 16)
	seqLen := idx.SeqLen()
	for i, b := range blots {
		if *bc.offsets {
			beg := toks[i].Pos
			end := toks[i+seqLen].Pos + uint32(len(toks[i+seqLen].Lit))
			fmt.Printf("%x %d:%d\n", b%N, beg, end)
		} else {
			fmt.Printf("%x\n", b%N)
		}
	}
	return nil
}
