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

// Command dupi is the dupi command line.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/gops/agent"
)

// TBD: write and use newIndex(), newExtract()
// instead of init() ugliness
var scMap = map[string]Verb{
	"index":   newIndexCmd(),
	"extract": newExtractCmd(),
	"blot":    newBlotCmd(),
	"unblot":  newUnblotCmd(),
	"inspect": newInspectCmd(),
	"like":    newLikeCmd()}

var gFlags = flag.NewFlagSet("dupi", flag.ExitOnError)

var root = gFlags.String("r", "", "index root")

func getIndexRoot() string {
	if *root != "" {
		return *root
	}
	r := os.Getenv("DUPIROOT")
	if r != "" {
		return r
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("couldn't get home dir: %s, using 'dupi'\n", err)
		return "dupi"
	}
	return filepath.Join(home, ".dupi")
}

// returns global argument list
func splitArgs(args []string) ([]string, []string) {
	var i int
	var arg string
	var sc Verb
	for i, arg = range args {
		sc = scMap[arg]
		if sc != nil {
			break
		}
	}
	return args[:i], args[i:]
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "usage: dupi [global opts] <verb> <args>\n")
	fmt.Fprintf(w, "verbs are:\n")
	for k, v := range scMap {
		fmt.Fprintf(w, "\t%-10s %30s\n", k, v.Usage())
	}
	fmt.Fprintf(w, "\nglobal options:\n")
	gFlags.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(w, "\t-%-16s default=%-16q %s\n", f.Name, f.DefValue, f.Usage)
	})
	fmt.Fprintf(w, "\nTo get help on a verb, try dupi <verb> -h.\n")
}

func usageFatal(w io.Writer) {
	usage(w)
	os.Exit(1)
}

func main() {
	log.SetPrefix("[dupi] ")
	log.SetFlags(log.LstdFlags)
	if err := agent.Listen(agent.Options{
		ShutdownCleanup: true, // automatically closes on os.Interrupt
	}); err != nil {
		log.Fatal(err)
	}
	gargs, largs := splitArgs(os.Args[1:])
	if len(largs) == 0 {
		usageFatal(os.Stderr)
	}
	sc := scMap[largs[0]]
	if sc == nil {
		usageFatal(os.Stderr)
	}
	gFlags.Usage = func() { usageFatal(os.Stderr) }
	gFlags.Parse(gargs)
	if err := sc.Run(largs[1:]); err != nil {
		log.Fatal(err)
	}
}
