// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/joinserver.proto

package ttnpb

import testing "testing"
import math_rand "math/rand"
import time "time"
import github_com_gogo_protobuf_proto "github.com/gogo/protobuf/proto"
import github_com_gogo_protobuf_jsonpb "github.com/gogo/protobuf/jsonpb"
import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func TestAppSKeyRequestProto(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := math_rand.New(math_rand.NewSource(seed))
	p := NewPopulatedAppSKeyRequest(popr, false)
	dAtA, err := github_com_gogo_protobuf_proto.Marshal(p)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	msg := &AppSKeyRequest{}
	if err := github_com_gogo_protobuf_proto.Unmarshal(dAtA, msg); err != nil {
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
		_ = github_com_gogo_protobuf_proto.Unmarshal(littlefuzz, msg)
	}
}

func TestAppSKeyRequestMarshalTo(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := math_rand.New(math_rand.NewSource(seed))
	p := NewPopulatedAppSKeyRequest(popr, false)
	size := p.Size()
	dAtA := make([]byte, size)
	for i := range dAtA {
		dAtA[i] = byte(popr.Intn(256))
	}
	_, err := p.MarshalTo(dAtA)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	msg := &AppSKeyRequest{}
	if err := github_com_gogo_protobuf_proto.Unmarshal(dAtA, msg); err != nil {
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

func BenchmarkAppSKeyRequestProtoMarshal(b *testing.B) {
	popr := math_rand.New(math_rand.NewSource(616))
	total := 0
	pops := make([]*AppSKeyRequest, 10000)
	for i := 0; i < 10000; i++ {
		pops[i] = NewPopulatedAppSKeyRequest(popr, false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dAtA, err := github_com_gogo_protobuf_proto.Marshal(pops[i%10000])
		if err != nil {
			panic(err)
		}
		total += len(dAtA)
	}
	b.SetBytes(int64(total / b.N))
}

func BenchmarkAppSKeyRequestProtoUnmarshal(b *testing.B) {
	popr := math_rand.New(math_rand.NewSource(616))
	total := 0
	datas := make([][]byte, 10000)
	for i := 0; i < 10000; i++ {
		dAtA, err := github_com_gogo_protobuf_proto.Marshal(NewPopulatedAppSKeyRequest(popr, false))
		if err != nil {
			panic(err)
		}
		datas[i] = dAtA
	}
	msg := &AppSKeyRequest{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		total += len(datas[i%10000])
		if err := github_com_gogo_protobuf_proto.Unmarshal(datas[i%10000], msg); err != nil {
			panic(err)
		}
	}
	b.SetBytes(int64(total / b.N))
}

func TestAppSKeyResponseProto(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := math_rand.New(math_rand.NewSource(seed))
	p := NewPopulatedAppSKeyResponse(popr, false)
	dAtA, err := github_com_gogo_protobuf_proto.Marshal(p)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	msg := &AppSKeyResponse{}
	if err := github_com_gogo_protobuf_proto.Unmarshal(dAtA, msg); err != nil {
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
		_ = github_com_gogo_protobuf_proto.Unmarshal(littlefuzz, msg)
	}
}

func TestAppSKeyResponseMarshalTo(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := math_rand.New(math_rand.NewSource(seed))
	p := NewPopulatedAppSKeyResponse(popr, false)
	size := p.Size()
	dAtA := make([]byte, size)
	for i := range dAtA {
		dAtA[i] = byte(popr.Intn(256))
	}
	_, err := p.MarshalTo(dAtA)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	msg := &AppSKeyResponse{}
	if err := github_com_gogo_protobuf_proto.Unmarshal(dAtA, msg); err != nil {
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

func BenchmarkAppSKeyResponseProtoMarshal(b *testing.B) {
	popr := math_rand.New(math_rand.NewSource(616))
	total := 0
	pops := make([]*AppSKeyResponse, 10000)
	for i := 0; i < 10000; i++ {
		pops[i] = NewPopulatedAppSKeyResponse(popr, false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dAtA, err := github_com_gogo_protobuf_proto.Marshal(pops[i%10000])
		if err != nil {
			panic(err)
		}
		total += len(dAtA)
	}
	b.SetBytes(int64(total / b.N))
}

func BenchmarkAppSKeyResponseProtoUnmarshal(b *testing.B) {
	popr := math_rand.New(math_rand.NewSource(616))
	total := 0
	datas := make([][]byte, 10000)
	for i := 0; i < 10000; i++ {
		dAtA, err := github_com_gogo_protobuf_proto.Marshal(NewPopulatedAppSKeyResponse(popr, false))
		if err != nil {
			panic(err)
		}
		datas[i] = dAtA
	}
	msg := &AppSKeyResponse{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		total += len(datas[i%10000])
		if err := github_com_gogo_protobuf_proto.Unmarshal(datas[i%10000], msg); err != nil {
			panic(err)
		}
	}
	b.SetBytes(int64(total / b.N))
}

func TestAppSKeyRequestJSON(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := math_rand.New(math_rand.NewSource(seed))
	p := NewPopulatedAppSKeyRequest(popr, true)
	marshaler := github_com_gogo_protobuf_jsonpb.Marshaler{}
	jsondata, err := marshaler.MarshalToString(p)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	msg := &AppSKeyRequest{}
	err = github_com_gogo_protobuf_jsonpb.UnmarshalString(jsondata, msg)
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
func TestAppSKeyResponseJSON(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := math_rand.New(math_rand.NewSource(seed))
	p := NewPopulatedAppSKeyResponse(popr, true)
	marshaler := github_com_gogo_protobuf_jsonpb.Marshaler{}
	jsondata, err := marshaler.MarshalToString(p)
	if err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	msg := &AppSKeyResponse{}
	err = github_com_gogo_protobuf_jsonpb.UnmarshalString(jsondata, msg)
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
func TestAppSKeyRequestProtoText(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := math_rand.New(math_rand.NewSource(seed))
	p := NewPopulatedAppSKeyRequest(popr, true)
	dAtA := github_com_gogo_protobuf_proto.MarshalTextString(p)
	msg := &AppSKeyRequest{}
	if err := github_com_gogo_protobuf_proto.UnmarshalText(dAtA, msg); err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("seed = %d, %#v !VerboseProto %#v, since %v", seed, msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("seed = %d, %#v !Proto %#v", seed, msg, p)
	}
}

func TestAppSKeyRequestProtoCompactText(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := math_rand.New(math_rand.NewSource(seed))
	p := NewPopulatedAppSKeyRequest(popr, true)
	dAtA := github_com_gogo_protobuf_proto.CompactTextString(p)
	msg := &AppSKeyRequest{}
	if err := github_com_gogo_protobuf_proto.UnmarshalText(dAtA, msg); err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("seed = %d, %#v !VerboseProto %#v, since %v", seed, msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("seed = %d, %#v !Proto %#v", seed, msg, p)
	}
}

func TestAppSKeyResponseProtoText(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := math_rand.New(math_rand.NewSource(seed))
	p := NewPopulatedAppSKeyResponse(popr, true)
	dAtA := github_com_gogo_protobuf_proto.MarshalTextString(p)
	msg := &AppSKeyResponse{}
	if err := github_com_gogo_protobuf_proto.UnmarshalText(dAtA, msg); err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("seed = %d, %#v !VerboseProto %#v, since %v", seed, msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("seed = %d, %#v !Proto %#v", seed, msg, p)
	}
}

func TestAppSKeyResponseProtoCompactText(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := math_rand.New(math_rand.NewSource(seed))
	p := NewPopulatedAppSKeyResponse(popr, true)
	dAtA := github_com_gogo_protobuf_proto.CompactTextString(p)
	msg := &AppSKeyResponse{}
	if err := github_com_gogo_protobuf_proto.UnmarshalText(dAtA, msg); err != nil {
		t.Fatalf("seed = %d, err = %v", seed, err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("seed = %d, %#v !VerboseProto %#v, since %v", seed, msg, p, err)
	}
	if !p.Equal(msg) {
		t.Fatalf("seed = %d, %#v !Proto %#v", seed, msg, p)
	}
}

func TestAppSKeyRequestVerboseEqual(t *testing.T) {
	popr := math_rand.New(math_rand.NewSource(time.Now().UnixNano()))
	p := NewPopulatedAppSKeyRequest(popr, false)
	dAtA, err := github_com_gogo_protobuf_proto.Marshal(p)
	if err != nil {
		panic(err)
	}
	msg := &AppSKeyRequest{}
	if err := github_com_gogo_protobuf_proto.Unmarshal(dAtA, msg); err != nil {
		panic(err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("%#v !VerboseEqual %#v, since %v", msg, p, err)
	}
}
func TestAppSKeyResponseVerboseEqual(t *testing.T) {
	popr := math_rand.New(math_rand.NewSource(time.Now().UnixNano()))
	p := NewPopulatedAppSKeyResponse(popr, false)
	dAtA, err := github_com_gogo_protobuf_proto.Marshal(p)
	if err != nil {
		panic(err)
	}
	msg := &AppSKeyResponse{}
	if err := github_com_gogo_protobuf_proto.Unmarshal(dAtA, msg); err != nil {
		panic(err)
	}
	if err := p.VerboseEqual(msg); err != nil {
		t.Fatalf("%#v !VerboseEqual %#v, since %v", msg, p, err)
	}
}
func TestAppSKeyRequestSize(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := math_rand.New(math_rand.NewSource(seed))
	p := NewPopulatedAppSKeyRequest(popr, true)
	size2 := github_com_gogo_protobuf_proto.Size(p)
	dAtA, err := github_com_gogo_protobuf_proto.Marshal(p)
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
	size3 := github_com_gogo_protobuf_proto.Size(p)
	if size3 != size {
		t.Errorf("seed = %d, size %v != after marshal proto.Size %v", seed, size, size3)
	}
}

func BenchmarkAppSKeyRequestSize(b *testing.B) {
	popr := math_rand.New(math_rand.NewSource(616))
	total := 0
	pops := make([]*AppSKeyRequest, 1000)
	for i := 0; i < 1000; i++ {
		pops[i] = NewPopulatedAppSKeyRequest(popr, false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		total += pops[i%1000].Size()
	}
	b.SetBytes(int64(total / b.N))
}

func TestAppSKeyResponseSize(t *testing.T) {
	seed := time.Now().UnixNano()
	popr := math_rand.New(math_rand.NewSource(seed))
	p := NewPopulatedAppSKeyResponse(popr, true)
	size2 := github_com_gogo_protobuf_proto.Size(p)
	dAtA, err := github_com_gogo_protobuf_proto.Marshal(p)
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
	size3 := github_com_gogo_protobuf_proto.Size(p)
	if size3 != size {
		t.Errorf("seed = %d, size %v != after marshal proto.Size %v", seed, size, size3)
	}
}

func BenchmarkAppSKeyResponseSize(b *testing.B) {
	popr := math_rand.New(math_rand.NewSource(616))
	total := 0
	pops := make([]*AppSKeyResponse, 1000)
	for i := 0; i < 1000; i++ {
		pops[i] = NewPopulatedAppSKeyResponse(popr, false)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		total += pops[i%1000].Size()
	}
	b.SetBytes(int64(total / b.N))
}

//These tests are generated by github.com/gogo/protobuf/plugin/testgen
