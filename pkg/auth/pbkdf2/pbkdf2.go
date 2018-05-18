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

// Package pbkdf2 implements the PBKDF2 algorithm method used to hash passwords.
package pbkdf2

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/random"
	"golang.org/x/crypto/pbkdf2"
)

var defaultInstance = &PBKDF2{
	iterations: 20000,
	keyLength:  512,
	algorithm:  Sha512,
	saltLength: 64,
}

// SetDefaultIterations sets the number of iterations for the default PBKDF2 instance.
// This should typically only be changed for testing purposes.
func SetDefaultIterations(iterations int) {
	defaultInstance.iterations = iterations
}

// Default returns a the default PBKDF2 instance.
func Default() *PBKDF2 { return defaultInstance }

// PBKDF2 is a password derivation method.
type PBKDF2 struct {
	// Iterations is the number of iterations to use in the PBKDF2 algorithm.
	iterations int

	// Algorithm is the hashing algorithm used.
	algorithm Algorithm

	// SaltLength is the length of the salt used.
	saltLength int

	// KeyLength is the length of the desired key.
	keyLength int
}

// Name returns the name of the PBKDF2 hashing method.
func (*PBKDF2) Name() string {
	return "PBKDF2"
}

// Hash hashes a plain text password.
func (p *PBKDF2) Hash(plain string) (string, error) {
	if p.saltLength == 0 {
		return "", errors.Errorf("Salts can not have zero length")
	}

	salt := random.String(p.saltLength)
	hash := hash64([]byte(plain), []byte(salt), p.iterations, p.keyLength, p.algorithm)
	pass := fmt.Sprintf("%s$%s$%v$%s$%s", p.Name(), p.algorithm, p.iterations, salt, string(hash))

	return pass, nil
}

// Validate validates a plaintext password against a hashed one.
// The format of the hashed password should be:
//
//     PBKDF2$<algorithm>$<iterations>$<salt>$<key in base64>
//
func (*PBKDF2) Validate(hashed, plain string) (bool, error) {
	parts := strings.Split(hashed, "$")
	if len(parts) != 5 {
		return false, errors.Errorf("Invalid PBKDF2 format")
	}

	alg := parts[1]
	algorithm, err := parseAlgorithm(alg)
	if err != nil {
		return false, err
	}

	iter, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		return false, errors.Errorf("Invalid number of iterations: %s", parts[2])
	}
	salt := parts[3]
	key := parts[4]

	// Get the key length.
	keylen, err := keyLen(key)
	if err != nil {
		return false, errors.Errorf("Could not get key length: %s", err)
	}

	// Hash the plaintext.
	hash := hash64([]byte(plain), []byte(salt), int(iter), keylen, algorithm)

	// Compare the hashed plaintext and the stored hash.
	return subtle.ConstantTimeCompare(hash, []byte(key)) == 1, nil
}

// hash64 hashes a plain password and encodes it to base64.
func hash64(plain, salt []byte, iterations int, keyLength int, algorithm Algorithm) []byte {
	key := pbkdf2.Key(plain, salt, iterations, keyLength, algorithm.Hash)
	res := make([]byte, base64.RawURLEncoding.EncodedLen(len(key)))
	base64.RawURLEncoding.Encode(res, key)
	return res
}

// get the key length from the key.
func keyLen(key string) (int, error) {
	buf, err := base64.RawURLEncoding.DecodeString(key)
	if err != nil {
		return 0, err
	}
	return len(buf), nil
}
