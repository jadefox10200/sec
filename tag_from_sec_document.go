// Copyright 2019 Miles Barr <milesbarr2@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sec

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// tagFromSECDocumentReader reads a tag from an SEC document.
type tagFromSECDocumentReader struct {
	scanner *bufio.Scanner
	tag     string
}

// Read implements the io.Reader interface.
func (r *tagFromSECDocumentReader) Read(p []byte) (n int, err error) {
	if r.scanner == nil {
		return 0, io.EOF
	}

	if !r.scanner.Scan() {
		if err := r.scanner.Err(); err != nil {
			return 0, err
		}
		return 0, io.EOF
	}

	// Read until the tag ends.
	if strings.TrimSpace(r.scanner.Text()) == "</"+r.tag+">" {
		r.scanner = nil
		return 0, io.EOF
	}

	copy(p, r.scanner.Bytes())
	return len(r.scanner.Bytes()), nil
}

// ExtractTagFromSECDocument extracts a tag from an SEC document read from r,
// returning a reader to the tag content.
func ExtractTagFromSECDocument(r io.Reader, tag string) (io.Reader, error) {
	scanner := bufio.NewScanner(r)

	// Skip until the tag starts.
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "<"+tag+">" {
			found = true
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf(
			"sec.ExtractTagFromSECDocument: missing tag \"%s\"", tag)
	}
	return &tagFromSECDocumentReader{scanner, tag}, nil
}
