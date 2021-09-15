package main

import (
	"flag"
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
	tokfn := idx.TokenFunc()
	blotter := idx.Blotter()
	_ = tokfn
	_ = blotter
	return nil
}
