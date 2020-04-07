// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
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
	"fmt"
	"reflect"

	"github.com/kr/pretty"
)

// StringEqual returns true iff strings x and y are equal and false otherwise.
func StringEqual(x, y string) bool {
	return x == y
}

// DiffEqual returns true iff pretty.Diff of x and y is empty and false otherwise.
func DiffEqual(x, y interface{}) bool {
	return len(pretty.Diff(x, y)) == 0
}

// Ranger represents an entity, which can be ranged over(e.g. sync.Map).
type Ranger interface {
	Range(f func(k, v interface{}) bool)
}

type indexRanger struct {
	reflect.Value
}

func (rv indexRanger) Range(f func(k, v interface{}) bool) {
	for i := 0; i < rv.Len(); i++ {
		if !f(i, rv.Index(i).Interface()) {
			return
		}
	}
}

type mapRanger struct {
	reflect.Value
}

func (rv mapRanger) Range(f func(k, v interface{}) bool) {
	for _, k := range rv.MapKeys() {
		if !f(k.Interface(), rv.MapIndex(k).Interface()) {
			return
		}
	}
}

// WrapRanger returns Ranger, true if v can be ranged over and nil, false otherwise.
func WrapRanger(v interface{}) (Ranger, bool) {
	r, ok := v.(Ranger)
	if ok {
		return r, ok
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String, reflect.Slice, reflect.Array:
		return indexRanger{rv}, true
	case reflect.Map:
		return mapRanger{rv}, true
	}
	return nil, false
}

// MapKeys returns the keys of a map as a slice.
func MapKeys(m interface{}) []interface{} {
	if m == nil {
		return nil
	}
	rv := reflect.ValueOf(m)
	if rv.Kind() != reflect.Map {
		panic(fmt.Errorf("expected %T to be of map type", m))
	}
	ks := make([]interface{}, 0, rv.Len())
	for _, k := range rv.MapKeys() {
		ks = append(ks, k.Interface())
	}
	return ks
}

// IsSubsetOfElements returns true iff a multiset sub represents a subset of
// multiset super under equality given by eq.
// Signature of eq must be func(A, B) bool, where A, B are types, which
// elements of sub and super can be assigned to respectively.
// It panics if either sub or super is not one of:
// 1. string, slice, array or map kind
// 2. value, which implements Ranger interface(e.g. sync.Map)
// NOTE: Map key values are not taken into account.
func IsSubsetOfElements(eq interface{}, sub, super interface{}) bool {
	if sub == nil {
		// NOTE: Empty set is a subset of any set.
		return true
	}
	if super == nil {
		// NOTE: No non-empty set is a subset of empty set.
		return false
	}

	ev := reflect.ValueOf(eq)
	if ev.Kind() != reflect.Func {
		panic(fmt.Errorf("expected kind of eq to be a function, got: %s", ev.Kind()))
	}
	subR, ok := WrapRanger(sub)
	if !ok {
		panic(fmt.Errorf("cannot range over values of type %T", sub))
	}
	supR, ok := WrapRanger(super)
	if !ok {
		panic(fmt.Errorf("cannot range over values of type %T", super))
	}

	type entry struct {
		value reflect.Value
		found uint
	}
	entries := map[*entry]struct{}{}

	findEntry := func(v reflect.Value) *entry {
		for e := range entries {
			if ev.Call([]reflect.Value{e.value, v})[0].Bool() {
				return e
			}
		}
		return nil
	}

	subR.Range(func(_, v interface{}) bool {
		rv := reflect.ValueOf(v)
		e := findEntry(rv)
		if e == nil {
			entries[&entry{
				value: rv,
				found: 1,
			}] = struct{}{}
		} else {
			e.found++
		}
		return true
	})
	supR.Range(func(_, v interface{}) bool {
		rv := reflect.ValueOf(v)
		e := findEntry(rv)
		if e == nil {
			return true
		}
		if e.found == 1 {
			delete(entries, e)
		} else {
			e.found--
		}
		return true
	})
	return len(entries) == 0
}

// IsProperSubsetOfElements is like IsSubsetOfElements, but checks for proper subset.
func IsProperSubsetOfElements(eq interface{}, sub, super interface{}) bool {
	return IsSubsetOfElements(eq, sub, super) && !IsSubsetOfElements(eq, super, sub)
}

// SameElements returns true iff IsSubsetOfElements(eq, xs, ys) returns true and IsSubsetOfElements(eq, ys, xs) returns true and false otherwise.
func SameElements(eq interface{}, xs, ys interface{}) bool {
	return IsSubsetOfElements(eq, xs, ys) && IsSubsetOfElements(eq, ys, xs)
}
