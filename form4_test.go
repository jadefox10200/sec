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
	"compress/gzip"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/tradyfinance/httpext"
	"github.com/tradyfinance/marshaler"
)

// See: https://www.sec.gov/Archives/edgar/data/1000045/0001357521-18-000008.txt
const sampleForm4XML = `
<?xml version="1.0"?>
<ownershipDocument>

    <schemaVersion>X0306</schemaVersion>

    <documentType>4</documentType>

    <periodOfReport>2018-10-15</periodOfReport>

    <notSubjectToSection16>0</notSubjectToSection16>

    <issuer>
        <issuerCik>0001000045</issuerCik>
        <issuerName>NICHOLAS FINANCIAL INC</issuerName>
        <issuerTradingSymbol>NICK</issuerTradingSymbol>
    </issuer>

    <reportingOwner>
        <reportingOwnerId>
            <rptOwnerCik>0001357521</rptOwnerCik>
            <rptOwnerName>MALSON KELLY M</rptOwnerName>
        </reportingOwnerId>
        <reportingOwnerAddress>
            <rptOwnerStreet1>2454 MCMULLEN BOOTH ROAD</rptOwnerStreet1>
            <rptOwnerStreet2>BUILDING C</rptOwnerStreet2>
            <rptOwnerCity>CLEARWATER</rptOwnerCity>
            <rptOwnerState>FL</rptOwnerState>
            <rptOwnerZipCode>33759</rptOwnerZipCode>
            <rptOwnerStateDescription></rptOwnerStateDescription>
        </reportingOwnerAddress>
        <reportingOwnerRelationship>
            <isDirector>0</isDirector>
            <isOfficer>1</isOfficer>
            <isTenPercentOwner>0</isTenPercentOwner>
            <isOther>0</isOther>
            <officerTitle>CFO</officerTitle>
        </reportingOwnerRelationship>
    </reportingOwner>

    <nonDerivativeTable>
        <nonDerivativeTransaction>
            <securityTitle>
                <value>Common</value>
            </securityTitle>
            <transactionDate>
                <value>2018-10-15</value>
            </transactionDate>
            <transactionCoding>
                <transactionFormType>4</transactionFormType>
                <transactionCode>P</transactionCode>
                <equitySwapInvolved>0</equitySwapInvolved>
            </transactionCoding>
            <transactionTimeliness>
                <value></value>
            </transactionTimeliness>
            <transactionAmounts>
                <transactionShares>
                    <value>1569</value>
                    <footnoteId id="F1"/>
                </transactionShares>
                <transactionPricePerShare>
                    <value>11.98</value>
                    <footnoteId id="F2"/>
                </transactionPricePerShare>
                <transactionAcquiredDisposedCode>
                    <value>A</value>
                </transactionAcquiredDisposedCode>
            </transactionAmounts>
            <postTransactionAmounts>
                <sharesOwnedFollowingTransaction>
                    <value>15989</value>
                </sharesOwnedFollowingTransaction>
            </postTransactionAmounts>
            <ownershipNature>
                <directOrIndirectOwnership>
                    <value>D</value>
                </directOrIndirectOwnership>
            </ownershipNature>
        </nonDerivativeTransaction>
        <nonDerivativeTransaction>
            <securityTitle>
                <value>Common</value>
            </securityTitle>
            <transactionDate>
                <value>2018-10-15</value>
            </transactionDate>
            <transactionCoding>
                <transactionFormType>4</transactionFormType>
                <transactionCode>A</transactionCode>
                <equitySwapInvolved>0</equitySwapInvolved>
                <footnoteId id="F3"/>
            </transactionCoding>
            <transactionTimeliness>
                <value></value>
            </transactionTimeliness>
            <transactionAmounts>
                <transactionShares>
                    <value>1569</value>
                </transactionShares>
                <transactionPricePerShare>
                    <value>0</value>
                    <footnoteId id="F3"/>
                </transactionPricePerShare>
                <transactionAcquiredDisposedCode>
                    <value>A</value>
                </transactionAcquiredDisposedCode>
            </transactionAmounts>
            <postTransactionAmounts>
                <sharesOwnedFollowingTransaction>
                    <value>17558</value>
                </sharesOwnedFollowingTransaction>
            </postTransactionAmounts>
            <ownershipNature>
                <directOrIndirectOwnership>
                    <value>D</value>
                </directOrIndirectOwnership>
            </ownershipNature>
        </nonDerivativeTransaction>
    </nonDerivativeTable>

    <footnotes>
        <footnote id="F1">Purchases of shares was made in accordance with a 10b5-1 Plan previously executed.</footnote>
        <footnote id="F2">Represents the average purchase price.</footnote>
        <footnote id="F3">These shares were awarded pursuant to the reporting person's employment agreement.  The closing stock price of the issuer's common stock on NASDAQ on 10/15/2018 was $11.79.</footnote>
    </footnotes>

    <ownerSignature>
        <signatureName>/s/ Kelly M. Malson</signatureName>
        <signatureDate>2018-10-15</signatureDate>
    </ownerSignature>
</ownershipDocument>`

