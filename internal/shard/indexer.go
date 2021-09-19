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

package shard

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/go-air/dupi/post"
)

type Indexer struct {
	id      uint32
	root    string
	postChn chan []post.T
	ind     [1 << 16]poster

	postFile *os.File
}

func (x *Indexer) initCommon(id uint32, root string, flushRate uint32) {
	x.root = root
	x.id = id
	x.postChn = make(chan []post.T)
}

func (x *Indexer) InitCreate(id uint32, root string, flushRate uint32) error {
	x.initCommon(id, root, flushRate)
	// maybe this is unnecessary, just write on close...
	iixName := fmt.Sprintf("%s.iix", x.root)
	f, err := os.OpenFile(iixName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	for i := range x.ind {
		z := &x.ind[i]
		if err := z.initCreate(f, uint16(i)); err != nil {
			return err
		}
	}
	x.postFile, err = os.OpenFile(x.root, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	return err
}

func (x *Indexer) InitAppend(id uint32, root string, flushRate uint32) error {
	x.initCommon(id, root, flushRate)
	return x.read()
}

func (x *Indexer) PostChan() chan []post.T {
	return x.postChn
}

func (x *Indexer) read() error {
	var err error
	if err = x.readIix(); err != nil {
		return err
	}
	x.postFile, err = os.OpenFile(x.root, os.O_RDWR|os.O_SYNC, 0644)
	if err != nil {
		return err
	}
	for i := range x.ind {
		p := &x.ind[i]
		err = p.readToLink(x.postFile)
		if err != nil {
			return fmt.Errorf("error reading to link for blot %x: %w", i, err)
		}
	}
	return err
}

func (x *Indexer) readIix() error {
	iixName := fmt.Sprintf("%s.iix", x.root)
	f, err := os.OpenFile(iixName, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	br := bufio.NewReader(f)
	for i := range x.ind {
		z := &x.ind[i]
		err = z.initAppend(br, uint16(i))
		if err != nil {
			return err
		}
	}
	return nil
}

//
// Flush everything and close
//
// list of doc info
// delta, delta, delta fileid, start-lastend, end-start
// this assumes the postChn is closed.
func (x *Indexer) Close() error {
	if err := x.flushInd(); err != nil {
		return err
	}
	return x.flushIix()
}

func (x *Indexer) flushIix() error {
	iix, err := os.OpenFile(fmt.Sprintf("%s.iix", x.root),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer iix.Close()
	for i := range x.ind {
		up := &x.ind[i]
		_, err = up.writeVarint64(iix, up.head)
		if err != nil {
			return err
		}
		_, err = up.writeVarint64(iix, int64(up.total))
		if err != nil {
			return err
		}
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[:], up.current)
		_, err = iix.Write(buf[:])
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *Indexer) flushInd() error {
	var err error
	defer x.postFile.Close()
	for i := range x.ind {
		up := &x.ind[i]
		err = up.flushTo(x.postFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *Indexer) Serve() {
	for {
		ps, ok := <-x.postChn
		if !ok {
			x.postChn = nil
			return
		}
		for _, p := range ps {
			docid, hash := p.Docid(), p.Blot()
			hash &= 0xffff
			fmt.Printf("adding post (%d,%x)\n", docid, hash)
			err := x.ind[hash].AddPost(docid, x.postFile)
			if err != nil {
				log.Printf("couldn't add post: %s", err)
			}
		}
		x.postChn <- nil
	}
}

func (x *Indexer) String() string {
	return x.root
}
