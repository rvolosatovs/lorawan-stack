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

package devicetemplates

import (
	"bytes"
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"strings"

	pbtypes "github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/gogoproto"
	"go.thethings.network/lorawan-stack/v3/pkg/provisioning"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
	jose "gopkg.in/square/go-jose.v2"
)

const (
	microchipATECC608AMAHTNTPart = "ATECC608A-MAHTN-T"
	microchipATECC608TNGLORAPart = "ATECC608A-TNGLORA"
)

var microchipPublicKeys = map[string]map[string][]byte{
	microchipATECC608AMAHTNTPart: {
		// Certificate:
		// Data:
		// 		Version: 3 (0x2)
		// 		Serial Number:
		// 				76:5f:7b:3e:38:02:7b:6b:8f:8a:18:49:5f:51:55:5d
		// Signature Algorithm: ecdsa-with-SHA256
		// 		Issuer: O=Microchip Technology Inc, CN=Log Signer Test
		// 		Validity
		// 				Not Before: Jan 18 20:29:41 2019 GMT
		// 				Not After : Feb 18 20:29:41 2019 GMT
		// 		Subject: O=Microchip Technology Inc, CN=Log Signer Test
		"B0uaDLLiyKe-SorP71F3BNoVrfY": []byte(`-----BEGIN CERTIFICATE-----
MIIByDCCAW6gAwIBAgIQdl97PjgCe2uPihhJX1FVXTAKBggqhkjOPQQDAjA9MSEw
HwYDVQQKDBhNaWNyb2NoaXAgVGVjaG5vbG9neSBJbmMxGDAWBgNVBAMMD0xvZyBT
aWduZXIgVGVzdDAeFw0xOTAxMTgyMDI5NDFaFw0xOTAyMTgyMDI5NDFaMD0xITAf
BgNVBAoMGE1pY3JvY2hpcCBUZWNobm9sb2d5IEluYzEYMBYGA1UEAwwPTG9nIFNp
Z25lciBUZXN0MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEHaTaoWgx1zG1JhnP
NbueEtfe926WJwbkHIyTBTB2aDzBUf3oRSFleYCJOEaRZlbEoQ4WiDDwCDeBd8GK
R70mN6NQME4wHQYDVR0OBBYEFAdLmgyy4sinvkqKz+9RdwTaFa32MB8GA1UdIwQY
MBaAFAdLmgyy4sinvkqKz+9RdwTaFa32MAwGA1UdEwEB/wQCMAAwCgYIKoZIzj0E
AwIDSAAwRQIhAI9jMSnc+HKKnjZ5ghmYVXYgPn9M9ae6gfE4AN5xekEZAiBNk7Pz
FVV78rUrxt7igKFg3mMLfE8Qeoh6dDKmRkbAEA==
-----END CERTIFICATE-----`),
		// Certificate:
		//     Data:
		//         Version: 3 (0x2)
		//         Serial Number:
		//             64:62:16:c8:c6:48:f5:c3:1c:05:98:a9:5f:14:ce:58
		//     Signature Algorithm: ecdsa-with-SHA256
		//         Issuer: O=Microchip Technology Inc, CN=Log Signer 001
		//         Validity
		//             Not Before: Jan 22 00:27:42 2019 GMT
		//             Not After : Jul 22 00:27:42 2019 GMT
		//         Subject: O=Microchip Technology Inc, CN=Log Signer 001
		"7cCILlAOwYo1-PChGuoyUISMK3g": []byte(`-----BEGIN CERTIFICATE-----
MIIBxjCCAWygAwIBAgIQZGIWyMZI9cMcBZipXxTOWDAKBggqhkjOPQQDAjA8MSEw
HwYDVQQKDBhNaWNyb2NoaXAgVGVjaG5vbG9neSBJbmMxFzAVBgNVBAMMDkxvZyBT
aWduZXIgMDAxMB4XDTE5MDEyMjAwMjc0MloXDTE5MDcyMjAwMjc0MlowPDEhMB8G
A1UECgwYTWljcm9jaGlwIFRlY2hub2xvZ3kgSW5jMRcwFQYDVQQDDA5Mb2cgU2ln
bmVyIDAwMTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABEu8/ZyRdTu4N0kuu76C
R1JR5vz04EuRqL4TQxMinRiUc3Htqy38O6HrXo2qmNoyrO0xd2I2pfQhXWYuLT35
MGWjUDBOMB0GA1UdDgQWBBTtwIguUA7BijX48KEa6jJQhIwreDAfBgNVHSMEGDAW
gBTtwIguUA7BijX48KEa6jJQhIwreDAMBgNVHRMBAf8EAjAAMAoGCCqGSM49BAMC
A0gAMEUCIQD9/x9zxmHkeWGwjEq67QsQqBVmoY8k6PvFVr4Bz1tYOwIgYfck+fv/
pno8+2vVTkQDhcinNrgoPLQORzV5/l/b4z4=
-----END CERTIFICATE-----`),
	},
	microchipATECC608TNGLORAPart: {
		// Certificate:
		//     Data:
		//         Version: 3 (0x2)
		//         Serial Number:
		//             73:a1:f2:32:3a:e4:4f:64:ce:63:5e:c5:01:7e:d7:56
		//     Signature Algorithm: ecdsa-with-SHA256
		//         Issuer: O=Microchip Technology Inc, CN=Log Signer 002
		//         Validity
		//             Not Before: Aug 15 19:47:59 2019 GMT
		//             Not After : Aug 15 19:47:59 2020 GMT
		// 				Subject: O=Microchip Technology Inc, CN=Log Signer 002
		"8VeKGdyU2d8wev6_VzNJOBOv-cA": []byte(`-----BEGIN CERTIFICATE-----
MIIBxzCCAWygAwIBAgIQc6HyMjrkT2TOY17FAX7XVjAKBggqhkjOPQQDAjA8MSEw
HwYDVQQKDBhNaWNyb2NoaXAgVGVjaG5vbG9neSBJbmMxFzAVBgNVBAMMDkxvZyBT
aWduZXIgMDAyMB4XDTE5MDgxNTE5NDc1OVoXDTIwMDgxNTE5NDc1OVowPDEhMB8G
A1UECgwYTWljcm9jaGlwIFRlY2hub2xvZ3kgSW5jMRcwFQYDVQQDDA5Mb2cgU2ln
bmVyIDAwMjBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABLCLrgPlT3OezntD9lC2
ShwUhlx07fiq/VETJ+ITUAwbgrPjB/Xi9GchLIM7FwZSUGOEqRA6KtH32XMpTGHK
mCCjUDBOMB0GA1UdDgQWBBTxV4oZ3JTZ3zB6/r9XM0k4E6/5wDAfBgNVHSMEGDAW
gBTxV4oZ3JTZ3zB6/r9XM0k4E6/5wDAMBgNVHRMBAf8EAjAAMAoGCCqGSM49BAMC
A0kAMEYCIQDKHgctLnq/zNqfB+1v0KRhDVPvRf6Dimt8aW9WLS0NWAIhAJvUe3uJ
pkMG4zpov9FCoj4G340idEadm7mVbAd5GOB9
-----END CERTIFICATE-----`),
	},
}

