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

package marshaling

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"reflect"

	"github.com/gogo/protobuf/proto"
	"go.thethings.network/lorawan-stack/pkg/errors"
)

var (
	reflectValueType = reflect.TypeOf(reflect.Value{})

	jsonMarshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	mapMarshalerType  = reflect.TypeOf((*MapMarshaler)(nil)).Elem()
)

// FlattenedValue is like Flattened, but it operates on maps containing reflect.Value.
func FlattenedValue(m map[string]reflect.Value) (out map[string]reflect.Value) {
	out = make(map[string]reflect.Value, len(m))
	for k, v := range m {
		if !v.IsValid() {
			out[k] = v
			continue
		}

		vt := v.Type()
		if vt.Kind() == reflect.Map && vt.Key().Kind() == reflect.String && vt.Elem() == reflectValueType {
			for sk, sv := range FlattenedValue(v.Interface().(map[string]reflect.Value)) {
				out[k+Separator+sk] = sv
			}
		} else {
			out[k] = v
		}
	}
	return out
}

// Flattened returns a copy of m with keys 'Flattened'.
// If the map contains sub-maps, the values of these sub-maps are set under the root map, each level separated by Separator.
func Flattened(m map[string]interface{}) (out map[string]interface{}) {
	out = make(map[string]interface{}, len(m))
	for k, v := range m {
		if sm, ok := v.(map[string]interface{}); ok {
			sm = Flattened(sm)
			for sk, sv := range sm {
				out[k+Separator+sk] = sv
			}
		} else {
			out[k] = v
		}
	}
	return out
}

// ToValueMap converts the input map m into a map[string]reflect.Value by calling reflect.ValueOf with each value in m.
func ToValueMap(m map[string]interface{}) map[string]reflect.Value {
	vm := make(map[string]reflect.Value, len(m))
	for k, iv := range m {
		vm[k] = reflect.ValueOf(iv)
	}
	return vm
}

// marshalNested retrieves recursively all types for the given value
// and returns the marshaled nested value.
func marshalNested(v reflect.Value) (reflect.Value, error) {
	t := v.Type()
	if t.Implements(mapMarshalerType) {
		if !reflect.Indirect(v).IsValid() {
			return reflect.Value{}, nil
		}

		im, err := v.Interface().(MapMarshaler).MarshalMap()
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(ToValueMap(im)), nil
	}

	if IsNillableKind(t.Kind()) && v.IsNil() {
		return v, nil
	}

	iv := reflect.Indirect(v)
	if iv.Kind() != reflect.Struct {
		return v, nil
	}

	t = iv.Type()
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).PkgPath == "" {
			// Only attempt to marshal structs with exported fields
			m, err := marshal(v)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(m), nil
		}
	}
	return v, nil
}

var errExpectedStruct = errors.DefineCorruption("expected_struct", "expected argument to be a struct, got `{result}` (kind: `{kind}`)")

// marhshal converts the given struct s to a map[string]reflect.Value
func marshal(v reflect.Value) (m map[string]reflect.Value, err error) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return nil, errExpectedStruct.WithAttributes("result", t.String(), "kind", t.Kind().String())
	}

	m = make(map[string]reflect.Value, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}

		m[f.Name], err = marshalNested(v.FieldByName(f.Name))
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

// MapMarshaler is the interface implemented by an object that can
// marshal itself into a flattened map[string]interface{}
//
// MarshalMap encodes the receiver into map[string]interface{} and returns the result.
type MapMarshaler interface {
	MarshalMap() (map[string]interface{}, error)
}

// MarshalMap returns the map encoding of v, where v is either a struct or a map with string keys.
//
// MarshalMap traverses the value v recursively. If v implements the MapMarshaler interface, MarshalMap calls its MarshalMap method to produce map[string]interface{}.
// Otherwise, MarshalMap first encodes the value v as a map[string]interface{}. Default marshaler marshals slices as maps with string keys, where all keys represent integers.
// The map produced by any of the methods will be flattened by joining sub-map values with a dot(note that slices produced by custom MarshalMap implementations won't be flattened).
func MarshalMap(v interface{}) (m map[string]interface{}, err error) {
	if mm, ok := v.(MapMarshaler); ok {
		return mm.MarshalMap()
	}

	vm, err := marshal(reflect.ValueOf(v))
	if err != nil {
		return nil, err
	}
	if len(vm) == 0 {
		return nil, nil
	}

	vm = FlattenedValue(vm)

	m = make(map[string]interface{}, len(vm))
	for k, v := range vm {
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64, reflect.Bool, reflect.String:
			m[k] = v.Interface()

		default:
			bv, err := ToBytesValue(v)
			if err != nil {
				return nil, err
			}
			m[k] = bv
		}
	}
	return m, nil
}