// See: https://www.sec.gov/Archives/edgar/data/1000045/0001357521-18-000008.txt
const sampleForm4SECDocument = `
<SEC-DOCUMENT>0001357521-18-000008.txt : 20181015
<SEC-HEADER>0001357521-18-000008.hdr.sgml : 20181015
<ACCEPTANCE-DATETIME>20181015171243
ACCESSION NUMBER:		0001357521-18-000008
CONFORMED SUBMISSION TYPE:	4
PUBLIC DOCUMENT COUNT:		1
CONFORMED PERIOD OF REPORT:	20181015
FILED AS OF DATE:		20181015
DATE AS OF CHANGE:		20181015

REPORTING-OWNER:	

	OWNER DATA:	
		COMPANY CONFORMED NAME:			MALSON KELLY M
		CENTRAL INDEX KEY:			0001357521

	FILING VALUES:
		FORM TYPE:		4
		SEC ACT:		1934 Act
		SEC FILE NUMBER:	000-26680
		FILM NUMBER:		181122886

	MAIL ADDRESS:	
		STREET 1:		2454 MCMULLEN BOOTH ROAD
		STREET 2:		BUILDING C
		CITY:			CLEARWATER
		STATE:			FL
		ZIP:			33759

	FORMER NAME:	
		FORMER CONFORMED NAME:	Snape Kelly Malson
		DATE OF NAME CHANGE:	20060327

ISSUER:		

	COMPANY DATA:	
		COMPANY CONFORMED NAME:			NICHOLAS FINANCIAL INC
		CENTRAL INDEX KEY:			0001000045
		STANDARD INDUSTRIAL CLASSIFICATION:	SHORT-TERM BUSINESS CREDIT INSTITUTIONS [6153]
		IRS NUMBER:				593019317
		STATE OF INCORPORATION:			FL
		FISCAL YEAR END:			0331

	BUSINESS ADDRESS:	
		STREET 1:		2454 MCMULLEN BOOTH RD
		STREET 2:		BLDG C SUITE 501 B
		CITY:			CLEARWATER
		STATE:			FL
		ZIP:			33759
		BUSINESS PHONE:		7277260763

	MAIL ADDRESS:	
		STREET 1:		2454 MCMULLEN BOOTH RD
		STREET 2:		BLDG C SUITE 501B
		CITY:			CLEARWATER
		STATE:			FL
		ZIP:			33759
</SEC-HEADER>
<DOCUMENT>
<TYPE>4
<SEQUENCE>1
<FILENAME>primary_doc.xml
<DESCRIPTION>PRIMARY DOCUMENT
<TEXT>
<XML>
<?xml version="1.0"?>
<ownershipDocument>

    <schemaVersion>X0306</schemaVersion>

    <documentType>4</documentType>

    <periodOfReport>2018-10-15</periodOfReport>

    <notSubjectToSection16>0</notSubjectToSection16>

    <issuer>
        <issuerCik>0001000045</issuerCik>
        <issuerName>NICHOLAS FINANCIAL INC</issuerName>
        <issuerTradingSymbol>NICK</issuerTradingSymbol>
    </issuer>

    <reportingOwner>
        <reportingOwnerId>
            <rptOwnerCik>0001357521</rptOwnerCik>
            <rptOwnerName>MALSON KELLY M</rptOwnerName>
        </reportingOwnerId>
        <reportingOwnerAddress>
            <rptOwnerStreet1>2454 MCMULLEN BOOTH ROAD</rptOwnerStreet1>
            <rptOwnerStreet2>BUILDING C</rptOwnerStreet2>
            <rptOwnerCity>CLEARWATER</rptOwnerCity>
            <rptOwnerState>FL</rptOwnerState>
            <rptOwnerZipCode>33759</rptOwnerZipCode>
            <rptOwnerStateDescription></rptOwnerStateDescription>
        </reportingOwnerAddress>
        <reportingOwnerRelationship>
            <isDirector>0</isDirector>
            <isOfficer>1</isOfficer>
            <isTenPercentOwner>0</isTenPercentOwner>
            <isOther>0</isOther>
            <officerTitle>CFO</officerTitle>
        </reportingOwnerRelationship>
    </reportingOwner>

    <nonDerivativeTable>
        <nonDerivativeTransaction>
            <securityTitle>
                <value>Common</value>
            </securityTitle>
            <transactionDate>
                <value>2018-10-15</value>
            </transactionDate>
            <transactionCoding>
                <transactionFormType>4</transactionFormType>
                <transactionCode>P</transactionCode>
                <equitySwapInvolved>0</equitySwapInvolved>
            </transactionCoding>
            <transactionTimeliness>
                <value></value>
            </transactionTimeliness>
            <transactionAmounts>
                <transactionShares>
                    <value>1569</value>
                    <footnoteId id="F1"/>
                </transactionShares>
                <transactionPricePerShare>
                    <value>11.98</value>
                    <footnoteId id="F2"/>
                </transactionPricePerShare>
                <transactionAcquiredDisposedCode>
                    <value>A</value>
                </transactionAcquiredDisposedCode>
            </transactionAmounts>
            <postTransactionAmounts>
                <sharesOwnedFollowingTransaction>
                    <value>15989</value>
                </sharesOwnedFollowingTransaction>
            </postTransactionAmounts>
            <ownershipNature>
                <directOrIndirectOwnership>
                    <value>D</value>
                </directOrIndirectOwnership>
            </ownershipNature>
        </nonDerivativeTransaction>
        <nonDerivativeTransaction>
            <securityTitle>
                <value>Common</value>
            </securityTitle>
            <transactionDate>
                <value>2018-10-15</value>
            </transactionDate>
            <transactionCoding>
                <transactionFormType>4</transactionFormType>
                <transactionCode>A</transactionCode>
                <equitySwapInvolved>0</equitySwapInvolved>
                <footnoteId id="F3"/>
            </transactionCoding>
            <transactionTimeliness>
                <value></value>
            </transactionTimeliness>
            <transactionAmounts>
                <transactionShares>
                    <value>1569</value>
                </transactionShares>
                <transactionPricePerShare>
                    <value>0</value>
                    <footnoteId id="F3"/>
                </transactionPricePerShare>
                <transactionAcquiredDisposedCode>
                    <value>A</value>
                </transactionAcquiredDisposedCode>
            </transactionAmounts>
            <postTransactionAmounts>
                <sharesOwnedFollowingTransaction>
                    <value>17558</value>
                </sharesOwnedFollowingTransaction>
            </postTransactionAmounts>
            <ownershipNature>
                <directOrIndirectOwnership>
                    <value>D</value>
                </directOrIndirectOwnership>
            </ownershipNature>
        </nonDerivativeTransaction>
    </nonDerivativeTable>

    <footnotes>
        <footnote id="F1">Purchases of shares was made in accordance with a 10b5-1 Plan previously executed.</footnote>
        <footnote id="F2">Represents the average purchase price.</footnote>
        <footnote id="F3">These shares were awarded pursuant to the reporting person's employment agreement.  The closing stock price of the issuer's common stock on NASDAQ on 10/15/2018 was $11.79.</footnote>
    </footnotes>

    <ownerSignature>
        <signatureName>/s/ Kelly M. Malson</signatureName>
        <signatureDate>2018-10-15</signatureDate>
    </ownerSignature>
</ownershipDocument>
</XML>
</TEXT>
</DOCUMENT>
</SEC-DOCUMENT>`

