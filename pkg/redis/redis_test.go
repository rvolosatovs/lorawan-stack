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

package redis_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gogo/protobuf/proto"
	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	. "go.thethings.network/lorawan-stack/v3/pkg/redis"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestInitTaskGroup(t *testing.T) {
	for _, tc := range []struct {
		Name           string
		Populate       func(context.Context, *Client) bool
		Group, Key     string
		ErrorAssertion func(*testing.T, error) bool
	}{
		{
			Name:           "no streams/no groups",
			Populate:       func(ctx context.Context, cl *Client) bool { return true },
			Group:          "testGroup",
			Key:            "testKey",
			ErrorAssertion: func(t *testing.T, err error) bool { return assertions.New(t).So(err, should.BeNil) },
		},
		{
			Name: "streams exist/groups exist",
			Populate: func(ctx context.Context, cl *Client) bool {
				_, a := test.MustNewTFromContext(ctx)
				_, err := cl.XGroupCreateMkStream(ctx, InputTaskKey(cl.Key("testKey")), cl.Key("testGroup"), "0").Result()
				if !a.So(err, should.BeNil) {
					return false
				}
				_, err = cl.XGroupCreateMkStream(ctx, ReadyTaskKey(cl.Key("testKey")), cl.Key("testGroup"), "0").Result()
				return a.So(err, should.BeNil)
			},
			Group:          "testGroup",
			Key:            "testKey",
			ErrorAssertion: func(t *testing.T, err error) bool { return assertions.New(t).So(err, should.BeNil) },
		},
	} {
		test.RunSubtest(t, test.SubtestConfig{
			Name:     tc.Name,
			Parallel: true,
			Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
				cl, flush := test.NewRedis(ctx, "redis_test")
				defer flush()
				defer cl.Close()

				a.So(tc.Populate(ctx, cl), should.BeTrue)

				err := InitTaskGroup(ctx, cl, cl.Key(tc.Group), cl.Key(tc.Key))
				a.So(tc.ErrorAssertion(t, err), should.BeTrue)
			},
		})
	}
}

func TestAddTask(t *testing.T) {
	a, ctx := test.New(t)

	cl, flush := test.NewRedis(ctx, "redis_test")
	defer flush()
	defer cl.Close()

	err := AddTask(ctx, cl, cl.Key("testKey"), 10, "testPayload", time.Unix(0, 42), false)
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}

	rets, err := cl.Client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{InputTaskKey(cl.Key("testKey")), "0"},
		Count:   10,
		Block:   -1,
	}).Result()
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}

	if a.So(rets, should.HaveLength, 1) {
		a.So(rets[0].Stream, should.Equal, InputTaskKey(cl.Key("testKey")))
		if a.So(rets[0].Messages, should.HaveLength, 1) {
			msg := rets[0].Messages[0]
			a.So(msg, should.Resemble, redis.XMessage{
				ID: msg.ID,
				Values: map[string]interface{}{
					"start_at": fmt.Sprintf("%d", time.Unix(0, 42).UnixNano()),
					"payload":  "testPayload",
				},
			})
		}
	}

	err = AddTask(ctx, cl, cl.Key("testKey"), 10, "testPayload", time.Unix(0, 42), true)
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}

	rets, err = cl.Client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{InputTaskKey(cl.Key("testKey")), "0"},
		Count:   10,
		Block:   -1,
	}).Result()
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}

	if a.So(rets, should.HaveLength, 1) {
		a.So(rets[0].Stream, should.Equal, InputTaskKey(cl.Key("testKey")))
		if a.So(rets[0].Messages, should.HaveLength, 2) {
			msg0 := rets[0].Messages[0]
			a.So(msg0, should.Resemble, redis.XMessage{
				ID: msg0.ID,
				Values: map[string]interface{}{
					"start_at": fmt.Sprintf("%d", time.Unix(0, 42).UnixNano()),
					"payload":  "testPayload",
				},
			})
			msg1 := rets[0].Messages[1]
			a.So(msg1, should.Resemble, redis.XMessage{
				ID: msg1.ID,
				Values: map[string]interface{}{
					"start_at": fmt.Sprintf("%d", time.Unix(0, 42).UnixNano()),
					"payload":  "testPayload",
					"replace":  "1",
				},
			})
		}
	}
}

