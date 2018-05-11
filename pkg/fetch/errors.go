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

import "go.thethings.network/lorawan-stack/pkg/errors"

var (
	// ErrFileFailedToOpen indicates the file could not be opened.
	ErrFileFailedToOpen = &errors.ErrDescriptor{
		MessageFormat:  "File `{filename}` failed to open",
		Code:           1,
		Type:           errors.Internal,
		SafeAttributes: []string{"filename"},
	}
	// ErrFileNotFound indicates the file could not be found.
	ErrFileNotFound = &errors.ErrDescriptor{
		MessageFormat:  "File `{filename}` not found",
		Code:           2,
		Type:           errors.NotFound,
		SafeAttributes: []string{"filename"},
	}
)

func init() {
	ErrFileFailedToOpen.Register()
	ErrFileNotFound.Register()
}
