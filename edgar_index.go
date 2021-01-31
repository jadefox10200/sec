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
	"compress/gzip"
	"fmt"
	"io"
	"time"

	"github.com/jadefox10200/httpext"
)

// ParseEDGARIndex parses an EDGAR index read from r, calling f for each entry.
func ParseEDGARIndex(r io.Reader, f func(EDGARIndexEntry) error) error {
	scanner := bufio.NewScanner(r)

	// Skip 11 lines.
	for i := 0; scanner.Scan() && i < 11; i++ {
	}

	// Parse index entries.
	for scanner.Scan() {
		e, err := ParseEDGARIndexEntry(scanner.Text())
		if err != nil {
			return err
		}

		if err := f(*e); err != nil {
			return err
		}
	}

	return scanner.Err()
}

// GetEDGARIndexEntries gets EDGAR index entries between start and end, calling f for
// each entry. The end time will default to the current time when zero.
//
// GetEDGARIndexEntries is a wrapper around DefaultClient.GetEDGARIndexEntries.
//
// See: https://www.sec.gov/edgar/searchedgar/accessing-edgar-data.htm
func GetEDGARIndexEntries(start, end time.Time, f func(EDGARIndexEntry) error) error {
	return DefaultClient.GetEDGARIndexEntries(start, end, f)
}

// GetEDGARIndexEntries gets EDGAR index entries between start and end, calling
// f for each entry. The end time will default to the current time when zero.
//
// See: https://www.sec.gov/edgar/searchedgar/accessing-edgar-data.htm
func (c *Client) GetEDGARIndexEntries(start, end time.Time, f func(EDGARIndexEntry) error) error {
	// Use DefaultClient if nil.
	if c == nil {
		c = DefaultClient
	}

	// Default the end time to the current time when zero.
	if end.IsZero() {
		end = time.Now()
	}

	for year := end.Year(); year >= start.Year(); year-- {
		endQuarter := 4
		if year == end.Year() {
			switch end.Month() {
			case 1, 2, 3:
				endQuarter = 1
			case 4, 5, 6:
				endQuarter = 2
			case 7, 8, 9:
				endQuarter = 3
			case 10, 11, 12:
				endQuarter = 4
			}
		}

		startQuarter := 1
		if year == start.Year() {
			switch start.Month() {
			case 1, 2, 3:
				startQuarter = 1
			case 4, 5, 6:
				startQuarter = 2
			case 7, 8, 9:
				startQuarter = 3
			case 10, 11, 12:
				startQuarter = 4
			}
		}

		for quarter := endQuarter; quarter >= startQuarter; quarter-- {
			url := fmt.Sprintf("https://www.sec.gov/Archives/edgar/full-index/%d/QTR%d/master.gz", year, quarter)
			resp, err := c.client.Get(url)
			if err != nil {
				return err
			}
			if httpext.IsErrorStatus(resp.StatusCode) {
				resp.Body.Close()
				return httpext.StatusError{
					URL:        url,
					StatusCode: resp.StatusCode,
				}
			}

			zr, err := gzip.NewReader(resp.Body)
			if err != nil {
				resp.Body.Close()
				if err == io.EOF {
					return nil
				}
				return err
			}

			if err := ParseEDGARIndex(zr, func(e EDGARIndexEntry) error {
				if !e.DateFiled.Before(start) && !e.DateFiled.After(end) {
					return f(e)
				}
				return nil
			}); err != nil {
				zr.Close()
				resp.Body.Close()
				return err
			}

			if err := zr.Close(); err != nil {
				resp.Body.Close()
				return err
			}

			if err := resp.Body.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}