// binaryEncodable returns v as value, which can be encoded using binary package functionality.
// binaryEncodable returns nil, false if conversion is not possible.
func binaryEncodable(v reflect.Value) (interface{}, bool) {
	if !v.IsValid() {
		return nil, false
	}

	switch v.Kind() {
	case reflect.Int:
		return v.Int(), true

	case reflect.Uint, reflect.Uintptr:
		return v.Uint(), true

	case reflect.Bool,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return v.Interface(), true

	case reflect.String:
		return []byte(v.String()), true

	case reflect.Slice, reflect.Array:
		switch v.Type().Elem().Kind() {
		case reflect.Bool,
			reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
			return v.Interface(), true
		}
	}
	return nil, false
}

var errNoCustomDecoder = errors.DefineInternal("no_custom_decoder", "values of type `{type}` (kind `{kind}`), do not implement custom decoder")

// ToBytesValue is like ToBytes, but operates on values of type reflect.Value.
func ToBytesValue(v reflect.Value) (b []byte, err error) {
	var enc Encoding
	defer func() {
		if err != nil {
			return
		}
		b = append([]byte{byte(DefaultVersion), byte(enc)}, b...)
	}()

	if !v.IsValid() ||
		(IsNillableKind(v.Kind()) && v.IsNil()) ||
		reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface()) {
		return nil, nil
	}

	t := v.Type()

	conv, ok := binaryEncodable(v)
	if ok {
		enc = BigEndianEncoding

		if conv == nil {
			return nil, nil
		}

		buf := &bytes.Buffer{}
		err := binary.Write(buf, binary.BigEndian, conv)
		return buf.Bytes(), err
	}

	switch {
	case t.Implements(jsonMarshalerType):
		enc = JSONEncoding
		return json.Marshal(v.Interface())

	case reflect.PtrTo(t).Implements(jsonMarshalerType):
		enc = JSONEncoding
		ptr := reflect.New(t)
		ptr.Elem().Set(v)
		return json.Marshal(ptr.Interface().(proto.Message))

	case t.Implements(protoMessageType):
		enc = ProtoEncoding
		return proto.Marshal(v.Interface().(proto.Message))

	case reflect.PtrTo(t).Implements(protoMessageType):
		enc = ProtoEncoding
		ptr := reflect.New(t)
		ptr.Elem().Set(v)
		return proto.Marshal(ptr.Interface().(proto.Message))

	case t.Kind() == reflect.Chan, t.Kind() == reflect.Func:
		return nil, errNoCustomDecoder.WithAttributes("type", t.String(), "kind", t.Kind().String())
	}

	enc = GobEncoding

	// Encode the value as a pointer to include type info.
	pv := reflect.New(t)
	pv.Elem().Set(v)

	buf := &bytes.Buffer{}
	if err := gob.NewEncoder(buf).EncodeValue(pv); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ToBytes marshals v into a []byte value and returns the result.
// Slices and arrays of bytes, strings, booleans and numeric types are stored in a human-readable
// format, if value implements proto.Marshaler, result of Marshal() method is stored, otherwise encoding/gob is used.
// Encoded values have the according Encoding byte prepended.
func ToBytes(v interface{}) (b []byte, err error) {
	return ToBytesValue(reflect.ValueOf(v))
}

// ByteMapMarshaler is the interface implemented by an object that can
// marshal itself into a map[string][]byte.
//
// MarshalByteMap encodes the receiver into map[string][]byte and returns the result.
type ByteMapMarshaler interface {
	MarshalByteMap() (map[string][]byte, error)
}

// MarshalByteMap returns the byte map encoding of v.
//
// MarshalByteMap traverses map returned by Marshal and converts all values to bytes.
func MarshalByteMap(v interface{}) (bm map[string][]byte, err error) {
	if bmm, ok := v.(ByteMapMarshaler); ok {
		return bmm.MarshalByteMap()
	}

	var vm map[string]reflect.Value
	if mm, ok := v.(MapMarshaler); ok {
		im, err := mm.MarshalMap()
		if err != nil {
			return nil, err
		}
		vm = ToValueMap(im)
	} else {
		vm, err = marshal(reflect.ValueOf(v))
		if err != nil {
			return nil, err
		}
		vm = FlattenedValue(vm)
	}
	if len(vm) == 0 {
		return nil, nil
	}

	bm = make(map[string][]byte, len(vm))
	for k, v := range vm {
		b, err := ToBytesValue(v)
		if err != nil {
			return nil, err
		}
		bm[k] = b
	}
	return bm, nil
}