var joinEUI = types.EUI64{0x70, 0xb3, 0xd5, 0x7e, 0xd0, 0x00, 0x00, 0x00}

type microchipEntry struct {
	jose.JSONWebSignature
}

func (m *microchipEntry) UnmarshalJSON(data []byte) error {
	jws, err := jose.ParseSigned(string(data))
	if err != nil {
		return err
	}
	*m = microchipEntry{JSONWebSignature: *jws}
	return nil
}

var (
	errMicrochipData      = errors.DefineInvalidArgument("microchip_data", "invalid Microchip data")
	errMicrochipPublicKey = errors.DefineInvalidArgument("microchip_public_key", "unknown Microchip public key ID `{id}`")
)

// microchipATECC608AMAHTNT is a Microchip ATECC608A-MAHTN-T device provisioner.
type microchipATECC608AMAHTNT struct {
	keys map[string]interface{}
}

func (m *microchipATECC608AMAHTNT) Format() *ttnpb.EndDeviceTemplateFormat {
	return &ttnpb.EndDeviceTemplateFormat{
		Name:           "Microchip ATECC608A-MAHTN-T Manifest File",
		Description:    "JSON manifest file received through Microchip Purchasing & Client Services.",
		FileExtensions: []string{".json"},
	}
}

// Convert decodes the given manifest data.
// The input data is an array of JWS (JSON Web Signatures).
func (m *microchipATECC608AMAHTNT) Convert(ctx context.Context, r io.Reader, ch chan<- *ttnpb.EndDeviceTemplate) error {
	defer close(ch)

	dec := json.NewDecoder(r)
	delim, err := dec.Token()
	if err != nil {
		return errMicrochipData.WithCause(err)
	}
	if _, ok := delim.(json.Delim); !ok {
		return errMicrochipData.New()
	}

	for dec.More() {
		var jws microchipEntry
		if err := dec.Decode(&jws); err != nil {
			return errMicrochipData.WithCause(err)
		}
		kid := jws.Signatures[0].Header.KeyID
		key, ok := m.keys[kid]
		if !ok {
			return errMicrochipPublicKey.WithAttributes("id", kid)
		}
		buf, err := jws.Verify(key)
		if err != nil {
			return errMicrochipData.WithCause(err)
		}
		m := make(map[string]interface{})
		if err := json.Unmarshal(buf, &m); err != nil {
			return errMicrochipData.WithCause(err)
		}
		s, err := gogoproto.Struct(m)
		if err != nil {
			return errMicrochipData.WithCause(err)
		}
		sn := s.Fields["uniqueId"].GetStringValue()
		tmpl := &ttnpb.EndDeviceTemplate{
			EndDevice: ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					JoinEUI: &joinEUI,
				},
				ProvisionerID:    provisioning.Microchip,
				ProvisioningData: s,
				RootKeys: &ttnpb.RootKeys{
					RootKeyID: sn,
				},
				SupportsJoin: true,
			},
			FieldMask: pbtypes.FieldMask{
				Paths: []string{
					"ids.join_eui",
					"provisioner_id",
					"provisioning_data",
					"root_keys.root_key_id",
					"supports_join",
				},
			},
			MappingKey: sn,
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- tmpl:
		}
	}
	return nil
}

