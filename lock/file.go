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

package lock

import "os"

type File struct {
	path   string
	handle *os.File
}

func New(path string) (*File, error) {
	f, e := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if e != nil {
		return nil, e
	}
	return &File{path, f}, nil
}

// Close unlocks and then closes the file, returning any
// error.  The file handle is closed whether or not
// unlocking fails with an error.
func (f *File) Close() error {
	erru := f.Unlock()
	errc := f.handle.Close()
	if erru == nil {
		return errc
	}
	return erru
}
