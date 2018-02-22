// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package store_test

import (
	"encoding/gob"
	"reflect"
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/errors"
	. "github.com/TheThingsNetwork/ttn/pkg/store"
	"github.com/gogo/protobuf/proto"
	"github.com/mitchellh/mapstructure"
)

type InterfaceStructA struct {
	A int
}

type InterfaceStructB struct {
	A map[string]bool
}

type SubSubStruct struct {
	ByteArray byteArray
	Bytes     []byte
	Empty     interface{}
	Int       int
	Interface interface{}
	String    string
}
type SubStruct struct {
	ByteArray    byteArray
	Bytes        []byte
	Empty        interface{}
	Int          int
	Interface    interface{}
	String       string
	SubSubStruct SubSubStruct
}
type Struct struct {
	ByteArray    byteArray
	Bytes        []byte
	Empty        interface{}
	Int          int
	Interface    interface{}
	String       string
	StructSlice  []SubStruct
	SubStruct    SubStruct
	SubStructPtr *SubStruct
}

type byteArray [2]byte

func mustToBytes(v interface{}) []byte {
	b, err := ToBytes(v)
	if err != nil {
		panic(err)
	}
	return b
}

func mustToBytesValue(v reflect.Value) []byte {
	b, err := ToBytesValue(v)
	if err != nil {
		panic(err)
	}
	return b
}

func wrapValue(v reflect.Value, t reflect.Type) reflect.Value {
	wv := reflect.New(t)
	wv.Elem().Set(v)
	return wv.Elem()
}

type ProtoMarshaler struct {
	a int
}

var _ proto.Marshaler = ProtoMarshaler{}
var _ proto.Unmarshaler = &ProtoMarshaler{}

func (m ProtoMarshaler) Marshal() ([]byte, error) {
	return []byte{byte(m.a), byte(ProtoEncoding)}, nil
}

func (m *ProtoMarshaler) Unmarshal(b []byte) error {
	if len(b) != 2 {
		return errors.Errorf("Encoded length must be 2, got %d", len(b))
	}
	if Encoding(b[1]) != ProtoEncoding {
		return errors.Errorf("Second byte must be %d, got %d", ProtoEncoding, b[1])
	}
	*m = ProtoMarshaler{
		a: int(b[0]),
	}
	return nil
}

var _ MapMarshaler = CustomMarshaler{}
var _ MapUnmarshaler = &CustomMarshaler{}
var _ ByteMapMarshaler = CustomMarshaler{}
var _ ByteMapUnmarshaler = &CustomMarshaler{}

type CustomMarshaler struct {
	a uint8
	b byte
	c []byte
}

func (cm CustomMarshaler) MarshalMap() (map[string]interface{}, error) {
	return map[string]interface{}{
		"aField": cm.a,
		"bField": cm.b,
		"cField": append(cm.c, 'X'),
	}, nil
}

func (cm CustomMarshaler) MarshalByteMap() (map[string][]byte, error) {
	return map[string][]byte{
		"aField": {cm.a},
		"bField": {cm.b},
		"cField": append(cm.c, 'X', 'X'),
	}, nil
}

func (cm *CustomMarshaler) UnmarshalMap(m map[string]interface{}) error {
	*cm = CustomMarshaler{
		a: m["aField"].(uint8),
		b: m["bField"].(byte),
		c: m["cField"].([]byte),
	}
	cm.c = cm.c[:len(cm.c)-1]
	return nil
}

func (cm *CustomMarshaler) UnmarshalByteMap(m map[string][]byte) error {
	*cm = CustomMarshaler{
		a: m["aField"][0],
		b: m["bField"][0],
		c: m["cField"],
	}
	cm.c = cm.c[:len(cm.c)-2]
	return nil
}

type CustomMarshalerAB struct {
	A *CustomMarshaler
	B *CustomMarshaler
}

func (cm CustomMarshalerAB) MarshalMap() (map[string]interface{}, error) {
	return map[string]interface{}{
		"A.aField": cm.A.a,
		"A.bField": cm.A.b,
		"A.cField": append(cm.A.c, 'X'),
		"B.aField": cm.B.a,
		"B.bField": cm.B.b,
		"B.cField": append(cm.B.c, 'X'),
	}, nil
}

func (cm CustomMarshalerAB) MarshalByteMap() (map[string][]byte, error) {
	return map[string][]byte{
		"A.aField": {cm.A.a},
		"A.bField": {cm.A.b},
		"A.cField": append(cm.A.c, 'X', 'X'),
		"B.aField": {cm.B.a},
		"B.bField": {cm.B.b},
		"B.cField": append(cm.B.c, 'X', 'X'),
	}, nil
}

