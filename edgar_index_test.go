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
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/tradyfinance/httpext"
)

const sampleEDGARIndex = `
Description:           Master Index of EDGAR Dissemination Feed
Last Data Received:    October 20, 2018
Comments:              webmaster@sec.gov
Anonymous FTP:         ftp://ftp.sec.gov/edgar/
Cloud HTTP:            https://www.sec.gov/Archives/

 
 
 
CIK|Company Name|Form Type|Date Filed|Filename
--------------------------------------------------------------------------------
1000045|NICHOLAS FINANCIAL INC|4|2018-10-15|edgar/data/1000045/0001357521-18-000008.txt
1000184|SAP SE|6-K|2018-10-19|edgar/data/1000184/0001104659-18-062851.txt
`

func TestParseEDGARIndex(t *testing.T) {
	got := []EDGARIndexEntry{}
	if err := ParseEDGARIndex(bytes.NewBufferString(sampleEDGARIndex), func(e EDGARIndexEntry) error {
		got = append(got, e)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if want := []EDGARIndexEntry{
		EDGARIndexEntry{
			CIK:         1000045,
			CompanyName: "NICHOLAS FINANCIAL INC",
			FormType:    "4",
			DateFiled:   time.Date(2018, 10, 15, 0, 0, 0, 0, time.UTC),
			Filename:    "edgar/data/1000045/0001357521-18-000008.txt",
		},
		EDGARIndexEntry{
			CIK:         1000184,
			CompanyName: "SAP SE",
			FormType:    "6-K",
			DateFiled:   time.Date(2018, 10, 19, 0, 0, 0, 0, time.UTC),
			Filename:    "edgar/data/1000184/0001104659-18-062851.txt",
		},
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestClient_GetEDGARIndexEntries(t *testing.T) {
	c := NewClient(httpext.WithTransportFunc(nil, func(req *http.Request) (*http.Response, error) {
		r, w := io.Pipe()
		go func() {
			gz := gzip.NewWriter(w)
			if _, err := gz.Write([]byte(sampleEDGARIndex)); err != nil {
				t.Fatal(err)
			}
			if err := gz.Close(); err != nil {
				t.Fatal(err)
			}
			if err := w.Close(); err != nil {
				t.Fatal(err)
			}
		}()
		var res http.Response
		res.Body = r
		return &res, nil
	}))
	start := time.Date(2018, 10, 15, 0, 0, 0, 0, time.UTC)
	end := time.Date(2018, 10, 19, 0, 0, 0, 0, time.UTC)
	got := []EDGARIndexEntry{}
	if err := c.GetEDGARIndexEntries(start, end, func(e EDGARIndexEntry) error {
		got = append(got, e)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if want := []EDGARIndexEntry{
		EDGARIndexEntry{
			CIK:         1000045,
			CompanyName: "NICHOLAS FINANCIAL INC",
			FormType:    "4",
			DateFiled:   start,
			Filename:    "edgar/data/1000045/0001357521-18-000008.txt",
		},
		EDGARIndexEntry{
			CIK:         1000184,
			CompanyName: "SAP SE",
			FormType:    "6-K",
			DateFiled:   end,
			Filename:    "edgar/data/1000184/0001104659-18-062851.txt",
		},
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}
