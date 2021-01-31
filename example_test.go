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

package sec_test

import (
	"fmt"
	"log"
	"time"

	"github.com/jadefox10200/sec"
)

func ExampleGetEDGARIndexEntries() {
	end := time.Now()
	start := end.AddDate(0, -1, 0)
	if err := sec.GetEDGARIndexEntries(start, end, func(e sec.EDGARIndexEntry) error {
		fmt.Printf("%+v\n", e)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_GetEDGARIndexEntries() {
	c := sec.NewClient(nil)
	end := time.Now()
	start := end.AddDate(0, -1, 0)
	if err := c.GetEDGARIndexEntries(start, end, func(e sec.EDGARIndexEntry) error {
		fmt.Printf("%+v\n", e)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func ExampleGetForm4Filings() {
	end := time.Now()
	start := end.AddDate(0, -1, 0)
	if err := sec.GetForm4Filings(start, end, func(form sec.Form4) error {
		fmt.Printf("%+v\n", form)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_GetForm4Filings() {
	c := sec.NewClient(nil)
	end := time.Now()
	start := end.AddDate(0, -1, 0)
	if err := c.GetForm4Filings(start, end, func(form sec.Form4) error {
		fmt.Printf("%+v\n", form)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
