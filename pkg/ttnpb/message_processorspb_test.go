// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/message_processors.proto

package ttnpb

import testing "testing"
import rand "math/rand"
import time "time"
import proto "github.com/gogo/protobuf/proto"
import jsonpb "github.com/gogo/protobuf/jsonpb"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func TestProcessUplinkMessageRequestProto(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := rand.New(rand.NewSource(seed))
	p := NewPopulatedProcessUplinkMessageRequest(popr, false)
	dAtA, err := proto.Marshal(p)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	msg := &ProcessUplinkMessageRequest{}
	if err := proto.Unmarshal(dAtA, msg); err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	littlefuzz := make([]byte, len(dAtA))
	copy(littlefuzz, dAtA)
	for i := range dAtA {
		dAtA[i] = byte(popr.Intn(256))
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("seed = %d, %#v !VerboseProto %#v, since %v", seed, msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("seed = %d, %#v !Proto %#v", seed, msg, p)
	}
	if len(littlefuzz) > 0 {
		fuzzamount := 100
		for i := 0; i < fuzzamount; i++ {
			littlefuzz[popr.Intn(len(littlefuzz))] = byte(popr.Intn(256))
			littlefuzz = append(littlefuzz, byte(popr.Intn(256)))
		}
		// shouldn't panic
		_ = proto.Unmarshal(littlefuzz, msg)
	}
}

func TestProcessUplinkMessageRequestMarshalTo(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := rand.New(rand.NewSource(seed))
	p := NewPopulatedProcessUplinkMessageRequest(popr, false)
	size := p.Size()
	dAtA := make([]byte, size)
	for i := range dAtA {
		dAtA[i] = byte(popr.Intn(256))
	}
	_, err := p.MarshalTo(dAtA)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	msg := &ProcessUplinkMessageRequest{}
	if err := proto.Unmarshal(dAtA, msg); err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	for i := range dAtA {
		dAtA[i] = byte(popr.Intn(256))
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("seed = %d, %#v !VerboseProto %#v, since %v", seed, msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("seed = %d, %#v !Proto %#v", seed, msg, p)
	}
}

func BenchmarkProcessUplinkMessageRequestProtoMarshal(b *testing.B) {
	popr := rand.New(rand.NewSource(616))
	total := 0
	pops := make([]*ProcessUplinkMessageRequest, 10000)
	for i := 0; i < 10000; i++ {
		pops[i] = NewPopulatedProcessUplinkMessageRequest(popr, false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dAtA, err := proto.Marshal(pops[i%10000])
		if err != nil {
			panic(err)
		}
		total += len(dAtA)
	}
	b.SetBytes(int64(total / b.N))
}

func BenchmarkProcessUplinkMessageRequestProtoUnmarshal(b *testing.B) {
	popr := rand.New(rand.NewSource(616))
	total := 0
	datas := make([][]byte, 10000)
	for i := 0; i < 10000; i++ {
		dAtA, err := proto.Marshal(NewPopulatedProcessUplinkMessageRequest(popr, false))
		if err != nil {
			panic(err)
		}
		datas[i] = dAtA
	}
	msg := &ProcessUplinkMessageRequest{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		total += len(datas[i%10000])
		if err := proto.Unmarshal(datas[i%10000], msg); err != nil {
			panic(err)
		}
	}
	b.SetBytes(int64(total / b.N))
}

func TestProcessDownlinkMessageRequestProto(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := rand.New(rand.NewSource(seed))
	p := NewPopulatedProcessDownlinkMessageRequest(popr, false)
	dAtA, err := proto.Marshal(p)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	msg := &ProcessDownlinkMessageRequest{}
	if err := proto.Unmarshal(dAtA, msg); err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	littlefuzz := make([]byte, len(dAtA))
	copy(littlefuzz, dAtA)
	for i := range dAtA {
		dAtA[i] = byte(popr.Intn(256))
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("seed = %d, %#v !VerboseProto %#v, since %v", seed, msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("seed = %d, %#v !Proto %#v", seed, msg, p)
	}
	if len(littlefuzz) > 0 {
		fuzzamount := 100
		for i := 0; i < fuzzamount; i++ {
			littlefuzz[popr.Intn(len(littlefuzz))] = byte(popr.Intn(256))
			littlefuzz = append(littlefuzz, byte(popr.Intn(256)))
		}
		// shouldn't panic
		_ = proto.Unmarshal(littlefuzz, msg)
	}
}

func TestProcessDownlinkMessageRequestMarshalTo(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := rand.New(rand.NewSource(seed))
	p := NewPopulatedProcessDownlinkMessageRequest(popr, false)
	size := p.Size()
	dAtA := make([]byte, size)
	for i := range dAtA {
		dAtA[i] = byte(popr.Intn(256))
	}
	_, err := p.MarshalTo(dAtA)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	msg := &ProcessDownlinkMessageRequest{}
	if err := proto.Unmarshal(dAtA, msg); err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	for i := range dAtA {
		dAtA[i] = byte(popr.Intn(256))
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("seed = %d, %#v !VerboseProto %#v, since %v", seed, msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("seed = %d, %#v !Proto %#v", seed, msg, p)
	}
}

func BenchmarkProcessDownlinkMessageRequestProtoMarshal(b *testing.B) {
	popr := rand.New(rand.NewSource(616))
	total := 0
	pops := make([]*ProcessDownlinkMessageRequest, 10000)
	for i := 0; i < 10000; i++ {
		pops[i] = NewPopulatedProcessDownlinkMessageRequest(popr, false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dAtA, err := proto.Marshal(pops[i%10000])
		if err != nil {
			panic(err)
		}
		total += len(dAtA)
	}
	b.SetBytes(int64(total / b.N))
}

func BenchmarkProcessDownlinkMessageRequestProtoUnmarshal(b *testing.B) {
	popr := rand.New(rand.NewSource(616))
	total := 0
	datas := make([][]byte, 10000)
	for i := 0; i < 10000; i++ {
		dAtA, err := proto.Marshal(NewPopulatedProcessDownlinkMessageRequest(popr, false))
		if err != nil {
			panic(err)
		}
		datas[i] = dAtA
	}
	msg := &ProcessDownlinkMessageRequest{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		total += len(datas[i%10000])
		if err := proto.Unmarshal(datas[i%10000], msg); err != nil {
			panic(err)
		}
	}
	b.SetBytes(int64(total / b.N))
}

func TestProcessUplinkMessageRequestJSON(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := rand.New(rand.NewSource(seed))
	p := NewPopulatedProcessUplinkMessageRequest(popr, true)
	marshaler := jsonpb.Marshaler{}
	jsondata, err := marshaler.MarshalToString(p)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	msg := &ProcessUplinkMessageRequest{}
	err = jsonpb.UnmarshalString(jsondata, msg)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("seed = %d, %#v !VerboseProto %#v, since %v", seed, msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("seed = %d, %#v !Json Equal %#v", seed, msg, p)
	}
}
func TestProcessDownlinkMessageRequestJSON(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := rand.New(rand.NewSource(seed))
	p := NewPopulatedProcessDownlinkMessageRequest(popr, true)
	marshaler := jsonpb.Marshaler{}
	jsondata, err := marshaler.MarshalToString(p)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	msg := &ProcessDownlinkMessageRequest{}
	err = jsonpb.UnmarshalString(jsondata, msg)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("seed = %d, %#v !VerboseProto %#v, since %v", seed, msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("seed = %d, %#v !Json Equal %#v", seed, msg, p)
	}
}
func TestProcessUplinkMessageRequestProtoText(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := rand.New(rand.NewSource(seed))
	p := NewPopulatedProcessUplinkMessageRequest(popr, true)
	dAtA := proto.MarshalTextString(p)
	msg := &ProcessUplinkMessageRequest{}
	if err := proto.UnmarshalText(dAtA, msg); err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("seed = %d, %#v !VerboseProto %#v, since %v", seed, msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("seed = %d, %#v !Proto %#v", seed, msg, p)
	}
}

func TestProcessUplinkMessageRequestProtoCompactText(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := rand.New(rand.NewSource(seed))
	p := NewPopulatedProcessUplinkMessageRequest(popr, true)
	dAtA := proto.CompactTextString(p)
	msg := &ProcessUplinkMessageRequest{}
	if err := proto.UnmarshalText(dAtA, msg); err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("seed = %d, %#v !VerboseProto %#v, since %v", seed, msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("seed = %d, %#v !Proto %#v", seed, msg, p)
	}
}

func TestProcessDownlinkMessageRequestProtoText(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := rand.New(rand.NewSource(seed))
	p := NewPopulatedProcessDownlinkMessageRequest(popr, true)
	dAtA := proto.MarshalTextString(p)
	msg := &ProcessDownlinkMessageRequest{}
	if err := proto.UnmarshalText(dAtA, msg); err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("seed = %d, %#v !VerboseProto %#v, since %v", seed, msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("seed = %d, %#v !Proto %#v", seed, msg, p)
	}
}

func TestProcessDownlinkMessageRequestProtoCompactText(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := rand.New(rand.NewSource(seed))
	p := NewPopulatedProcessDownlinkMessageRequest(popr, true)
	dAtA := proto.CompactTextString(p)
	msg := &ProcessDownlinkMessageRequest{}
	if err := proto.UnmarshalText(dAtA, msg); err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("seed = %d, %#v !VerboseProto %#v, since %v", seed, msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("seed = %d, %#v !Proto %#v", seed, msg, p)
	}
}

func TestProcessUplinkMessageRequestVerboseEqual(t *testing.T) {
	popr := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := NewPopulatedProcessUplinkMessageRequest(popr, false)
	dAtA, err := proto.Marshal(p)
	if err != nil {
		panic(err)
	}
	msg := &ProcessUplinkMessageRequest{}
	if err := proto.Unmarshal(dAtA, msg); err != nil {
		panic(err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("%#v !VerboseEqual %#v, since %v", msg, p, err)
	}
}
func TestProcessDownlinkMessageRequestVerboseEqual(t *testing.T) {
	popr := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := NewPopulatedProcessDownlinkMessageRequest(popr, false)
	dAtA, err := proto.Marshal(p)
	if err != nil {
		panic(err)
	}
	msg := &ProcessDownlinkMessageRequest{}
	if err := proto.Unmarshal(dAtA, msg); err != nil {
		panic(err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("%#v !VerboseEqual %#v, since %v", msg, p, err)
	}
}
func TestProcessUplinkMessageRequestSize(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := rand.New(rand.NewSource(seed))
	p := NewPopulatedProcessUplinkMessageRequest(popr, true)
	size2 := proto.Size(p)
	dAtA, err := proto.Marshal(p)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	size := p.Size()
	if len(dAtA) != size {
		t.Errorf("seed = %d, size %v != marshalled size %v", seed, size, len(dAtA))
	}
	if size2 != size {
		t.Errorf("seed = %d, size %v != before marshal proto.Size %v", seed, size, size2)
	}
	size3 := proto.Size(p)
	if size3 != size {
		t.Errorf("seed = %d, size %v != after marshal proto.Size %v", seed, size, size3)
	}
}

func BenchmarkProcessUplinkMessageRequestSize(b *testing.B) {
	popr := rand.New(rand.NewSource(616))
	total := 0
	pops := make([]*ProcessUplinkMessageRequest, 1000)
	for i := 0; i < 1000; i++ {
		pops[i] = NewPopulatedProcessUplinkMessageRequest(popr, false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		total += pops[i%1000].Size()
	}
	b.SetBytes(int64(total / b.N))
}

func TestProcessDownlinkMessageRequestSize(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := rand.New(rand.NewSource(seed))
	p := NewPopulatedProcessDownlinkMessageRequest(popr, true)
	size2 := proto.Size(p)
	dAtA, err := proto.Marshal(p)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	size := p.Size()
	if len(dAtA) != size {
		t.Errorf("seed = %d, size %v != marshalled size %v", seed, size, len(dAtA))
	}
	if size2 != size {
		t.Errorf("seed = %d, size %v != before marshal proto.Size %v", seed, size, size2)
	}
	size3 := proto.Size(p)
	if size3 != size {
		t.Errorf("seed = %d, size %v != after marshal proto.Size %v", seed, size, size3)
	}
}

func BenchmarkProcessDownlinkMessageRequestSize(b *testing.B) {
	popr := rand.New(rand.NewSource(616))
	total := 0
	pops := make([]*ProcessDownlinkMessageRequest, 1000)
	for i := 0; i < 1000; i++ {
		pops[i] = NewPopulatedProcessDownlinkMessageRequest(popr, false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		total += pops[i%1000].Size()
	}
	b.SetBytes(int64(total / b.N))
}

//These tests are generated by github.com/gogo/protobuf/plugin/testgen
