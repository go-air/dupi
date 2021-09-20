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

import "fmt"

type Stats struct {
	Root      string
	NumDocs   uint64
	NumPaths  uint64
	NumPosts  uint64
	NumBlots  uint64
	BlotMean  float64
	BlotSigma float64
}

const stFmt = `dupi index at %s:
	- %d docs
	- %d nodes in path tree
	- %d posts
	- %d blots
	- %.2f mean docs per blot
	- %.2f sigma (std deviation)
`

func (st *Stats) String() string {
	return fmt.Sprintf(stFmt, st.Root, st.NumDocs,
		st.NumPaths, st.NumPosts, st.NumBlots,
		st.BlotMean, st.BlotSigma)
}
