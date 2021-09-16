package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/go-air/dupi"
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
			start, end, err := idx.FindBlot(hex, doc)
			if err != nil {
				log.Printf("warning: %s", err)
				continue
			}
			dat := string(doc.Dat[start:end])
			doc.Dat = nil
			m[dat] = append(m[dat], doc)
		}
		for k, ds := range m {
			fmt.Printf("text:\n'''\n%s'''\n", k)
			for _, d := range ds {
				fmt.Printf("\t%s %d:%d\n", d.Path, d.Start, d.End)
			}
		}
	}
	return nil
}
