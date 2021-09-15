// Copyright 2018 The Reach Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package main

import "flag"

type SubCmd interface {
	Name() string
	Flags() *flag.FlagSet
	Run(args []string) error
	Usage() string
}

type subCmd struct {
	name  string
	flags *flag.FlagSet
	usage string
}

func (s *subCmd) Name() string {
	return s.name
}

func (s *subCmd) Flags() *flag.FlagSet {
	return s.flags
}
