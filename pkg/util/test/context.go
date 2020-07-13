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

package test

import (
	"context"
	"testing"
)

type tbKey struct{}

// ContextWithTB saves the testing.TB in the context.
func ContextWithTB(ctx context.Context, tb testing.TB) context.Context {
	return context.WithValue(ctx, tbKey{}, tb)
}

// TBFromContext returns the testing.TB saved using ContextWithTB from the context.
func TBFromContext(ctx context.Context) (testing.TB, bool) {
	tb, ok := ctx.Value(tbKey{}).(testing.TB)
	if !ok {
		return nil, false
	}
	return tb, true
}

// MustTBFromContext returns the testing.TB from the context, and panics if it was not saved in the context.
func MustTBFromContext(ctx context.Context) testing.TB {
	tb, ok := TBFromContext(ctx)
	if !ok {
		panic("testing.TB not present in the context")
	}
	return tb
}

// TFromContext returns the *testing.T saved using ContextWithTB from the context.
func TFromContext(ctx context.Context) (*testing.T, bool) {
	t, ok := ctx.Value(tbKey{}).(*testing.T)
	if !ok {
		return nil, false
	}
	return t, true
}

// MustTFromContext returns the *testing.T from the context, and panics if it was not saved in the context.
func MustTFromContext(ctx context.Context) *testing.T {
	t, ok := TFromContext(ctx)
	if !ok {
		panic("*testing.T not present in the context")
	}
	return t
}
