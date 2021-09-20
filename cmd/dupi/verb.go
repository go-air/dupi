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

import "flag"

type Verb interface {
	Name() string
	Flags() *flag.FlagSet
	Run(args []string) error
	Usage() string
}

type verb struct {
	name  string
	flags *flag.FlagSet
	usage string
}

func (s *verb) Name() string {
	return s.name
}

func (s *verb) Flags() *flag.FlagSet {
	return s.flags
}
