// Copyright © 2021 The Things Network Foundation, The Things Industries B.V.
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

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	urlutil "go.thethings.network/lorawan-stack/v3/pkg/util/url"
	"go.thethings.network/lorawan-stack/v3/pkg/version"
)

// Option is an option for the API client.
type Option interface {
	apply(*Client)
}

// OptionFunc is an Option implemented as a function.
type OptionFunc func(*Client)

func (f OptionFunc) apply(c *Client) { f(c) }

// Client is an API client for the LoRa Cloud Device Management v1 service.
type Client struct {
	token   string
	baseURL *url.URL
	cl      *http.Client
}

const (
	contentType      = "application/json"
	defaultServerURL = "https://gls.loracloud.com"
	basePath         = "/api/v3"
)

var (
	userAgent        = "ttn-lw-application-server/" + version.TTN
	DefaultServerURL *url.URL
)

func (c *Client) newRequest(ctx context.Context, method, category, operation string, body io.Reader) (*http.Request, error) {
	u := urlutil.CloneURL(c.baseURL)
	u.Path = path.Join(basePath, category, operation)
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", userAgent)
	if c.token != "" {
		req.Header.Set("Ocp-Apim-Subscription-Key", c.token)
	}
	return req, nil
}

// Do executes a new HTTP request with the given parameters and body and returns the response.
func (c *Client) Do(ctx context.Context, method, category, operation string, body io.Reader) (*http.Response, error) {
	req, err := c.newRequest(ctx, method, category, operation, body)
	if err != nil {
		return nil, err
	}
	return c.cl.Do(req)
}

// WithToken uses the given authentication token in the client.
func WithToken(token string) Option {
	return OptionFunc(func(c *Client) {
		c.token = token
	})
}

// WithBaseURL uses the given base URL for the requests of the client.
func WithBaseURL(baseURL *url.URL) Option {
	return OptionFunc(func(c *Client) {
		c.baseURL = baseURL
	})
}

// SolveSingleFrame attempts to solve the location of the end-device using the provided request.
func (c *Client) SolveSingleFrame(ctx context.Context, request *SingleFrameRequest) (*ExtendedSingleFrameResponse, error) {
	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(request); err != nil {
		return nil, err
	}
	resp, err := c.Do(ctx, http.MethodPost, "solve", "singleframe", buffer)
	if err != nil {
		return nil, err
	}
	response := &ExtendedSingleFrameResponse{}
	err = parse(&response, resp)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// New creates a new Client with the given options.
func New(cl *http.Client, opts ...Option) (*Client, error) {
	client := &Client{
		cl:      cl,
		baseURL: urlutil.CloneURL(DefaultServerURL),
	}
	for _, opt := range opts {
		opt.apply(client)
	}
	return client, nil
}

func init() {
	var err error
	DefaultServerURL, err = url.Parse(defaultServerURL)
	if err != nil {
		panic(fmt.Sprintf("loragls: failed to parse base URL: %v", err))
	}
}