func (cm *CustomMarshalerAB) UnmarshalMap(m map[string]interface{}) error {
	*cm = CustomMarshalerAB{
		A: &CustomMarshaler{
			a: m["A.aField"].(uint8),
			b: m["A.bField"].(byte),
			c: m["A.cField"].([]byte),
		},
		B: &CustomMarshaler{
			a: m["B.aField"].(uint8),
			b: m["B.bField"].(byte),
			c: m["B.cField"].([]byte),
		},
	}
	cm.A.c = cm.A.c[:len(cm.A.c)-1]
	cm.B.c = cm.B.c[:len(cm.B.c)-1]
	return nil
}

func (cm *CustomMarshalerAB) UnmarshalByteMap(m map[string][]byte) error {
	*cm = CustomMarshalerAB{
		A: &CustomMarshaler{
			a: m["A.aField"][0],
			b: m["A.bField"][0],
			c: m["A.cField"],
		},
		B: &CustomMarshaler{
			a: m["B.aField"][0],
			b: m["B.bField"][0],
			c: m["B.cField"],
		},
	}
	cm.A.c = cm.A.c[:len(cm.A.c)-2]
	cm.B.c = cm.B.c[:len(cm.B.c)-2]
	return nil
}

// Trick to register types before values is declared.
var _ interface{} = func() interface{} {
	for _, v := range []interface{}{
		InterfaceStructA{},
		InterfaceStructB{},
		time.Time{},
		struct{ A int }{},
	} {
		gob.Register(v)
	}
	return nil
}()

