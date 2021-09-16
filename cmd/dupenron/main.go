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
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	_ "net/http/pprof"

	"github.com/go-air/dupi"
	"github.com/google/gops/agent"
)

var N = flag.Int("n", 16, "buckets")
var S = flag.Int("s", 10, "similarity seq length")
var O = flag.Bool("o", false, "o")
var Root = flag.String("r", "dupi.idx", "index dir")
var pprof = flag.String("pprof", "", "profile address")

func csvReader(fname string) (*csv.Reader, func() error, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening '%s': %s", fname, err)
	}
	rdr := csv.NewReader(f)
	rdr.LazyQuotes = true
	rdr.FieldsPerRecord = 2
	rdr.ReuseRecord = true
	return rdr, f.Close, nil
}

func main() {
	flag.Parse()
	if *O {
		oEnron(flag.Args()[0])
		return
	}
	if err := agent.Listen(agent.Options{
		ShutdownCleanup: true, // automatically closes on os.Interrupt
	}); err != nil {
		log.Fatal(err)
	}
	index, err := dupi.CreateIndexer(*Root, *N, *S)
	if err != nil {
		log.Fatal(err)
	}
	fname := flag.Args()[0]
	err = indexEnron(index, fname)
	if err != nil {
		log.Fatal(err)
	}

	//log.Fatal(qEnron(index, fname))
}

func oEnron(fname string) error {
	rdr, closer, err := csvReader(fname)
	if err != nil {
		return fmt.Errorf("error opening '%s': %s", fname, err)
	}
	defer closer()
	writer := csv.NewWriter(os.Stdout)
	writer.Comma = rdr.Comma
	for i := 0; i < *N; i++ {
		rcd, err := rdr.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error reading csv: %s", err)
		}
		err = writer.Write(rcd)
		if err != nil {
			return err
		}
	}
	return nil
}

func indexEnron(index *dupi.Indexer, fname string) error {
	rdr, closer, err := csvReader(fname)
	if err != nil {
		return fmt.Errorf("error opening '%s': %s", fname, err)
	}
	defer closer()
	rdr.ReuseRecord = true
	rcds := 0
	offset := uint32(0)
	for {
		rcd, err := rdr.Read()
		if err != nil {
			if err == io.EOF {
				return index.Close()
			}
			return fmt.Errorf("error reading csv: %s", err)
		}
		msg := []byte(rcd[1])
		ibody := bytes.Index(msg, []byte("\n\n"))
		if ibody != -1 {
			msg = msg[ibody+2:]
		}
		doc := &dupi.Doc{
			Path:  fname, //"/" + rcd[0],
			Dat:   msg,
			Start: offset + uint32(ibody) + 2 + uint32(len(rcd[0]))}
		doc.End = offset + uint32(len(doc.Dat))
		offset = doc.End
		index.Add(doc)
		rcds++
		if rcds%10000 == 0 {
			fmt.Printf("%d rcds\n", rcds)
		}
	}
}

func qEnron(index *dupi.Indexer, fname string) error {
	rdr, closer, err := csvReader(fname)
	if err != nil {
		return fmt.Errorf("error opening '%s': %s", fname, err)
	}
	defer closer()
	rcds := 0
	offset := uint32(0)
	for {
		rcd, err := rdr.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error reading csv: %s", err)
		}
		doc := &dupi.Doc{
			Path:  "/" + rcd[0],
			Dat:   []byte(rcd[1]),
			Start: offset}
		//index.Query(doc.Dat)
		_ = doc

		rcds++
		if rcds%10000 == 0 {
			log.Printf("%d rcds\n", rcds)
		}
	}
}
