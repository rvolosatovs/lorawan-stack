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

// Package randutil provides pseudo-random number generator utilities.
package randutil

import (
	"math/rand"
	"sync"
)

// LockedSource is a rand.Source, which is safe for concurrent use. Adapted from the non-exported
// lockedSource from stdlib rand.
type LockedSource struct {
	mu  sync.Mutex
	src rand.Source
	s64 rand.Source64 // non-nil if src is source64
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (r *LockedSource) Int63() (n int64) {
	r.mu.Lock()
	n = r.src.Int63()
	r.mu.Unlock()
	return
}

// Uint64 returns a pseudo-random 64-bit value as a uint64.
func (r *LockedSource) Uint64() (n uint64) {
	if r.s64 != nil {
		r.mu.Lock()
		n = r.s64.Uint64()
		r.mu.Unlock()
		return
	}
	return uint64(r.Int63())>>31 | uint64(r.Int63())<<32
}

// Seed uses the provided seed value to initialize the generator to a deterministic state.
// Seed should not be called concurrently with any other Rand method.
func (r *LockedSource) Seed(seed int64) {
	r.mu.Lock()
	r.src.Seed(seed)
	r.mu.Unlock()
}

// NewLockedSource returns a rand.Source, which is safe for concurrent use.
func NewLockedSource(src rand.Source) *LockedSource {
	s64, _ := src.(rand.Source64)
	return &LockedSource{
		src: src,
		s64: s64,
	}
}
