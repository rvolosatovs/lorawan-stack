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

package errors

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
)

// httpStatuscodes maps status codes to HTTP codes.
// See package google.golang.org/genproto/googleapis/rpc/code and google.golang.org/grpc/codes for details.
var httpStatuscodes = map[uint32]int{
	uint32(codes.OK):                 http.StatusOK,
	uint32(codes.Canceled):           499, // Client Closed Request
	uint32(codes.Unknown):            http.StatusInternalServerError,
	uint32(codes.InvalidArgument):    http.StatusBadRequest,
	uint32(codes.DeadlineExceeded):   http.StatusGatewayTimeout,
	uint32(codes.NotFound):           http.StatusNotFound,
	uint32(codes.AlreadyExists):      http.StatusConflict,
	uint32(codes.PermissionDenied):   http.StatusForbidden,
	uint32(codes.Unauthenticated):    http.StatusUnauthorized,
	uint32(codes.ResourceExhausted):  http.StatusTooManyRequests,
	uint32(codes.FailedPrecondition): http.StatusBadRequest,
	uint32(codes.Aborted):            http.StatusConflict,
	uint32(codes.OutOfRange):         http.StatusBadRequest,
	uint32(codes.Unimplemented):      http.StatusNotImplemented,
	uint32(codes.Internal):           http.StatusInternalServerError,
	uint32(codes.Unavailable):        http.StatusServiceUnavailable,
	uint32(codes.DataLoss):           http.StatusInternalServerError,
}

// HTTPStatusCode maps an error to HTTP response codes.
func HTTPStatusCode(err error) int {
	if status, ok := httpStatuscodes[Code(err)]; ok {
		return status
	}
	return http.StatusInternalServerError
}

// ToHTTP writes the error to the HTTP response.
func ToHTTP(in error, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	if ttnErr, ok := From(in); ok {
		w.WriteHeader(HTTPStatusCode(ttnErr))
		return json.NewEncoder(w).Encode(ttnErr)
	}
	w.WriteHeader(http.StatusInternalServerError)
	return json.NewEncoder(w).Encode(in)
}

// FromHTTP reads an error from the HTTP response.
func FromHTTP(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}
	defer resp.Body.Close()
	var err Error
	if decErr := json.NewDecoder(resp.Body).Decode(&err); decErr != nil {
		return decErr
	}
	return err
}