// microchipATECC608TNGLORA is a Microchip ATECC608A-TNGLORA-B and -C device provisioner.
type microchipATECC608TNGLORA struct {
	keys map[string]interface{}
}

func (m *microchipATECC608TNGLORA) Format() *ttnpb.EndDeviceTemplateFormat {
	return &ttnpb.EndDeviceTemplateFormat{
		Name:           "Microchip ATECC608A-TNGLORA Manifest File",
		Description:    "JSON manifest file received through Microchip Purchasing & Client Services.",
		FileExtensions: []string{".json"},
	}
}

var (
	errMicrochipNoCertificate  = errors.DefineInvalidArgument("microchip_no_certificate", "no Microchip certificate found")
	errMicrochipCertificateSAN = errors.DefineInvalidArgument("microchip_certificate_san", "invalid Microchip certificate Subject Alternate Name")
	errMicrochipUnknownDevEUI  = errors.DefineInvalidArgument("microchip_unknown_dev_eui", "unknown Microchip DevEUI")
)

// Convert decodes the given manifest data.
// The input data is an array of JWS (JSON Web Signatures).
func (m *microchipATECC608TNGLORA) Convert(ctx context.Context, r io.Reader, ch chan<- *ttnpb.EndDeviceTemplate) error {
	defer close(ch)

	dec := json.NewDecoder(r)
	delim, err := dec.Token()
	if err != nil {
		return errMicrochipData.WithCause(err)
	}
	if _, ok := delim.(json.Delim); !ok {
		return errMicrochipData.New()
	}

	for dec.More() {
		var jws microchipEntry
		if err := dec.Decode(&jws); err != nil {
			return errMicrochipData.WithCause(err)
		}
		kid := jws.Signatures[0].Header.KeyID
		key, ok := m.keys[kid]
		if !ok {
			return errMicrochipPublicKey.WithAttributes("id", kid)
		}
		buf, err := jws.Verify(key)
		if err != nil {
			return errMicrochipData.WithCause(err)
		}
		// publicKeySet.keys[0].x5c[0] contains the base64 certificate of the secure element.
		data := struct {
			PublicKeySet struct {
				Keys []struct {
					X5C []string `json:"x5c"`
				} `json:"keys"`
			} `json:"publicKeySet"`
		}{}
		if err := json.Unmarshal(buf, &data); err != nil {
			return errMicrochipData.WithCause(err)
		}
		if len(data.PublicKeySet.Keys) < 1 || len(data.PublicKeySet.Keys[0].X5C) < 1 {
			return errMicrochipData.WithCause(errMicrochipNoCertificate)
		}
		certBuf, err := base64.StdEncoding.DecodeString(data.PublicKeySet.Keys[0].X5C[0])
		if err != nil {
			return errMicrochipData.WithCause(err)
		}
		cert, err := x509.ParseCertificate(certBuf)
		if err != nil {
			return errMicrochipData.WithCause(err)
		}
		// In the B-generation, the DevEUI is the second element of a three-part CN.
		// In the C-generation, only the serial number is in the CN and the DevEUI is the SAN prepended with `eui64_`.
		cnParts := strings.SplitN(cert.Subject.CommonName, " ", 3)
		var devEUI types.EUI64
		switch len(cnParts) {
		case 1:
			for _, ext := range cert.Extensions {
				if !ext.Id.Equal(asn1.ObjectIdentifier{2, 5, 29, 17}) {
					continue
				}
				// The extension is a DirectoryName with a serial number which contains the illegal ASN.1 PrintableString character `_`:
				//  0  37: SEQUENCE {
				// 	2  35:   [4] {
				// 	4  33:     SEQUENCE {
				// 	6  31:       SET {
				// 	8  29:         SEQUENCE {
				// 10   3:           OBJECT IDENTIFIER '2 5 4 5'
				// 15  22:           PrintableString 'eui64_0004A310001AA90A'
				//       :             Error: PrintableString contains illegal character(s).
				//       :           }
				//       :         }
				//       :       }
				//       :     }
				//       :   }
				// Therefore, this needs to get binary replaced by HYPHEN-MINUS before parsing the content.
				// Note: do not remove the prefix here because the length of PrintableString needs to be preserved.
				buf := bytes.ReplaceAll(ext.Value, []byte(`eui64_`), []byte(`eui64-`))
				s := struct {
					DirectoryName []pkix.RDNSequence `asn1:"tag:4"`
				}{}
				if rest, err := asn1.Unmarshal(buf, &s); err != nil {
					return errMicrochipData.WithCause(errMicrochipCertificateSAN.WithCause(err))
				} else if len(rest) != 0 {
					return errMicrochipData.WithCause(errMicrochipCertificateSAN)
				} else if len(s.DirectoryName) == 0 {
					return errMicrochipData.WithCause(errMicrochipCertificateSAN)
				}
				var name pkix.Name
				name.FillFromRDNSequence(&s.DirectoryName[0])
				if err := devEUI.UnmarshalText([]byte(strings.TrimPrefix(name.SerialNumber, "eui64-"))); err != nil {
					return errMicrochipData.WithCause(errMicrochipCertificateSAN.WithCause(err))
				}
			}
		case 3:
			if err := devEUI.UnmarshalText([]byte(cnParts[1])); err != nil {
				return errMicrochipData.WithCause(err)
			}
		}
		if devEUI.IsZero() {
			return errMicrochipData.WithCause(errMicrochipUnknownDevEUI)
		}
		m := make(map[string]interface{})
		if err := json.Unmarshal(buf, &m); err != nil {
			return errMicrochipData.WithCause(err)
		}
		s, err := gogoproto.Struct(m)
		if err != nil {
			return errMicrochipData.WithCause(err)
		}
		sn := s.Fields["uniqueId"].GetStringValue()
		tmpl := &ttnpb.EndDeviceTemplate{
			EndDevice: ttnpb.EndDevice{
				EndDeviceIdentifiers: ttnpb.EndDeviceIdentifiers{
					DeviceID: strings.ToLower(fmt.Sprintf("eui-%s", devEUI)),
					JoinEUI:  &joinEUI,
					DevEUI:   &devEUI,
				},
				ProvisionerID:    provisioning.Microchip,
				ProvisioningData: s,
				RootKeys: &ttnpb.RootKeys{
					RootKeyID: sn,
				},
				SupportsJoin: true,
			},
			FieldMask: pbtypes.FieldMask{
				Paths: []string{
					"ids.device_id",
					"ids.dev_eui",
					"ids.join_eui",
					"provisioner_id",
					"provisioning_data",
					"root_keys.root_key_id",
					"supports_join",
				},
			},
			MappingKey: sn,
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- tmpl:
		}
	}
	return nil
}

func init() {
	getKeys := func(raw map[string][]byte) map[string]interface{} {
		keys := make(map[string]interface{}, len(raw))
		for kid, key := range raw {
			block, _ := pem.Decode(key)
			if block == nil {
				panic(fmt.Sprintf("invalid Microchip public key %v", kid))
			}
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				panic(err)
			}
			keys[kid] = cert.PublicKey
		}
		return keys
	}

	RegisterConverter("microchip-atecc608a-mahtn-t", &microchipATECC608AMAHTNT{
		keys: getKeys(microchipPublicKeys[microchipATECC608AMAHTNTPart]),
	})
	RegisterConverter("microchip-atecc608a-tnglora", &microchipATECC608TNGLORA{
		keys: getKeys(microchipPublicKeys[microchipATECC608TNGLORAPart]),
	})
}