var sampleForm4 = &Form4{
	XMLName:                 xml.Name{Local: "ownershipDocument"},
	PeriodOfReport:          marshaler.Date(time.Date(2018, 10, 15, 0, 0, 0, 0, time.UTC)),
	IssuerCIK:               1000045,
	IssuerName:              "NICHOLAS FINANCIAL INC",
	IssuerTradingSymbol:     "NICK",
	ReportingOwnerCIK:       1357521,
	ReportingOwnerName:      "MALSON KELLY M",
	ReportingOwnerIsOfficer: true,
	NonDerivativeTransactions: []Form4Transaction{
		Form4Transaction{
			SecurityTitle: "Common",
			Date:          marshaler.Date(time.Date(2018, 10, 15, 0, 0, 0, 0, time.UTC)),
			ConversionOrExercisePrice:       0.000000,
			FormType:                        "4",
			TransactionCode:                 "P",
			Shares:                          1569,
			PricePerShare:                   11.980000,
			AcquiredDisposedCode:            "A",
			SharesOwnedFollowingTransaction: 15989,
			DirectOrIndirectOwnership:       "D",
		},
		Form4Transaction{
			SecurityTitle: "Common",
			Date:          marshaler.Date(time.Date(2018, 10, 15, 0, 0, 0, 0, time.UTC)),
			ConversionOrExercisePrice:       0.000000,
			FormType:                        "4",
			TransactionCode:                 "A",
			Shares:                          1569,
			PricePerShare:                   0.000000,
			AcquiredDisposedCode:            "A",
			SharesOwnedFollowingTransaction: 17558,
			DirectOrIndirectOwnership:       "D",
		},
	},
}

func TestParseForm4Filing(t *testing.T) {
	got, err := ParseForm4(strings.NewReader(sampleForm4XML))
	if err != nil {
		t.Fatal(err)
	}
	if want := sampleForm4; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestParseForm4FilingFromSECDocument(t *testing.T) {
	got, err := ParseForm4FromSECDocument(strings.NewReader(sampleForm4SECDocument))
	if err != nil {
		t.Fatal(err)
	}
	if want := sampleForm4; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestClient_GetForm4Filings(t *testing.T) {
	c := NewClient(httpext.WithTransportFunc(nil, func(req *http.Request) (*http.Response, error) {
		var res http.Response
		if strings.Contains(req.URL.Path, "edgar/data") {
			res.Body = ioutil.NopCloser(strings.NewReader(sampleForm4SECDocument))
		} else {
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
			res.Body = r
		}
		return &res, nil
	}))
	start := time.Date(2018, 10, 15, 0, 0, 0, 0, time.UTC)
	end := time.Date(2018, 10, 19, 0, 0, 0, 0, time.UTC)
	got := []Form4{}
	if err := c.GetForm4Filings(start, end, func(form Form4) error {
		got = append(got, form)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if want := []Form4{*sampleForm4}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}
