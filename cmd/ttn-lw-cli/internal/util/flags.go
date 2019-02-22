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

package util

import (
	"encoding"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/spf13/pflag"
	"go.thethings.network/lorawan-stack/pkg/errors"
)

var (
	toDash       = strings.NewReplacer("_", "-")
	toUnderscore = strings.NewReplacer("-", "_")
)

func NormalizeFlags(f *pflag.FlagSet, name string) pflag.NormalizedName {
	return pflag.NormalizedName(toDash.Replace(name))
}

func SelectFieldMask(cmdFlags *pflag.FlagSet, fieldMaskFlags ...*pflag.FlagSet) (paths []string) {
	cmdFlags.Visit(func(flag *pflag.Flag) {
		flagName := toUnderscore.Replace(flag.Name)
		for _, fieldMaskFlags := range fieldMaskFlags {
			if b, err := fieldMaskFlags.GetBool(flagName); err == nil && b {
				paths = append(paths, flagName)
				return
			}
		}
	})
	return
}

func UpdateFieldMask(cmdFlags *pflag.FlagSet, fieldMaskFlags ...*pflag.FlagSet) (paths []string) {
	cmdFlags.Visit(func(flag *pflag.Flag) {
		flagName := toUnderscore.Replace(flag.Name)
		for _, fieldMaskFlags := range fieldMaskFlags {
			if fieldMaskFlags.Lookup(flagName) != nil {
				paths = append(paths, flagName)
				return
			}
		}
	})
	return
}

func FieldFlags(v interface{}, prefix ...string) *pflag.FlagSet {
	t := reflect.Indirect(reflect.ValueOf(v)).Type()
	return fieldMaskFlags(prefix, t, false)
}

func FieldMaskFlags(v interface{}, prefix ...string) *pflag.FlagSet {
	t := reflect.Indirect(reflect.ValueOf(v)).Type()
	return fieldMaskFlags(prefix, t, true)
}

func isAtomicType(t reflect.Type, maskOnly bool) bool {
	switch t.PkgPath() {
	case "time":
		switch t.Name() {
		case "Time", "Duration":
			return true
		}
	case "github.com/gogo/protobuf/types":
		switch t.Name() {
		case "DoubleValue", "FloatValue", "Int64Value", "UInt64Value", "Int32Value", "UInt32Value", "BoolValue", "StringValue", "BytesValue":
			return true
		}
	case "go.thethings.network/lorawan-stack/pkg/ttnpb":
		switch t.Name() {
		case "Picture":
			return true
		}
	}
	return false
}

func isSelectableField(name string) bool {
	switch name {
	case "created_at", "updated_at", "ids":
		return false
	}
	return true
}

func isSettableField(name string) bool {
	switch name {
	case "attributes", "contact_info", "password_updated_at", "temporary_password_created_at",
		"temporary_password_expires_at", "antennas", "profile_picture":
		return false
	}
	return true
}

func enumValues(t reflect.Type) map[string]int32 {
	if t.PkgPath() == "go.thethings.network/lorawan-stack/pkg/ttnpb" {
		return proto.EnumValueMap(fmt.Sprintf("ttn.lorawan.v3.%s", t.Name()))
	}
	return nil
}

