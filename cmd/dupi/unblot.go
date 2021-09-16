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

type unblotCmd struct {
	subCmd
}

func newUnblotCmd() *unblotCmd {
	return &unblotCmd{
		subCmd: subCmd{name: "unblot", flags: flag.NewFlagSet("unblot", flag.ExitOnError)}}
}

func (ub *unblotCmd) Usage() string {
	return "unblot <blot>"
}

func (ub *unblotCmd) Run(args []string) error {
	ub.flags.Parse(args)
	root := getIndexRoot()
	idx, err := dupi.OpenIndex(root)
	if err != nil {
		log.Fatalf("couldn't open dupi index at '%s': %s", root, err)
	}
	defer idx.Close()
	tokfn := idx.TokenFunc()
	blotter := idx.Blotter()
	query := idx.StartQuery(dupi.QueryMaxBlot)
	for _, arg := range ub.flags.Args() {
		var hex uint32
		if _, err := fmt.Sscanf(arg, "%x", &hex); err != nil {
			return err
		}
		blot := &dupi.Blot{Blot: hex}
		if err := query.Get(blot); err != nil {
			return err
		}

		m := make(map[string][]*dupi.Doc)
		for i := range blot.Docs {
			doc := &blot.Docs[i]
			dat, err := findBlot(doc, tokfn, blotter, blot.Blot, uint32(idx.NumShards()), idx.SeqLen())
			if err != nil {
				log.Printf("warning: %s", err)
				continue
			}
			m[string(dat)] = append(m[string(dat)], doc)
		}
		for k, ds := range m {
			for _, d := range ds {
				fmt.Printf("\t%s %d:%d\n", d.Path, d.Start, d.End)
			}
		}
	}
	_ = tokfn
	_ = blotter
	return nil
}

func findBlot(doc *dupi.Doc, tokfn token.TokenizerFunc, blotter blotter.T, theBlot uint32, nShard uint32, seqLen int) ([]byte, error) {
	if doc.Dat == nil {
		f, err := os.Open(doc.Path)
		if err != nil {
			return nil, err
		}
		doc.Dat, err = ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		f.Close()
		doc.Dat = doc.Dat[doc.Start:doc.End]
	}
	toks := tokfn(nil, doc.Dat, doc.Start)
	for i, tok := range toks {
		if tok.Tag != token.Word {
			continue
		}
		blot := blotter.Blot(tok.Lit)
		if i < seqLen {
			continue
		}
		blot %= nShard * (1 << 16)
		if blot != theBlot {
			continue
		}
		start := toks[i-seqLen].Pos
		end := tok.Pos + uint32(len(tok.Lit))
		return doc.Dat[start:end], nil
	}
	return nil, nil
}
