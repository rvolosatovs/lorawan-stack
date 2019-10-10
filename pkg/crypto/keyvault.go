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

package crypto

import (
	"context"
	"crypto/tls"
	"crypto/x509"
)

// KeyVault provides wrapping and unwrapping keys using KEK labels.
type KeyVault interface {
	ComponentKEKLabeler

	Wrap(ctx context.Context, plaintext []byte, kekLabel string) ([]byte, error)
	Unwrap(ctx context.Context, ciphertext []byte, kekLabel string) ([]byte, error)

	// GetCertificate gets the X.509 certificate of the given identifier.
	GetCertificate(ctx context.Context, id string) (*x509.Certificate, error)
	// ExportCertificate exports the X.509 certificate and private key of the given identifier.
	ExportCertificate(ctx context.Context, id string) (*tls.Certificate, error)
}
