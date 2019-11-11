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
	"errors"
	"strconv"
	"strings"
	"time"
)

// An EDGARIndexEntry represents an entry in an a EDGAR index.
type EDGARIndexEntry struct {
	CIK         int
	CompanyName string
	FormType    string
	DateFiled   time.Time
	Filename    string
}

// ParseEDGARIndexEntry parses s as an entry in an EDGAR index.
func ParseEDGARIndexEntry(s string) (*EDGARIndexEntry, error) {
	cols := strings.Split(s, "|")
	if len(cols) != 5 {
		return nil, errors.New("sec.ParseEDGARIndexEntry: wrong number of columns for EDGAR index entry")
	}

	cik, err := strconv.Atoi(cols[0])
	if err != nil {
		return nil, err
	}

	dateFiled, err := time.Parse("2006-01-02", cols[3])
	if err != nil {
		return nil, err
	}

	return &EDGARIndexEntry{
		CIK:         cik,
		CompanyName: cols[1],
		FormType:    cols[2],
		DateFiled:   dateFiled,
		Filename:    cols[4],
	}, nil
}

// URL returns the URL for the EDGAR index entry.
func (e EDGARIndexEntry) URL() string {
	return "https://www.sec.gov/Archives/" + e.Filename
}

func (e EDGARIndexEntry) String() string {
	return strings.Join([]string{
		strconv.Itoa(e.CIK),
		e.CompanyName,
		e.FormType,
		e.DateFiled.Format("2006-01-02"),
		e.Filename,
	}, "|")
}
