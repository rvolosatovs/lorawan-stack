// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package store

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/TheThingsNetwork/ttn/pkg/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/mitchellh/mapstructure"
)

var (
	mapUnmarshalerType     = reflect.TypeOf((*MapUnmarshaler)(nil)).Elem()
	protoUnmarshalerType   = reflect.TypeOf((*proto.Unmarshaler)(nil)).Elem()
	jsonUnmarshalerType    = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	byteMapUnmarshalerType = reflect.TypeOf((*ByteMapUnmarshaler)(nil)).Elem()
)

// Unflattened unflattens m and returns the result
func Unflattened(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		skeys := strings.Split(k, Separator)
		parent := out
		for _, sk := range skeys[:len(skeys)-1] {
			sm, ok := parent[sk]
			if !ok {
				sm = make(map[string]interface{})
				parent[sk] = sm
			}
			parent = sm.(map[string]interface{})
		}
		parent[skeys[len(skeys)-1]] = v
	}
	return out
}

func isByteSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8
}

// MapUnmarshaler is the interface implemented by an object that can
// unmarshal a map[string]interface{} representation of itself.
//
// UnmarshalMap must be able to decode the form generated by MarshalMap.
// UnmarshalMap must deep copy the data if it wishes to retain the data after returning.
type MapUnmarshaler interface {
	UnmarshalMap(map[string]interface{}) error
}

func unmarshalMapHookFunc(f reflect.Type, t reflect.Type, v interface{}) (interface{}, error) {
	switch {
	case isByteSlice(f) && !isByteSlice(t):
		return BytesToType(v.([]byte), t)

	case t.Kind() == reflect.Interface,
		t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Interface,
		f.Kind() != reflect.Map:
		return v, nil

	case t.Implements(mapUnmarshalerType):
		out := reflect.New(t.Elem()).Interface().(MapUnmarshaler)
		return out, out.UnmarshalMap(v.(map[string]interface{}))

	case reflect.PtrTo(t).Implements(mapUnmarshalerType):
		out := reflect.New(t).Interface().(MapUnmarshaler)
		return out, out.UnmarshalMap(v.(map[string]interface{}))

	default:
		return v, nil
	}
}

// UnmarshalMap parses the map-encoded data and stores the result
// in the value pointed to by v.
//
// UnmarshalMap uses the inverse of the encodings that
// Marshal uses.
func UnmarshalMap(m map[string]interface{}, v interface{}, hooks ...mapstructure.DecodeHookFunc) error {
	if mu, ok := v.(MapUnmarshaler); ok {
		return mu.UnmarshalMap(m)
	}

	if len(m) == 0 {
		return nil
	}

	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		ZeroFields:       true,
		Result:           v,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(append(
			hooks,
			unmarshalMapHookFunc,
		)...),
	})
	if err != nil {
		panic(errors.NewWithCause(err, "Failed to intialize decoder"))
	}
	return dec.Decode(Unflattened(m))
}

