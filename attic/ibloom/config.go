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

package ibloom

import "fmt"

type Config struct {
	Name string
	M, K int
}

func FromConfig(cfg *Config) (Filter, error) {
	switch cfg.Name {
	case "fnv32":
		return NewFnv32(cfg.M/8, cfg.K), nil
	default:
		return nil, fmt.Errorf("unknown filter name: '%s'", cfg.Name)
	}
}