func TestPopTask(t *testing.T) {
	a, ctx := test.New(t)

	cl, flush := test.NewRedis(ctx, "redis_test")
	defer flush()
	defer cl.Close()

	const (
		testGroup = "testGroup"

		testKey1 = "testKey1"
		testKey2 = "testKey2"
	)

	assertPop := func(ctx context.Context, inputKey, expectedPayload string, expectedStartAt time.Time) bool {
		t, a := test.MustNewTFromContext(ctx)
		t.Helper()

		var called bool
		errCh := make(chan error, 1)
		go func() {
			errCh <- PopTask(ctx, cl.Client, testGroup, "testID", 10, func(p redis.Pipeliner, payload string, startAt time.Time) error {
				p.Ping(ctx)
				if !test.AllTrue(
					a.So(called, should.BeFalse),
					a.So(payload, should.Equal, expectedPayload),
					a.So(startAt, should.Resemble, expectedStartAt),
				) {
					t.Errorf("PopTask assertion failed for task with expected payload %s and expected starting time of %s", expectedPayload, expectedStartAt)
				}
				called = true
				return nil
			}, inputKey)
		}()

		select {
		case <-ctx.Done():
			t.Error("Timed out while waiting for Pop callback to be called")
			return false

		case err := <-errCh:
			if !a.So(err, should.BeNil) {
				t.Errorf("PopTask failed with: %s", test.FormatError(err))
			}
			return a.So(called, should.BeTrue) && !a.Failed()
		}
	}

	testKeys := [...]string{
		cl.Key("testKey1"),
		cl.Key("testKey2"),
	}
	for _, k := range testKeys {
		_, err := cl.XGroupCreateMkStream(ctx, InputTaskKey(k), testGroup, "0").Result()
		if !a.So(err, should.BeNil) {
			t.FailNow()
		}

		_, err = cl.XGroupCreateMkStream(ctx, ReadyTaskKey(k), testGroup, "0").Result()
		if !a.So(err, should.BeNil) {
			t.FailNow()
		}

		timeout := (1 >> 5) * test.Delay

		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		err = PopTask(timeoutCtx, cl.Client, testGroup, "testID", 10, func(redis.Pipeliner, string, time.Time) error {
			panic("must not be called")
		}, k)
		cancel()
		if a.So(err, should.BeError) {
			a.So(errors.IsDeadlineExceeded(err), should.BeTrue)
		}

		cancelCtx, cancel := context.WithCancel(ctx)
		time.AfterFunc(timeout, cancel)
		err = PopTask(cancelCtx, cl.Client, testGroup, "testID", 10, func(redis.Pipeliner, string, time.Time) error {
			panic("must not be called")
		}, k)
		cancel()
		if a.So(err, should.BeError) {
			a.So(errors.IsCanceled(err), should.BeTrue)
		}
	}

	inputKeys := [...]string{
		InputTaskKey(testKeys[0]),
		InputTaskKey(testKeys[1]),
	}

	payloads := [...]string{
		"testPayload",
		"testPayload2",
		"testPayload3",
	}

	now := time.Now()
	nextMin := now.Add(time.Hour)

	for _, x := range []*redis.XAddArgs{
		{
			Stream: inputKeys[0],
			Values: map[string]interface{}{
				"start_at": fmt.Sprintf("%d", time.Unix(0, 42).UnixNano()),
				"payload":  payloads[0],
			},
		},
		{
			Stream: inputKeys[0],
			Values: map[string]interface{}{
				"start_at": fmt.Sprintf("%d", time.Unix(0, 43).UnixNano()),
				"payload":  payloads[0],
			},
		},
		{
			Stream: inputKeys[0],
			Values: map[string]interface{}{
				"start_at": fmt.Sprintf("%d", time.Unix(0, 41).UnixNano()),
				"payload":  payloads[0],
			},
		},
		{
			Stream: inputKeys[0],
			Values: map[string]interface{}{
				"start_at": fmt.Sprintf("%d", time.Unix(0, 41).UnixNano()),
				"payload":  payloads[0],
			},
		},
		{
			Stream: inputKeys[1],
			Values: map[string]interface{}{
				"start_at": fmt.Sprintf("%d", time.Unix(0, 66).UnixNano()),
				"payload":  payloads[0],
			},
		},
		{
			Stream: inputKeys[0],
			Values: map[string]interface{}{
				"start_at": nextMin.UnixNano(),
				"payload":  payloads[1],
			},
		},
		{
			Stream: inputKeys[0],
			Values: map[string]interface{}{
				"start_at": "0",
				"payload":  payloads[2],
			},
		},
		{
			Stream: InputTaskKey(cl.Key("testKeyUnrelated")),
			Values: map[string]interface{}{
				"start_at": "0",
				"payload":  "testPayloadUnrelated",
			},
		},
	} {
		_, err := cl.Client.XAdd(ctx, x).Result()
		if !a.So(err, should.BeNil) {
			t.FailNow()
		}
	}

	a.So(assertPop(ctx, testKeys[0], payloads[2], time.Unix(0, 0).UTC()), should.BeTrue)
	a.So(assertPop(ctx, testKeys[0], payloads[0], time.Unix(0, 42).UTC()), should.BeTrue)
	a.So(assertPop(ctx, testKeys[1], payloads[0], time.Unix(0, 66).UTC()), should.BeTrue)
}