// BytesToType decodes []byte value in b into a new value of type typ.
func BytesToType(b []byte, typ reflect.Type) (interface{}, error) {
	if len(b) == 0 {
		return nil, ErrInvalidData
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	enc := Encoding(b[0])
	b = b[1:]

	switch enc {
	case RawEncoding:
		if len(b) == 0 {
			return reflect.New(typ).Elem().Interface(), nil
		}

		switch k := typ.Kind(); k {
		case reflect.String:
			return string(b), nil
		case reflect.Bool:
			return strconv.ParseBool(string(b))
		case reflect.Int:
			return strconv.ParseInt(string(b), 10, 64)
		case reflect.Int8:
			return strconv.ParseInt(string(b), 10, 8)
		case reflect.Int16:
			return strconv.ParseInt(string(b), 10, 16)
		case reflect.Int32:
			return strconv.ParseInt(string(b), 10, 32)
		case reflect.Int64:
			return strconv.ParseInt(string(b), 10, 64)
		case reflect.Uint:
			return strconv.ParseUint(string(b), 10, 64)
		case reflect.Uint8:
			return strconv.ParseUint(string(b), 10, 8)
		case reflect.Uint16:
			return strconv.ParseUint(string(b), 10, 16)
		case reflect.Uint32:
			return strconv.ParseUint(string(b), 10, 32)
		case reflect.Uint64:
			return strconv.ParseUint(string(b), 10, 64)
		case reflect.Float32:
			return strconv.ParseFloat(string(b), 32)
		case reflect.Float64:
			return strconv.ParseFloat(string(b), 64)
		case reflect.Slice, reflect.Array:
			elem := typ.Elem()
			if elem.Kind() == reflect.Uint8 {
				// Handle byte slices/arrays directly
				if k == reflect.Slice {
					return b, nil
				}
				rv := reflect.Indirect(reflect.New(typ))
				for i := 0; i < rv.Len(); i++ {
					rv.Index(i).SetUint(uint64(b[i]))
				}
				return rv.Interface(), nil
			}
		}
		return nil, errors.Errorf("Can not decode raw bytes to value of type %s", typ)
	case JSONEncoding:
		if !reflect.PtrTo(typ).Implements(jsonUnmarshalerType) {
			return nil, errors.Errorf("Expected %s to implement %s", typ, jsonUnmarshalerType)
		}
		v := reflect.New(typ).Interface().(json.Unmarshaler)
		return v, v.UnmarshalJSON(b)
	case ProtoEncoding:
		if !reflect.PtrTo(typ).Implements(protoUnmarshalerType) {
			return nil, errors.Errorf("Expected %s to implement %s", typ, protoUnmarshalerType)
		}
		v := reflect.New(typ).Interface().(proto.Unmarshaler)
		return v, v.Unmarshal(b)
	case GobEncoding:
		v := reflect.New(typ)
		if err := gob.NewDecoder(bytes.NewReader(b)).DecodeValue(v); err != nil {
			return nil, err
		}
		return v.Elem().Interface(), nil
	default:
		return nil, ErrInvalidData
	}
}

// ByteMapUnmarshaler is the interface implemented by an object that can
// unmarshal a map[string][]byte representation of itself.
//
// UnmarshalByteMap must be able to decode the form generated by MarshalByteMap.
// UnmarshalByteMap must deep copy the data if it wishes to retain the data after returning.
type ByteMapUnmarshaler interface {
	UnmarshalByteMap(map[string][]byte) error
}

func interfaceMapToByteMap(im map[string]interface{}) (bm map[string][]byte, ok bool) {
	bm = make(map[string][]byte, len(im))
	for k, v := range im {
		bm[k], ok = v.([]byte)
		if !ok {
			return nil, false
		}
	}
	return
}

func unmarshalByteMapHookFunc(f reflect.Type, t reflect.Type, v interface{}) (interface{}, error) {
	switch {
	case isByteSlice(f):
		return BytesToType(v.([]byte), t)

	case t.Kind() == reflect.Interface,
		t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Interface,
		f.Kind() != reflect.Map:
		return v, nil

	case t.Implements(byteMapUnmarshalerType):
		out := reflect.New(t.Elem()).Interface().(ByteMapUnmarshaler)
		bm, ok := interfaceMapToByteMap(v.(map[string]interface{}))
		if !ok {
			return nil, errors.New("Mismatched hook")
		}
		return out, out.UnmarshalByteMap(bm)

	case reflect.PtrTo(t).Implements(byteMapUnmarshalerType):
		out := reflect.New(t).Interface().(ByteMapUnmarshaler)
		bm, ok := interfaceMapToByteMap(v.(map[string]interface{}))
		if !ok {
			return nil, errors.New("Mismatched hook")
		}
		return out, out.UnmarshalByteMap(bm)

	default:
		return v, nil
	}
}

// UnmarshalByteMap parses the byte map-encoded data and stores the result in the value pointed to by v.
// UnmarshalByteMap uses the inverse of the encodings that MarshalByteMap uses.
func UnmarshalByteMap(bm map[string][]byte, v interface{}, hooks ...mapstructure.DecodeHookFunc) error {
	if mm, ok := v.(ByteMapUnmarshaler); ok {
		return mm.UnmarshalByteMap(bm)
	}

	if len(bm) == 0 {
		return nil
	}

	im := make(map[string]interface{}, len(bm))
	for k, v := range bm {
		im[k] = v
	}

	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		ZeroFields:       true,
		Result:           v,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(append(
			hooks,
			unmarshalByteMapHookFunc,
		)...),
	})
	if err != nil {
		panic(errors.NewWithCause(err, "Failed to intialize decoder"))
	}
	return dec.Decode(Unflattened(im))
}
