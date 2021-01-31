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
	"encoding/xml"
	"io"
	"time"

	"github.com/jadefox10200/httpext"
	"github.com/jadefox10200/marshaler"
)

// A Form4 represents a SEC form 4 filing.
type Form4 struct {
	XMLName                         xml.Name           `xml:"ownershipDocument"`
	PeriodOfReport                  marshaler.Date     `xml:"periodOfReport"`
	IssuerCIK                       int                `xml:"issuer>issuerCik"`
	IssuerName                      string             `xml:"issuer>issuerName"`
	IssuerTradingSymbol             string             `xml:"issuer>issuerTradingSymbol"`
	ReportingOwnerCIK               int                `xml:"reportingOwner>reportingOwnerId>rptOwnerCik"`
	ReportingOwnerName              string             `xml:"reportingOwner>reportingOwnerId>rptOwnerName"`
	ReportingOwnerTitle             string             `xml:"reportingOwner>reportingOwnerId>rptOwnerTitle"`
	ReportingOwnerIsDirector        bool               `xml:"reportingOwner>reportingOwnerRelationship>isDirector"`
	ReportingOwnerIsOfficer         bool               `xml:"reportingOwner>reportingOwnerRelationship>isOfficer"`
	ReportingOwnerIsTenPercentOwner bool               `xml:"reportingOwner>reportingOwnerRelationship>isTenPercentOwner"`
	NonDerivativeTransactions       []Form4Transaction `xml:"nonDerivativeTable>nonDerivativeTransaction"`
	DeriviativeTransactions         []Form4Transaction `xml:"derivativeTable>derivativeTransaction"`
}

// A Form4Transaction represents a transaction in a SEC form 4 filing.
type Form4Transaction struct {
	SecurityTitle                   string                  `xml:"securityTitle>value"`
	Date                            marshaler.Date          `xml:"transactionDate>value"`
	ConversionOrExercisePrice       marshaler.RobustFloat64 `xml:"conversionOrExercisePrice"`
	FormType                        string                  `xml:"transactionCoding>transactionFormType"`
	TransactionCode                 string                  `xml:"transactionCoding>transactionCode"`
	EquitySwapInvolved              bool                    `xml:"transactionCoding>equitySwapInvolved"`
	Shares                          float64                 `xml:"transactionAmounts>transactionShares>value"`
	PricePerShare                   marshaler.RobustFloat64 `xml:"transactionAmounts>transactionPricePerShare>value"`
	AcquiredDisposedCode            string                  `xml:"transactionAmounts>transactionAcquiredDisposedCode>value"`
	SharesOwnedFollowingTransaction float64                 `xml:"postTransactionAmounts>sharesOwnedFollowingTransaction>value"`
	DirectOrIndirectOwnership       string                  `xml:"ownershipNature>directOrIndirectOwnership>value"`
}

// ParseForm4 parses a form 4 filing read from r.
func ParseForm4(r io.Reader) (*Form4, error) {
	var form Form4
	if err := xml.NewDecoder(r).Decode(&form); err != nil {
		return nil, err
	}
	return &form, nil
}

// ParseForm4FromSECDocument parses a form 4 filing from an SEC document read
// from r.
func ParseForm4FromSECDocument(r io.Reader) (*Form4, error) {
	r, err := ExtractTagFromSECDocument(r, "XML")
	if err != nil {
		return nil, err
	}
	return ParseForm4(r)
}

// GetForm4Filings gets form 4 filings between start and end, calling f for each
// filing. The end time will default to the current time when zero.
//
// GetForm4Filings is a wrapper around DefaultClient.GetForm4Filings.
func GetForm4Filings(start, end time.Time, f func(Form4) error) error {
	return DefaultClient.GetForm4Filings(start, end, f)
}

// GetForm4Filings gets form 4 filings between start and end, calling f for each
// filing. The end time will default to the current time when zero.
func (c *Client) GetForm4Filings(start, end time.Time, f func(Form4) error) error {
	// Use DefaultClient when nil.
	if c == nil {
		c = DefaultClient
	}

	return c.GetEDGARIndexEntries(start, end, func(e EDGARIndexEntry) error {
		// Skip all forms except form 4 filings and amended form 4 filings.
		if e.FormType != FormType4 && e.FormType != FormType4A {
			return nil
		}

		// Send an HTTP request.
		url := e.URL()
		resp, err := c.client.Get(url)
		if err != nil {
			return nil
		}
		if httpext.IsErrorStatus(resp.StatusCode) {
			resp.Body.Close()
			return httpext.StatusError{URL: url, StatusCode: resp.StatusCode}
		}

		// Parse the form 4 filing from the SEC document.
		form, err := ParseForm4FromSECDocument(resp.Body)
		if err != nil {
			resp.Body.Close()
			return err
		}

		// Call f with the filing.
		if err := f(*form); err != nil {
			resp.Body.Close()
			return err
		}

		return resp.Body.Close()
	})
}
