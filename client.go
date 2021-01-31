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
	"net/http"

	"github.com/jadefox10200/httpext"
)

// A Client is an SEC client.  Its zero value (DefaultClient) is a usable client
// that uses http.DefaultClient.
type Client struct {
	client *http.Client
}

// NewClient returns a new Client given a HTTP client. The HTTP client will
// default to http.DefaultClient when nil.
func NewClient(c *http.Client) *Client {
	if c == nil {
		c = http.DefaultClient
	}
	rateLimitedClient := *c
	rateLimitedClient.Transport = httpext.NewLazyRateLimiter(c.Transport, nil)
	return &Client{&rateLimitedClient}
}

// DefaultClient is the default Client.
var DefaultClient = NewClient(nil)