func addField(fs *pflag.FlagSet, name string, t reflect.Type, maskOnly bool) {
	if maskOnly {
		if t.Kind() == reflect.Struct && !isAtomicType(t, maskOnly) {
			fs.Bool(name, false, fmt.Sprintf("select the %s field and all sub-fields", name))
		} else {
			fs.Bool(name, false, fmt.Sprintf("select the %s field", name))
		}
		return
	}
	if t.Kind() == reflect.Struct && !isAtomicType(t, maskOnly) {
		return
	}
	if t.PkgPath() == "" {
		switch t.Kind() {
		case reflect.Bool:
			fs.Bool(name, false, "")
		case reflect.String:
			fs.String(name, "", "")
		case reflect.Int32:
			fs.Int32(name, 0, "")
		case reflect.Int64:
			fs.Int64(name, 0, "")
		case reflect.Uint32:
			fs.Uint32(name, 0, "")
		case reflect.Uint64:
			fs.Uint64(name, 0, "")
		case reflect.Float32:
			fs.Float32(name, 0, "")
		case reflect.Float64:
			fs.Float64(name, 0, "")
		case reflect.Slice:
			switch t.Elem().Kind() {
			case reflect.Bool:
				fs.BoolSlice(name, nil, "")
			case reflect.String:
				fs.StringSlice(name, nil, "")
			case reflect.Int32:
				if valueMap := enumValues(t); valueMap != nil {
					values := make([]string, 0, len(valueMap))
					for value := range valueMap {
						values = append(values, value)
					}
					fs.StringSlice(name, nil, strings.Join(values, "|"))
				} else {
					fs.IntSlice(name, nil, "")
				}
			case reflect.Int64:
				fs.IntSlice(name, nil, "")
			case reflect.Uint8:
				fs.String(name, "", "(hex)")
			case reflect.Uint32, reflect.Uint64:
				fs.UintSlice(name, nil, "")
			case reflect.Ptr:
				// Not supported
			default:
				fmt.Printf("flags: %s slice not yet supported (%s)\n", t.Elem().Kind(), name)
			}
		case reflect.Map:
			// Not supported
		default:
			fmt.Printf("flags: %s not yet supported (%s)\n", t.Kind(), name)
		}
	} else if t.Kind() == reflect.Int32 && strings.HasSuffix(t.PkgPath(), "ttnpb") {
		if valueMap := enumValues(t); valueMap != nil {
			values := make([]string, 0, len(valueMap))
			for value := range valueMap {
				values = append(values, value)
			}
			fs.String(name, "", strings.Join(values, "|"))
		}
	} else if (t.Kind() == reflect.Slice || t.Kind() == reflect.Array) && t.Elem().Kind() == reflect.Uint8 {
		fs.String(name, "", "(hex)")
	} else {
		switch t.PkgPath() {
		case "time":
			switch t.Name() {
			case "Time":
				fs.String(name, "", "(YYYY-MM-DDTHH:MM:SSZ)")
			case "Duration":
				fs.Duration(name, 0, "(1h2m3s)")
			}
		case "github.com/gogo/protobuf/types":
			switch t.Name() {
			case "DoubleValue":
				fs.Float64(name, 0, "")
			case "FloatValue":
				fs.Float32(name, 0, "")
			case "Int64Value":
				fs.Int64(name, 0, "")
			case "UInt64Value":
				fs.Uint64(name, 0, "")
			case "Int32Value":
				fs.Int32(name, 0, "")
			case "UInt32Value":
				fs.Uint32(name, 0, "")
			case "BoolValue":
				fs.Bool(name, false, "")
			case "StringValue":
				fs.String(name, "", "")
			case "BytesValue":
				fs.String(name, "", "(hex)")
			}
		default:
			fmt.Printf("flags: %s.%s not yet supported (%s)\n", t.PkgPath(), t.Name(), name)
		}
	}
}

func fieldMaskFlags(prefix []string, t reflect.Type, maskOnly bool) *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	props := proto.GetProperties(t)
	for _, prop := range props.Prop {
		if prop.Tag == 0 {
			continue
		}
		field, ok := t.FieldByName(prop.Name)
		if !ok {
			continue
		}
		path := append(prefix, prop.OrigName)
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		name := strings.Join(path, ".")
		if !isSelectableField(name) {
			continue
		}
		if !maskOnly && !isSettableField(name) {
			continue
		}
		addField(flagSet, name, fieldType, maskOnly)
		if fieldType.Kind() == reflect.Struct && !isAtomicType(fieldType, maskOnly) {
			flagSet.AddFlagSet(fieldMaskFlags(path, fieldType, maskOnly))
		}
	}
	return flagSet
}

func trimPrefix(path []string, prefix ...string) []string {
	var nextElement int
	for i, prefix := range prefix {
		if i >= len(path) || path[i] != prefix {
			break
		}
		nextElement = i + 1
	}
	return path[nextElement:]
}

var flagValueError = errors.DefineInvalidArgument("flag_value", "invalid flag value")

func SetFields(dst interface{}, flags *pflag.FlagSet, prefix ...string) error {
	var flagValueErrorAttributes []interface{}
	rv := reflect.Indirect(reflect.ValueOf(dst))
	flags.VisitAll(func(flag *pflag.Flag) {
		if !flag.Changed {
			return
		}
		flagName := toUnderscore.Replace(flag.Name)
		var v interface{}
		switch flag.Value.Type() {
		case "bool":
			v, _ = flags.GetBool(flagName)
		case "string":
			v, _ = flags.GetString(flagName)
		case "int32":
			v, _ = flags.GetInt32(flagName)
		case "int64":
			v, _ = flags.GetInt64(flagName)
		case "uint32":
			v, _ = flags.GetUint32(flagName)
		case "uint64":
			v, _ = flags.GetUint64(flagName)
		case "float32":
			v, _ = flags.GetFloat32(flagName)
		case "float64":
			v, _ = flags.GetFloat64(flagName)
		case "stringSlice":
			v, _ = flags.GetStringSlice(flagName)
		case "intSlice":
			v, _ = flags.GetIntSlice(flagName)
		case "uintSlice":
			v, _ = flags.GetUintSlice(flagName)
		case "duration":
			v, _ = flags.GetDuration(flagName)
		}
		if v == nil {
			flagValueErrorAttributes = append(flagValueErrorAttributes,
				flag.Name, fmt.Errorf("can't set field to %s (%v)", flag.Value, flag.Value.Type()),
			)
		}
		if err := setField(rv, trimPrefix(strings.Split(flagName, "."), prefix...), reflect.ValueOf(v)); err != nil {
			flagValueErrorAttributes = append(flagValueErrorAttributes,
				flag.Name, err.Error(),
			)
		}
	})
	if len(flagValueErrorAttributes) > 0 {
		return flagValueError.WithAttributes(flagValueErrorAttributes...)
	}
	return nil
}