var values = []struct {
	unmarshaled interface{}
	marshaled   map[string]interface{}
	bytes       map[string][]byte
	decodeHooks []mapstructure.DecodeHookFunc
}{
	{
		Struct{
			String:    "42",
			Int:       42,
			Interface: InterfaceStructA{42},
			Bytes:     []byte("42"),
			ByteArray: byteArray([2]byte{'4', '2'}),
			StructSlice: []SubStruct{
				{
					String:    "42",
					ByteArray: byteArray([2]byte{'4', '2'}),
				},
				{
					Int:       42,
					Interface: float64(42),
					Bytes:     []byte("42"),
				},
				{
					Interface: InterfaceStructB{map[string]bool{"42": true}},
				},
			},
			SubStruct: SubStruct{
				String:    "42",
				Int:       42,
				Interface: "42",
				Bytes:     []byte("42"),
				ByteArray: byteArray([2]byte{'4', '2'}),
				SubSubStruct: SubSubStruct{
					String:    "42",
					Int:       42,
					Bytes:     []byte("42"),
					ByteArray: byteArray([2]byte{'4', '2'}),
				},
			},
			SubStructPtr: &SubStruct{
				String: "42",
				Int:    42,
			},
		},
		map[string]interface{}{
			"ByteArray": mustToBytes(byteArray{'4', '2'}),
			"Bytes":     []byte("42"),
			"Int":       int(42),
			"Interface": mustToBytesValue(
				wrapValue(reflect.ValueOf(InterfaceStructA{42}),
					reflect.TypeOf((*interface{})(nil)).Elem())),
			"String": string("42"),
			"StructSlice": mustToBytes([]SubStruct{
				{
					String:    "42",
					ByteArray: byteArray([2]byte{'4', '2'}),
				},
				{
					Int:       42,
					Interface: float64(42),
					Bytes:     []byte("42"),
				},
				{
					Interface: InterfaceStructB{map[string]bool{"42": true}},
				},
			}),
			"SubStruct.ByteArray": mustToBytes(byteArray{'4', '2'}),
			"SubStruct.Bytes":     []byte("42"),
			"SubStruct.Int":       int(42),
			"SubStruct.Interface": mustToBytesValue(
				wrapValue(reflect.ValueOf("42"),
					reflect.TypeOf((*interface{})(nil)).Elem())),
			"SubStruct.String":                 string("42"),
			"SubStruct.SubSubStruct.ByteArray": mustToBytes(byteArray{'4', '2'}),
			"SubStruct.SubSubStruct.Bytes":     []byte("42"),
			"SubStruct.SubSubStruct.Int":       int(42),
			"SubStruct.SubSubStruct.String":    string("42"),
			"SubStructPtr.Int":                 int(42),
			"SubStructPtr.String":              string("42"),
		},

		map[string][]byte{
			"ByteArray": mustToBytes("42"),
			"Bytes":     mustToBytes("42"),
			"Int":       mustToBytes(42),
			"Interface": mustToBytesValue(
				wrapValue(reflect.ValueOf(InterfaceStructA{42}),
					reflect.TypeOf((*interface{})(nil)).Elem())),
			"String": mustToBytes("42"),
			"StructSlice": mustToBytes([]SubStruct{
				{
					String:    "42",
					ByteArray: byteArray([2]byte{'4', '2'}),
				},
				{
					Int:       42,
					Interface: float64(42),
					Bytes:     []byte("42"),
				},
				{
					Interface: InterfaceStructB{map[string]bool{"42": true}},
				},
			}),
			"SubStruct.ByteArray": mustToBytes("42"),
			"SubStruct.Bytes":     mustToBytes("42"),
			"SubStruct.Int":       mustToBytes(42),
			"SubStruct.Interface": mustToBytesValue(
				wrapValue(reflect.ValueOf("42"),
					reflect.TypeOf((*interface{})(nil)).Elem())),
			"SubStruct.String":                 mustToBytes("42"),
			"SubStruct.SubSubStruct.ByteArray": mustToBytes("42"),
			"SubStruct.SubSubStruct.Bytes":     mustToBytes("42"),
			"SubStruct.SubSubStruct.Int":       mustToBytes(42),
			"SubStruct.SubSubStruct.String":    mustToBytes("42"),
			"SubStructPtr.Int":                 mustToBytes(42),
			"SubStructPtr.String":              mustToBytes("42"),
		},
		nil,
	},
	{
		struct {
			a int
			b int
		}{},
		(map[string]interface{})(nil),
		(map[string][]byte)(nil),
		nil,
	},
	{
		struct {
			A interface{}
			B interface{}
		}{
			42,
			struct{ A int }{42},
		},
		map[string]interface{}{
			"A": mustToBytesValue(
				wrapValue(reflect.ValueOf(42),
					reflect.TypeOf((*interface{})(nil)).Elem())),
			"B": mustToBytesValue(
				wrapValue(reflect.ValueOf(struct{ A int }{42}),
					reflect.TypeOf((*interface{})(nil)).Elem())),
		},
		map[string][]byte{
			"A": mustToBytesValue(
				wrapValue(reflect.ValueOf(42),
					reflect.TypeOf((*interface{})(nil)).Elem())),
			"B": mustToBytesValue(
				wrapValue(reflect.ValueOf(struct{ A int }{42}),
					reflect.TypeOf((*interface{})(nil)).Elem())),
		},
		nil,
	},
	{
		struct{ time.Time }{time.Unix(42, 42).UTC()},
		map[string]interface{}{"Time": mustToBytes(time.Unix(42, 42).UTC())},
		map[string][]byte{"Time": mustToBytes(time.Unix(42, 42).UTC())},
		nil,
	},
	{
		struct{ T time.Time }{time.Unix(42, 42).UTC()},
		map[string]interface{}{"T": mustToBytes(time.Unix(42, 42).UTC())},
		map[string][]byte{"T": mustToBytes(time.Unix(42, 42).UTC())},
		nil,
	},
	{
		struct{ Interfaces []interface{} }{[]interface{}{
			nil,
			nil,
			nil,
			time.Time{},
			&time.Time{},
			struct{ A int }{42},
		}},
		map[string]interface{}{
			"Interfaces": mustToBytes([]interface{}{
				nil,
				nil,
				nil,
				time.Time{},
				&time.Time{},
				struct{ A int }{42},
			}),
		},
		map[string][]byte{
			"Interfaces": mustToBytes([]interface{}{
				nil,
				nil,
				nil,
				time.Time{},
				&time.Time{},
				struct{ A int }{42},
			}),
		},
		nil,
	},
	{
		struct {
			A *ProtoMarshaler
		}{
			&ProtoMarshaler{42},
		},
		map[string]interface{}{
			"A": mustToBytes(ProtoMarshaler{42}),
		},
		map[string][]byte{
			"A": mustToBytes(ProtoMarshaler{42}),
		},
		nil,
	},
	{
		CustomMarshaler{
			a: 42,
			b: 43,
			c: []byte("foo"),
		},
		map[string]interface{}{
			"aField": uint8(42),
			"bField": byte(43),
			"cField": []byte("fooX"),
		},
		map[string][]byte{
			"aField": {42},
			"bField": {43},
			"cField": []byte("fooXX"),
		},
		nil,
	},
	{
		CustomMarshalerAB{
			&CustomMarshaler{
				a: 42,
				b: 43,
				c: []byte("foo"),
			},
			&CustomMarshaler{
				a: 4,
				b: 5,
				c: []byte("bar"),
			},
		},
		map[string]interface{}{
			"A.aField": uint8(42),
			"A.bField": byte(43),
			"A.cField": []byte("fooX"),
			"B.aField": uint8(4),
			"B.bField": byte(5),
			"B.cField": []byte("barX"),
		},
		map[string][]byte{
			"A.aField": {42},
			"A.bField": {43},
			"A.cField": []byte("fooXX"),
			"B.aField": {4},
			"B.bField": {5},
			"B.cField": []byte("barXX"),
		},
		nil,
	},
}
