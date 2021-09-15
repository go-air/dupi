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

package dupi

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/scott-cotton/dupi/blotter"
	"github.com/scott-cotton/dupi/token"
)

type Config struct {
	IndexRoot   string
	SeqLen      int
	NumBuckets  int
	NumShatters int

	// How frequently buckets write document
	// data to disk.  Higher= less memory,
	// more frequent i/o.
	// Frequency in terms of number of documents.
	DocFlushRate int

	TokenConfig token.Config
	BlotConfig  blotter.Config
}

func DefaultConfig(root string) (*Config, error) {
	cfg := &Config{}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	cfg.IndexRoot = abs
	cfg.DocFlushRate = 16384
	cfg.NumBuckets = 2
	cfg.NumShatters = 2
	cfg.SeqLen = 10
	cfg.TokenConfig = *token.DefaultConfig()
	cfg.BlotConfig = *blotter.DefaultConfig()
	return cfg, nil
}

func NewConfig(root string, nbuckets, seqLen int) (*Config, error) {
	cfg, err := DefaultConfig(root)
	if err != nil {
		return nil, err
	}
	cfg.NumBuckets = nbuckets
	cfg.SeqLen = seqLen
	return cfg, nil
}

func (cfg *Config) Path() string {
	return filepath.Join(cfg.IndexRoot, "cfg.json")
}

func (cfg *Config) LockPath() string {
	return cfg.IndexRoot + ".lock"
}

func (cfg *Config) DmdPath() string {
	return filepath.Join(cfg.IndexRoot, "dmd")
}

func (cfg *Config) FnamesPath() string {
	return filepath.Join(cfg.IndexRoot, "files.fnm")
}

func (cfg *Config) IixPath(i int) string {
	return filepath.Join(cfg.IndexRoot, fmt.Sprintf("b%d.iix", i))
}

func (cfg *Config) PostPath(i int) string {
	return filepath.Join(cfg.IndexRoot, fmt.Sprintf("b%d.pos", i))
}

func ReadConfig(root string) (*Config, error) {
	cfg, err := DefaultConfig(root)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(cfg.Path(), os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// TBD: check for huge file
	d, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(d, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *Config) Write() error {
	f, err := os.OpenFile(cfg.Path(), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	d, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}
	_, err = f.Write(d)
	if err != nil {
		return err
	}
	d[0] = '\n'
	_, err = f.Write(d[:1])
	if err != nil {
		return err
	}
	return nil
}