var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

func setField(rv reflect.Value, path []string, v reflect.Value) error {
	rt := rv.Type()
	vt := v.Type()
	props := proto.GetProperties(rt)
	for _, prop := range props.Prop {
		if prop.OrigName == path[0] {
			field := rv.FieldByName(prop.Name)
			if field.Type().Kind() == reflect.Ptr {
				if field.IsNil() {
					field.Set(reflect.New(field.Type().Elem()))
				}
				field = field.Elem()
			}
			ft := field.Type()
			if len(path) == 1 {
				switch {
				case ft.AssignableTo(vt):
					field.Set(v)
				case ft.Kind() == reflect.Int32 && vt.Kind() == reflect.String:
					if valueMap := enumValues(ft); valueMap != nil {
						if enumValue, ok := valueMap[v.String()]; ok {
							field.SetInt(int64(enumValue))
							break
						}
					}

					if reflect.PtrTo(ft).Implements(textUnmarshalerType) {
						fv := reflect.New(ft)
						if err := fv.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(v.String())); err == nil {
							field.SetInt(fv.Elem().Int())
							break
						}
					}
					return fmt.Errorf(`invalid value "%s" for %s`, v.String(), ft.Name())
				case ft.PkgPath() == "time":
					switch {
					case ft.Name() == "Time" && vt.Kind() == reflect.String:
						var t time.Time
						var err error
						if v.String() != "" {
							t, err = time.Parse(time.RFC3339Nano, v.String())
							if err != nil {
								return err
							}
						}
						field.Set(reflect.ValueOf(t))
					case ft.Name() == "Duration" && vt.Kind() == reflect.String:
						d, err := time.ParseDuration(v.String())
						if err != nil {
							return err
						}
						field.Set(reflect.ValueOf(d))
					}
				case ft.PkgPath() == "github.com/gogo/protobuf/types":
					switch ft.Name() {
					case "DoubleValue":
						field.Set(reflect.ValueOf(types.DoubleValue{Value: v.Float()}))
					case "FloatValue":
						field.Set(reflect.ValueOf(types.FloatValue{Value: float32(v.Float())}))
					case "Int64Value":
						field.Set(reflect.ValueOf(types.Int64Value{Value: v.Int()}))
					case "UInt64Value":
						field.Set(reflect.ValueOf(types.UInt64Value{Value: v.Uint()}))
					case "Int32Value":
						field.Set(reflect.ValueOf(types.Int32Value{Value: int32(v.Int())}))
					case "UInt32Value":
						field.Set(reflect.ValueOf(types.UInt32Value{Value: uint32(v.Uint())}))
					case "BoolValue":
						field.Set(reflect.ValueOf(types.BoolValue{Value: v.Bool()}))
					case "StringValue":
						field.Set(reflect.ValueOf(types.StringValue{Value: v.String()}))
					case "BytesValue":
						s := strings.TrimPrefix(v.String(), "0x")
						buf, err := hex.DecodeString(s)
						if err != nil {
							return err
						}
						field.Set(reflect.ValueOf(types.BytesValue{Value: buf}))
					}
				case ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Uint8 && vt.Kind() == reflect.String:
					s := strings.TrimPrefix(v.String(), "0x")
					buf, err := hex.DecodeString(s)
					if err != nil {
						return err
					}
					field.Set(reflect.ValueOf(buf))
				case ft.Kind() == reflect.Array && ft.Elem().Kind() == reflect.Uint8 && vt.Kind() == reflect.String:
					s := strings.TrimPrefix(v.String(), "0x")
					buf, err := hex.DecodeString(s)
					if err != nil {
						return err
					}
					if len(buf) > 0 {
						if len(buf) != ft.Len() {
							return fmt.Errorf(`bytes of "%s" do not fit in [%d]byte`, v.String(), ft.Len())
						}
						for i := 0; i < ft.Len(); i++ {
							field.Index(i).Set(reflect.ValueOf(buf[i]))
						}
					} else {
						field.Set(reflect.Zero(ft))
					}
				case ft.Kind() == reflect.Slice && vt.Kind() == reflect.Slice:
					if vt.Elem().ConvertibleTo(ft.Elem()) {
						slice := reflect.MakeSlice(ft, v.Len(), v.Len())
						for i := 0; i < v.Len(); i++ {
							slice.Index(i).Set(v.Index(i).Convert(ft.Elem()))
						}
						field.Set(slice)
					} else {
						return fmt.Errorf("%v is not convertible to %v", ft, vt)
					}
				default:
					return fmt.Errorf("%v is not assignable to %v", ft, vt)
				}
				return nil
			}
			return setField(field, path[1:], v)
		}
	}
	return fmt.Errorf("unknown field")
}