func TestTaskQueue(t *testing.T) {
	a, ctx := test.New(t)

	cl, flush := test.NewRedis(ctx, "redis_test")
	defer flush()
	defer cl.Close()

	q := &TaskQueue{
		Redis:  cl,
		MaxLen: 42,
		Group:  "testGroup",
		ID:     "testID",
		Key:    cl.Key("test"),
	}

	err := q.Init(ctx)
	a.So(err, should.BeNil)
	defer func() {
		err := q.Close(ctx)
		a.So(err, should.BeNil)
	}()

	assertPop := func(ctx context.Context, r redis.Cmdable, expectedPayload string, expectedStartAt time.Time) bool {
		t, a := test.MustNewTFromContext(ctx)
		t.Helper()

		type popFuncReq struct {
			Pipeliner redis.Pipeliner
			Payload   string
			Time      time.Time
			Response  chan<- error
		}

		var called bool
		errCh := make(chan error, 1)
		go func() {
			errCh <- q.Pop(ctx, r, func(p redis.Pipeliner, payload string, startAt time.Time) error {
				p.Ping(ctx)
				a.So(called, should.BeFalse)
				a.So(payload, should.Equal, expectedPayload)
				a.So(startAt, should.Resemble, expectedStartAt)
				called = true
				return nil
			})
		}()

		select {
		case <-ctx.Done():
			t.Error("Timed out while waiting for Pop callback to be called")
			return false

		case err := <-errCh:
			return test.AllTrue(
				a.So(err, should.BeNil),
				a.So(called, should.BeTrue),
			)
		}
	}

	p := cl.Pipeline()
	switch {
	case !a.So(q.Add(ctx, nil, "test", time.Now(), true), should.BeNil),
		!a.So(q.Add(ctx, p, "test", time.Unix(0, 42), true), should.BeNil),
		!a.So(q.Add(ctx, nil, "test", time.Unix(0, 24), false), should.BeNil),
		!a.So(q.Add(ctx, p, "test2", time.Unix(0, 43), false), should.BeNil),
		!a.So(q.Add(ctx, p, "test", time.Unix(0, 420), false), should.BeNil),
		!a.So(func() error {
			_, err := p.Exec(ctx)
			return err
		}(), should.BeNil),
		!a.So(assertPop(ctx, nil, "test", time.Unix(0, 42).UTC()), should.BeTrue),
		!a.So(q.Add(ctx, nil, "test2", time.Unix(0, 41), true), should.BeNil),
		!a.So(assertPop(ctx, nil, "test2", time.Unix(0, 43).UTC()), should.BeTrue),
		!a.So(assertPop(ctx, nil, "test2", time.Unix(0, 41).UTC()), should.BeTrue):
	}
}

func TestProtoDeduplicator(t *testing.T) {
	a, ctx := test.New(t)

	cl, flush := test.NewRedis(ctx, "redis_test")
	defer flush()
	defer cl.Close()

	makeMockProto := func(s string) proto.Message {
		return &test.MockProtoMessageMarshalUnmarshaler{
			MockProtoMarshaler: test.MockProtoMarshaler{
				MarshalFunc: func() ([]byte, error) {
					return []byte(s), nil
				},
			},
		}
	}

	ttl := (1 << 12) * test.Delay
	key1 := cl.Key("test1")
	key2 := cl.Key("test2")

	v, err := DeduplicateProtos(ctx, cl, key1, ttl)
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	a.So(v, should.BeTrue)

	v, err = DeduplicateProtos(ctx, cl, key1, ttl, makeMockProto("proto1"))
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	a.So(v, should.BeFalse)

	v, err = DeduplicateProtos(ctx, cl, key2, ttl, makeMockProto("proto1"))
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	a.So(v, should.BeTrue)

	v, err = DeduplicateProtos(ctx, cl, key1, ttl, makeMockProto("proto1"))
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	a.So(v, should.BeFalse)

	v, err = DeduplicateProtos(ctx, cl, key1, ttl, makeMockProto("proto2"), makeMockProto("proto3"))
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	a.So(v, should.BeFalse)

	v, err = DeduplicateProtos(ctx, cl, key2, ttl, makeMockProto("proto2"))
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	a.So(v, should.BeFalse)

	ss, err := cl.LRange(ctx, ListKey(key1), 0, -1).Result()
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	lockTTL, err := cl.PTTL(ctx, LockKey(key1)).Result()
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	listTTL, err := cl.PTTL(ctx, ListKey(key1)).Result()
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	a.So(ss, should.Resemble, []string{
		Encoding.EncodeToString([]byte("proto1")),
		Encoding.EncodeToString([]byte("proto1")),
		Encoding.EncodeToString([]byte("proto2")),
		Encoding.EncodeToString([]byte("proto3")),
	})
	a.So(lockTTL, should.BeGreaterThan, 0)
	a.So(lockTTL, should.BeLessThanOrEqualTo, ttl)
	a.So(listTTL, should.BeGreaterThan, 0)
	a.So(listTTL, should.BeLessThanOrEqualTo, ttl)

	ss, err = cl.LRange(ctx, ListKey(key2), 0, -1).Result()
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	lockTTL, err = cl.PTTL(ctx, LockKey(key2)).Result()
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	listTTL, err = cl.PTTL(ctx, ListKey(key2)).Result()
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	a.So(ss, should.Resemble, []string{
		Encoding.EncodeToString([]byte("proto1")),
		Encoding.EncodeToString([]byte("proto2")),
	})
	a.So(lockTTL, should.BeGreaterThan, 0)
	a.So(lockTTL, should.BeLessThanOrEqualTo, ttl)
	a.So(listTTL, should.BeGreaterThan, 0)
	a.So(listTTL, should.BeLessThanOrEqualTo, ttl)
}

