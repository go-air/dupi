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

package trigram

import (
	"fmt"
	"testing"

	"github.com/go-air/dupi/token"
)

func TestTrigram(t *testing.T) {
	txt := []byte("I am a fox.")
	toks := token.Tokenize(nil, txt, 0)
	for i := range toks {
		fmt.Printf("tok: %s\n", &toks[i])
	}
	tris := FromTokens(nil, toks)

	for _, tri := range tris {
		fmt.Printf("%s\n", tri)
	}

}
