// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
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

package fetch

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gregjones/httpcache"
	errors "go.thethings.network/lorawan-stack/pkg/errorsv3"
)

type httpFetcher struct {
	baseFetcher
	transport *http.Client
}

func (f httpFetcher) File(pathElements ...string) ([]byte, error) {
	start := time.Now()
	url := fmt.Sprintf("%s/%s", f.base, strings.TrimLeft(path.Join(pathElements...), "/"))

	resp, err := f.transport.Get(url)
	if err != nil {
		return nil, err
	}

	if err = errors.FromHTTP(resp); err != nil {
		return nil, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return nil, err
	}

	f.observeLatency(time.Since(start))
	return result, nil
}

// FromHTTP returns an object to fetch files from a webserver.
func FromHTTP(baseURL string, cache bool) Interface {
	baseURL = strings.TrimRight(baseURL, "/")
	f := httpFetcher{
		baseFetcher: baseFetcher{
			base:    baseURL,
			latency: fetchLatency.WithLabelValues("http", baseURL),
		},
		transport: http.DefaultClient,
	}
	if !cache {
		f.transport = httpcache.NewMemoryCacheTransport().Client()
	}
	return f
}
