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
}

func newBlotCmd() *blotCmd {
	return &blotCmd{
		subCmd: subCmd{name: "blot", flags: flag.NewFlagSet("blot", flag.ExitOnError)}}
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
		fmt.Printf("%s\n", fname)
		f, e := os.Open(fname)
		if e != nil {
			log.Printf("error opening '%s': %s", fname, e)
			continue
		}
		func(f *os.File) {
			defer f.Close()
			dat, err := ioutil.ReadAll(f)
			if err != nil {
				log.Printf("error reading '%s': %s", fname, e)
				return
			}
			toks := tokfn(nil, dat, 0)
			for i, tok := range toks {
				switch tok.Tag {
				case token.Word:
					blot := blotter.Blot(tok.Lit)
					if i >= idx.SeqLen() {
						end := tok.Pos + uint32(len(tok.Lit))
						beg := toks[i-idx.SeqLen()].Pos
						fmt.Printf("'%s' %x\n", string(dat[beg:end]), blot%N)
					}
				default:
				}
			}
		}(f)
	}
	return nil
}
