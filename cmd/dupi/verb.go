// Copyright 2018 The Reach Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

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
