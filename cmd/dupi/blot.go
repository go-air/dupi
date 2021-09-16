package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-air/dupi"
	"github.com/go-air/dupi/blotter"
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
	tokfn := idx.TokenFunc()
	blotter := idx.Blotter()
	N := uint32(idx.NumShards() * (1 << 16))

	for _, fname := range b.flags.Args() {
		if err = b.doFilename(fname, idx, tokfn, blotter, N); err != nil {
			log.Printf("error processing %s: %s", fname, err)
		}
	}
	return nil
}

func (b *blotCmd) doFilename(fname string, idx *dupi.Index, tokfn token.TokenizerFunc, blotter blotter.T, N uint32) error {
	f, e := os.Open(fname)
	if e != nil {
		return e
	}
	defer f.Close()
	dat, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	toks := tokfn(nil, dat, 0)
	for i, tok := range toks {
		switch tok.Tag {
		case token.Word:
			blot := blotter.Blot(tok.Lit)
			if i < idx.SeqLen() {
				continue
			}
			if *b.offsets {
				end := tok.Pos + uint32(len(tok.Lit))
				beg := toks[i-idx.SeqLen()].Pos
				fmt.Printf("%x %d:%d\n", blot%N, beg, end)
			} else {
				fmt.Printf("%x\n", blot%N)
			}
		default:
		}
	}
	return nil
}