func TestMutex(t *testing.T) {
	a, ctx := test.New(t)

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second+(1<<8)*test.Delay)
	defer cancel()

	cl, flush := test.NewRedis(ctx, "redis_test")
	defer flush()
	defer cl.Close()

	ttl := (1 << 8) * test.Delay
	key := cl.Key("test1")

	err := LockMutex(ctx, cl, key, "test-id-1", ttl)
	if !a.So(err, should.BeNil) {
		t.Fatalf("Failed to lock mutex: %s", err)
	}

	lockTTL, err := cl.PTTL(ctx, LockKey(key)).Result()
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	a.So(lockTTL, should.BeGreaterThan, 0)
	a.So(lockTTL, should.BeLessThanOrEqualTo, ttl)

	blockErrCh := make(chan error, 1)
	go func() {
		blockErrCh <- LockMutex(ctx, cl, key, "test-id-2", ttl)
	}()

	timeoutErrCh := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(ctx, test.Delay)
		defer cancel()
		timeoutErrCh <- LockMutex(ctx, cl, key, "test-id-3", ttl)
	}()

	select {
	case <-ctx.Done():
		t.Fatalf("Timed out while waiting for LockMutex with a deadline to return")
	case err := <-timeoutErrCh:
		a.So(errors.IsDeadlineExceeded(err), should.BeTrue)
	}
	select {
	case err := <-blockErrCh:
		t.Fatalf("LockMutex returned before previous caller unlocked: %s", err)
	default:
	}

	err = UnlockMutex(ctx, cl, key, "test-id-1", ttl)
	if !a.So(err, should.BeNil) {
		t.Fatalf("Failed to unlock mutex: %s", err)
	}

	select {
	case <-ctx.Done():
		t.Fatalf("Timed out while waiting for blocked LockMutex to return")
	case err := <-blockErrCh:
		a.So(err, should.BeNil)
	}
	lockTTL, err = cl.PTTL(ctx, LockKey(key)).Result()
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	a.So(lockTTL, should.BeGreaterThan, 0)
	a.So(lockTTL, should.BeLessThanOrEqualTo, ttl)

	err = UnlockMutex(ctx, cl, key, "test-id-2", ttl)
	if !a.So(err, should.BeNil) {
		t.Fatalf("Failed to unlock mutex: %s", err)
	}
	lockTTL, err = cl.PTTL(ctx, LockKey(key)).Result()
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	listTTL, err := cl.PTTL(ctx, ListKey(key)).Result()
	if !a.So(err, should.BeNil) {
		t.FailNow()
	}
	a.So(lockTTL, should.BeGreaterThan, 0)
	a.So(lockTTL, should.BeLessThanOrEqualTo, ttl)
	a.So(listTTL, should.BeLessThanOrEqualTo, lockTTL)

	err = UnlockMutex(ctx, cl, key, "non-existent-id", ttl)
	if !a.So(err, should.BeNil) {
		t.Fatalf("Failed to unlock existent mutex with non-existent ID: %s", err)
	}
	err = UnlockMutex(ctx, cl, cl.Key("non-existent-key"), "non-existent-id", ttl)
	if !a.So(err, should.BeNil) {
		t.Fatalf("Failed to unlock non-existent mutex: %s", err)
	}
}
