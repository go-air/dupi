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
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-air/dupi"
)

type inspectCmd struct {
	verb
	json  *bool
	files *bool
}

func newInspectCmd() *inspectCmd {
	sub := &verb{
		name:  "inspect",
		flags: flag.NewFlagSet("inspect", flag.ExitOnError)}
	res := &inspectCmd{
		verb: *sub,
		json:   sub.flags.Bool("json", false, "output json."),
		files:  sub.flags.Bool("files", false, "output file info.")}
	return res
}

func (in *inspectCmd) Usage() string {
	return "inspect the root index."
}

func (in *inspectCmd) Run(args []string) error {
	var (
		err error
		idx *dupi.Index
	)
	in.flags.Parse(args)
	idx, err = dupi.OpenIndex(getIndexRoot())
	if err != nil {
		return err
	}
	defer idx.Close()
	st, err := idx.Stats()
	if err != nil {
		log.Fatal(err)
	}
	if *in.json {
		d, err := json.MarshalIndent(st, "", "\t")
		if err != nil {
			log.Fatal(err)
		}
		os.Stdout.Write(d)
	} else {
		fmt.Print(st)
	}
	return nil
}
