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

package token

import "fmt"

// Config describes tokenizer configurations.
type Config struct {
	Name string
	// nothing here yet.
}

// TokenizerFunc is the type of a function used
// for tokenizing document data.
type TokenizerFunc func(dst []T, dat []byte, offset uint32) []T

// DefaultConfig returns the default tokenizer config
// for dupy.
func DefaultConfig() *Config {
	return &Config{Name: "words.simple"}
}

// FromConfig attempts to create a tokenizer function
// from a configuration.
func FromConfig(cfg *Config) (TokenizerFunc, error) {
	switch cfg.Name {
	case "words.simple":
		return Tokenize, nil
	case "go.tokens":
		return GoTokenize, nil
	default:
		return nil, fmt.Errorf("unrecognized token config: '%#v'", cfg)
	}
}
