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

package rights

import "go.thethings.network/lorawan-stack/pkg/ttnpb"

// noopCache implements cache and caches nothing.
type noopCache struct{}

func (c *noopCache) GetOrFetch(key string, fetch func() ([]ttnpb.Right, error)) ([]ttnpb.Right, error) {
	return fetch()
}
