// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
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

package gatewayconfigurationserver

import (
	"fmt"
	"net/http"
	"strings"
)

// rewriteAuthorization rewrites the Authorization header from The Things Network Stack V2 style to The Things Stack.
// Packet Forwarders designed for The Things Stack Network V2 pass the gateway access key via the Authorization header
// prepended by `key`, which this function rewrites to `bearer`.
func rewriteAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := r.Header.Get("Authorization")
		parts := strings.SplitN(value, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "key" {
			r.Header.Set("Authorization", fmt.Sprintf("bearer %v", parts[1]))
		}
		next.ServeHTTP(w, r)
	})
}
