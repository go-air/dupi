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
	"io"
	"log"
	"math"
	"os"

	"github.com/go-air/dupi"
)

type extractCmd struct {
	subCmd
	index *dupi.Index
	json  *bool
	sigma *float64
}

func newExtractCmd() *extractCmd {
	var extract = &extractCmd{subCmd: subCmd{
		name:  "extract",
		usage: "extract [args]",
		flags: flag.NewFlagSet("extract", flag.ExitOnError)}}

	extract.json = extract.flags.Bool("json", false, "output json")
	extract.sigma = extract.flags.Float64("sigma", 2.0, "explore blots within σ of average (higher=most probable dups, lower=more volume)")
	return extract
}

func (x *extractCmd) Usage() string {
	return `extract from the index root`
}

func (x *extractCmd) Run(args []string) error {
	var err error
	x.flags.Parse(args)
	x.index, err = dupi.OpenIndex(getIndexRoot())
	if err != nil {
		return err
	}
	defer x.index.Close()
	st, err := x.index.Stats()
	if err != nil {
		log.Fatal(err)
	}
	σ := *x.sigma
	N := int(math.Round(st.BlotMean + σ*st.BlotSigma))
	query := x.index.StartQuery(dupi.QueryMaxBlot)
	shape := []dupi.Blot{{Blot: 0}}
	for {
		n, err := query.Next(shape)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if n == 0 {
			return fmt.Errorf("Query.Next gave 0 and no error")
		}
		if len(shape[0].Docs) < N {
			return nil
		}
		if *x.json {
			shp2 := shape
			j := 0
			for i := range shp2 {
				if len(shp2[i].Docs) <= 1 {
					continue
				}
				shp2[j] = shp2[i]
			}
			d, err := json.MarshalIndent(shp2, "", "\t")
			if err != nil {
				return err
			}
			_, err = os.Stdout.Write(d)
			if err != nil {
				return err
			}
		} else {
			for _, blot := range shape {
				if len(blot.Docs) <= 1 {
					continue
				}
				fmt.Printf("%x", blot.Blot)
				for i := range blot.Docs {
					doc := &blot.Docs[i]
					fmt.Printf(" %s@%d:%d", doc.Path, doc.Start, doc.End)
				}
				fmt.Printf("\n")
			}
		}
		for i := range shape {
			shape[i].Docs = nil
		}
	}
}
