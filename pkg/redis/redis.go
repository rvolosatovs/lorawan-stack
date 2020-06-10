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

// Package redis provides a general Redis client and utilities.
package redis

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/gogo/protobuf/proto"
)

const (
	// separator is character used to separate the keys.
	separator = ':'
)

var (
	encoding = base64.RawStdEncoding
)

// WatchCmdable is transactional redis.Cmdable.
type WatchCmdable interface {
	redis.Cmdable
	Watch(fn func(*redis.Tx) error, keys ...string) error
}

// MarshalProto marshals pb into printable string.
func MarshalProto(pb proto.Message) (string, error) {
	b, err := proto.Marshal(pb)
	if err != nil {
		return "", errEncode.WithCause(err)
	}
	return encoding.EncodeToString(b), nil
}

// UnmarshalProto unmarshals string returned from MarshalProto into pb.
func UnmarshalProto(s string, pb proto.Message) error {
	b, err := encoding.DecodeString(s)
	if err != nil {
		return errDecode.WithCause(err)
	}
	if err = proto.Unmarshal(b, pb); err != nil {
		return errDecode.WithCause(err)
	}
	return nil
}

// Key constructs the full key for entity identified by ks by joining ks using the default separator.
func Key(ks ...string) string {
	return strings.Join(ks, string(separator))
}

// Client represents a Redis store client.
type Client struct {
	*redis.Client
	namespace string
}

// Config represents Redis configuration.
type Config struct {
	Address       string         `name:"address" description:"Address of the Redis server"`
	Password      string         `name:"password" description:"Password of the Redis server"`
	Database      int            `name:"database" description:"Redis database to use"`
	RootNamespace []string       `name:"namespace" description:"Namespace for Redis keys"`
	PoolSize      int            `name:"pool-size" description:"The maximum number of database connections"`
	Failover      FailoverConfig `name:"failover" description:"Redis failover configuration"`
	namespace     []string
}

func (c Config) WithNamespace(namespace ...string) *Config {
	deriv := c
	deriv.namespace = namespace
	return &deriv
}

// IsZero returns whether the Redis configuration is empty.
func (c Config) IsZero() bool {
	if c.Failover.Enable {
		return c.Failover.MasterName == "" && len(c.Failover.Addresses) == 0
	}
	return c.Address == ""
}

// FailoverConfig represents Redis failover configuration.
type FailoverConfig struct {
	Enable     bool     `name:"enable" description:"Enable failover using Redis Sentinel"`
	Addresses  []string `name:"addresses" description:"Redis Sentinel server addresses"`
	MasterName string   `name:"master-name" description:"Redis Sentinel master name"`
}

// newRedisClient returns a Redis client, which connects using correct client type.
func newRedisClient(conf *Config) *redis.Client {
	if conf.Failover.Enable {
		redis.SetLogger(log.New(ioutil.Discard, "", 0))
		return redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    conf.Failover.MasterName,
			SentinelAddrs: conf.Failover.Addresses,
			Password:      conf.Password,
			DB:            conf.Database,
			PoolSize:      conf.PoolSize,
		})
	}
	return redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Password: conf.Password,
		DB:       conf.Database,
		PoolSize: conf.PoolSize,
	})
}

// New returns a new initialized Redis store.
func New(conf *Config) *Client {
	return &Client{
		namespace: Key(append(conf.RootNamespace, conf.namespace...)...),
		Client:    newRedisClient(conf),
	}
}

// Key constructs the full key for entity identified by ks by prepending the configured namespace and joining ks using the default separator.
func (cl *Client) Key(ks ...string) string {
	return Key(append([]string{cl.namespace}, ks...)...)
}

// ProtoCmd is a command, which can unmarshal its result into a protocol buffer.
type ProtoCmd struct {
	result func() (string, error)
}

// ScanProto scans command result into proto.Message pb.
func (cmd ProtoCmd) ScanProto(pb proto.Message) error {
	s, err := cmd.result()
	if err != nil {
		return ConvertError(err)
	}
	return UnmarshalProto(s, pb)
}

// GetProto unmarshals protocol buffer message stored under key k in r into pb.
// Note, that GetProto passes k verbatim to the underlying store and hence, k must represent the full key(including namespace etc.).
func GetProto(r redis.Cmdable, k string) *ProtoCmd {
	return &ProtoCmd{r.Get(k).Result}
}

// SetProto marshals protocol buffer message represented by pb and stores it under key k in r.
// Note, that SetProto passes k verbatim to the underlying store and hence, k must represent the full key(including namespace etc.).
func SetProto(r redis.Cmdable, k string, pb proto.Message, expiration time.Duration) (*redis.StatusCmd, error) {
	s, err := MarshalProto(pb)
	if err != nil {
		return nil, err
	}
	return r.Set(k, s, expiration), nil
}

// FindProto finds the protocol buffer stored under the key stored under k.
// The external key is constructed using keyCmd.
func FindProto(r WatchCmdable, k string, keyCmd func(string) (string, error)) *ProtoCmd {
	var result func() (string, error)
	if err := r.Watch(func(tx *redis.Tx) error {
		id, err := tx.Get(k).Result()
		if err != nil {
			return err
		}
		ik, err := keyCmd(id)
		if err != nil {
			return err
		}
		result = tx.Get(ik).Result
		return nil
	}, k); err != nil {
		return &ProtoCmd{result: func() (string, error) { return "", err }}
	}
	return &ProtoCmd{result: result}
}

type stringSliceCmd struct {
	result func() ([]string, error)
}

// ProtosCmd is a command, which can unmarshal its result into multiple protocol buffers.
type ProtosCmd stringSliceCmd

// Range ranges over command result and unmarshals it into a protocol buffer.
// f must return a new empty proto.Message of the type expected to be present in the command.
// The function returned by f will be called after the commands result is unmarshaled into the message returned by f.
// If both the function returned by f and the message are nil, the entry is skipped.
func (cmd ProtosCmd) Range(f func() (proto.Message, func() (bool, error))) error {
	ss, err := cmd.result()
	if err != nil {
		return err
	}
	for _, s := range ss {
		if s == "" {
			continue
		}

		pb, cb := f()
		if pb == nil && cb == nil {
			continue
		}
		if err := UnmarshalProto(s, pb); err != nil {
			return err
		}
		if ok, err := cb(); err != nil {
			return err
		} else if !ok {
			return nil
		}
	}
	return nil
}

// ProtosWithKeysCmd is a command, which can unmarshal its result into multiple protocol buffers given a key.
type ProtosWithKeysCmd stringSliceCmd

// Range ranges over command result and unmarshals it into a protocol buffer.
// f must return a new empty proto.Message of the type expected to be present in the command given the key.
// The function returned by f will be called after the commands result is unmarshaled into the message returned by f.
// If both the function returned by f and the message are nil, the entry is skipped.
func (cmd ProtosWithKeysCmd) Range(f func(string) (proto.Message, func() (bool, error))) error {
	ss, err := cmd.result()
	if err != nil {
		return err
	}
	if len(ss)%2 != 0 {
		panic(fmt.Sprintf("odd slice length: %d", len(ss)))
	}
	for i := 0; i < len(ss); i += 2 {
		if ss[i+1] == "" {
			continue
		}

		pb, cb := f(ss[i])
		if pb == nil && cb == nil {
			continue
		}
		if err := UnmarshalProto(ss[i+1], pb); err != nil {
			return err
		}
		if ok, err := cb(); err != nil {
			return err
		} else if !ok {
			return nil
		}
	}
	return nil
}

type redisSort struct {
	*redis.Sort
}

// FindProtosOption is an option for the FindProtos query.
type FindProtosOption func(redisSort)

// FindProtosSorted ensures that entries are sorted. If alpha is true, lexicographical sorting is used, otherwise - numerical.
func FindProtosSorted(alpha bool) FindProtosOption {
	return func(s redisSort) {
		s.Alpha = alpha
		s.By = ""
	}
}

// FindProtosWithOffsetAndCount changes the offset and the limit of the query.
func FindProtosWithOffsetAndCount(offset, count int64) FindProtosOption {
	return func(s redisSort) {
		s.Offset, s.Count = offset, count
	}
}

func findProtos(r redis.Cmdable, k string, keyCmd func(string) string, opts ...FindProtosOption) stringSliceCmd {
	s := &redis.Sort{
		Get: []string{keyCmd("*")},
		By:  "nosort", // see https://redis.io/commands/sort#skip-sorting-the-elements
	}
	for _, opt := range opts {
		opt(redisSort{s})
	}
	return stringSliceCmd{
		result: r.Sort(k, s).Result,
	}
}

// FindProtos gets protos stored under keys in k.
func FindProtos(r redis.Cmdable, k string, keyCmd func(string) string, opts ...FindProtosOption) ProtosCmd {
	return ProtosCmd(findProtos(r, k, keyCmd, opts...))
}

// FindProtosWithKeys gets protos stored under keys in k including the keys.
func FindProtosWithKeys(r redis.Cmdable, k string, keyCmd func(string) string, opts ...FindProtosOption) ProtosWithKeysCmd {
	return ProtosWithKeysCmd(findProtos(r, k, keyCmd, append([]FindProtosOption{func(s redisSort) { s.Get = append([]string{"#"}, s.Get...) }}, opts...)...))
}

// ListProtos gets list of protos stored under key k.
func ListProtos(ctx context.Context, r redis.Cmdable, k string) ProtosCmd {
	return ProtosCmd{
		result: r.LRange(k, 0, -1).Result,
	}
}

const (
	payloadKey = "payload"
	replaceKey = "replace"
	startAtKey = "start_at"
)

// InputTaskKey returns the subkey of k, where input tasks are stored.
func InputTaskKey(k string) string {
	return Key(k, "input")
}

// ReadyTaskKey returns the subkey of k, where ready tasks are stored.
func ReadyTaskKey(k string) string {
	return Key(k, "ready")
}

// WaitingTaskKey returns the subkey of k, where waiting tasks are stored.
func WaitingTaskKey(k string) string {
	return Key(k, "waiting")
}

// IsConsumerGroupExistsErr returns true if error represents the redis BUSYGROUP error.
func IsConsumerGroupExistsErr(err error) bool {
	return err != nil && err.Error() == "BUSYGROUP Consumer Group name already exists"
}

// InitTaskGroup initializes the task group for streams at InputTaskKey(k) and ReadyTaskKey(k).
// It must be called before all other task-related functions at subkeys of k.
func InitTaskGroup(r redis.Cmdable, group, k string) error {
	_, err := r.Pipelined(func(p redis.Pipeliner) error {
		p.XGroupCreateMkStream(InputTaskKey(k), group, "$")
		p.XGroupCreateMkStream(ReadyTaskKey(k), group, "$")
		return nil
	})
	if IsConsumerGroupExistsErr(err) {
		return nil
	}
	return ConvertError(err)
}

// AddTask adds a task identified by payload with timestamp startAt to the stream at InputTaskKey(k).
// maxLen is the approximate length of the stream, to which it may be trimmed.
func AddTask(r redis.Cmdable, k string, maxLen int64, payload string, startAt time.Time, replace bool) error {
	m := make(map[string]interface{}, 2)
	m[payloadKey] = payload
	if replace {
		m[replaceKey] = replace
	}
	if !startAt.IsZero() {
		m[startAtKey] = startAt.UnixNano()
	}
	return ConvertError(r.XAdd(&redis.XAddArgs{
		Stream:       InputTaskKey(k),
		MaxLenApprox: maxLen,
		Values:       m,
	}).Err())
}

// DispatchTasks dispatches ready-to-execute tasks from input task streams and waiting task sets to ready task streams.
// It first attempts to read at most maxLen tasks from streams at input task keys corresponding to ks as a consumer id from group group.
// It blocks until deadline, if it is not zero, otherwise it blocks forever.
// It then adds all the tasks read from the stream to the sorted set
// at corresponding waiting task key and acks them.
// Note that task payload is used as the key in the sorted set.
// It then proceeds to add all the tasks from the sorted set, for which execution time is at or before time.Now() to corresponding ready task stream.
func DispatchTasks(r WatchCmdable, group, id string, maxLen int64, deadline time.Time, ks ...string) (time.Time, error) {
	readStreams := make([]string, 0, len(ks))
	for _, k := range ks {
		readStreams = append(readStreams, InputTaskKey(k))
	}
	for range readStreams {
		readStreams = append(readStreams, ">")
	}

	var block time.Duration
	if !deadline.IsZero() {
		block = time.Until(deadline)
		if block <= 0 {
			block = time.Duration(-1)
		}
	}

	rets, err := r.XReadGroup(&redis.XReadGroupArgs{
		Group:    group,
		Consumer: id,
		Streams:  readStreams,
		Count:    maxLen,
		Block:    block,
	}).Result()
	if err != nil && err != redis.Nil {
		return time.Time{}, ConvertError(err)
	}

	if err != redis.Nil {
		_, err := r.Pipelined(func(p redis.Pipeliner) error {
			for i, ret := range rets {
				toAdd := make([]*redis.Z, 0, len(ret.Messages))
				toAddNX := make([]*redis.Z, 0, len(ret.Messages))
				toAck := make([]string, 0, len(ret.Messages))
				for _, msg := range ret.Messages {
					var score float64
					if v, ok := msg.Values[startAtKey]; ok {
						s, ok := v.(string)
						if !ok {
							return errInvalidKeyValueType.WithAttributes("key", startAtKey)
						}

						p, err := strconv.ParseInt(s, 10, 64)
						if err != nil {
							return errInvalidKeyValueType.WithAttributes("key", startAtKey).WithCause(err)
						}
						score = float64(p)
					}

					var member interface{}
					if v, ok := msg.Values[payloadKey]; ok {
						s, ok := v.(string)
						if !ok {
							return errInvalidKeyValueType.WithAttributes("key", payloadKey)
						}
						member = s
					}

					toAck = append(toAck, msg.ID)

					var replace bool
					if v, ok := msg.Values[replaceKey]; ok {
						s, ok := v.(string)
						if !ok {
							return errInvalidKeyValueType.WithAttributes("key", replaceKey)
						}

						p, err := strconv.ParseBool(s)
						if err != nil {
							return errInvalidKeyValueType.WithAttributes("key", replaceKey).WithCause(err)
						}
						replace = p
					}

					if replace {
						toAdd = append(toAdd, &redis.Z{
							Member: member,
							Score:  score,
						})
					} else {
						toAddNX = append(toAddNX, &redis.Z{
							Member: member,
							Score:  score,
						})
					}
				}
				if len(toAdd) > 0 {
					p.ZAdd(WaitingTaskKey(ks[i]), toAdd...)
				}
				if len(toAddNX) > 0 {
					p.ZAddNX(WaitingTaskKey(ks[i]), toAddNX...)
				}
				p.XAck(ret.Stream, group, toAck...)
			}
			return nil
		})
		if err != nil {
			return time.Time{}, ConvertError(err)
		}
	}

	var min time.Time
	for _, k := range ks {
		if err := r.Watch(func(tx *redis.Tx) error {
			zs, err := tx.ZRangeByScoreWithScores(WaitingTaskKey(k), &redis.ZRangeBy{
				Min: "-inf",
				Max: fmt.Sprintf("%d", time.Now().UnixNano()),
			}).Result()
			if err != nil {
				return err
			}

			var minCmd *redis.ZSliceCmd
			_, err = tx.TxPipelined(func(p redis.Pipeliner) error {
				toDel := make([]interface{}, 0, len(zs))
				for _, z := range zs {
					toDel = append(toDel, z.Member)
					p.XAdd(&redis.XAddArgs{
						Stream:       ReadyTaskKey(k),
						MaxLenApprox: maxLen,
						Values: map[string]interface{}{
							payloadKey: z.Member,
							startAtKey: z.Score,
						},
					})
				}
				if len(toDel) > 0 {
					p.ZRem(WaitingTaskKey(k), toDel...)
				}
				minCmd = p.ZRangeWithScores(WaitingTaskKey(k), 0, 0)
				return nil
			})
			if err != nil {
				return err
			}
			if v := minCmd.Val(); len(v) == 1 {
				t := time.Unix(0, int64(v[0].Score))
				if min.IsZero() || t.Before(min) {
					min = t
				}
			}
			return nil
		}, WaitingTaskKey(k)); err != nil {
			return time.Time{}, ConvertError(err)
		}
	}
	return min, nil
}

// PopTask calls f on the most recent task in the queue, for which timestamp is in range [0, time.Now()]
// If timeout value is 0 - PopTask blocks forever
// If timeout value is negative - PopTask does not block
// If timeout value is positive - PopTask blocks until either a task is popped or timeout has passed.
// group is the consumer group name.
// id is the consumer group ID.
// ks are the keys to pop from.
// Tasks are acked if f returns without error.
func PopTask(r redis.Cmdable, group, id string, timeout time.Duration, f func(k string, payload string, startAt time.Time) error, ks ...string) error {
	if len(ks) == 0 {
		return nil
	}

	readStreams := make([]string, 0, len(ks))
	for _, k := range ks {
		readStreams = append(readStreams, ReadyTaskKey(k))
	}
	for range readStreams {
		readStreams = append(readStreams, ">")
	}

	rets, err := r.XReadGroup(&redis.XReadGroupArgs{
		Group:    group,
		Consumer: id,
		Streams:  readStreams,
		Count:    1,
		Block:    timeout,
	}).Result()
	if err != nil && err != redis.Nil {
		return ConvertError(err)
	}
	for i, ret := range rets {
		for _, msg := range ret.Messages {
			var startAt time.Time
			if v, ok := msg.Values[startAtKey]; ok {
				s, ok := v.(string)
				if !ok {
					return errInvalidKeyValueType.WithAttributes("key", startAtKey)
				}
				i, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return errInvalidKeyValueType.WithAttributes("key", startAtKey).WithCause(err)
				}
				startAt = time.Unix(0, i).UTC()
			}

			var payload string
			if v, ok := msg.Values[payloadKey]; ok {
				payload, ok = v.(string)
				if !ok {
					return errInvalidKeyValueType.WithAttributes("key", payloadKey)
				}
			}
			if err := f(ks[i], payload, startAt); err != nil {
				return err
			}
			_, err = r.XAck(ret.Stream, group, msg.ID).Result()
			return ConvertError(err)
		}
	}
	return nil
}

// TaskQueue is a task queue.
type TaskQueue struct {
	Redis     WatchCmdable
	MaxLen    int64
	Group, ID string
	Key       string
}

// Init initializes the task queue.
// It must be called at least once before using the queue.
func (q *TaskQueue) Init() error {
	return InitTaskGroup(q.Redis, q.Group, q.Key)
}

// Run dispatches tasks until ctx.Deadline() is reached(if present) or read on ctx.Done() succeeds.
func (q *TaskQueue) Run(ctx context.Context) error {
	if err := q.Init(); err != nil {
		return err
	}

	var hasDeadline bool
	dl, ok := ctx.Deadline()
	min := dl
	if !ok {
		min = time.Now()
	} else {
		hasDeadline = !dl.IsZero()
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var err error
		min, err = DispatchTasks(q.Redis, q.Group, q.ID, q.MaxLen, min, q.Key)
		if err != nil {
			return err
		}
		if min.IsZero() || hasDeadline && dl.Before(min) {
			min = dl
		}
	}
}

// Add adds a task s to the queue with a timestamp startAt.
func (q *TaskQueue) Add(s string, startAt time.Time, replace bool) error {
	return AddTask(q.Redis, q.Key, q.MaxLen, s, startAt, replace)
}

// Pop calls f on the most recent task in the queue, for which timestamp is in range [0, time.Now()],
// if such is available, otherwise it blocks until it is.
// If ctx.Deadline() is present, Pop will return at or shortly after it.
func (q *TaskQueue) Pop(ctx context.Context, f func(string, time.Time) error) error {
	var timeout time.Duration
	dl, ok := ctx.Deadline()
	if ok {
		timeout = time.Until(dl)
	}
	return PopTask(q.Redis, q.Group, q.ID, timeout, func(_ string, payload string, startAt time.Time) error {
		return f(payload, startAt)
	}, q.Key)
}

// Scripter is redis.scripter.
type Scripter interface {
	Eval(script string, keys []string, args ...interface{}) *redis.Cmd
	EvalSha(sha1 string, keys []string, args ...interface{}) *redis.Cmd
	ScriptExists(hashes ...string) *redis.BoolSliceCmd
	ScriptLoad(script string) *redis.StringCmd
}

var deduplicationScript = redis.NewScript(`local exp = ARGV[1]
local ok = redis.call('set', KEYS[1], '', 'px', exp, 'nx')
if #ARGV > 1 then
	table.remove(ARGV, 1)
	redis.call('rpush', KEYS[2], unpack(ARGV))
	if ok then
		redis.call('pexpire', KEYS[2], exp)
	end
end
if ok then
	return 1
else
	return 0
end`)

// DeduplicationLockKey returns the key deduplication lock for k is stored under.
func DeduplicationLockKey(k string) string {
	return Key(k, "lock")
}

// DeduplicationLockKey returns the key deduplication list for k is stored under.
func DeduplicationListKey(k string) string {
	return Key(k, "list")
}

// DeduplicateProtos deduplicates protos using key k. It stores a lock at DeduplicationLockKey(k) and the list of collected protos at DeduplicationListKey(k).
func DeduplicateProtos(ctx context.Context, r Scripter, k string, window time.Duration, msgs ...proto.Message) (bool, error) {
	args := make([]interface{}, 0, 1+len(msgs))
	args = append(args, window.Milliseconds())
	for _, msg := range msgs {
		s, err := MarshalProto(msg)
		if err != nil {
			return false, err
		}
		args = append(args, s)
	}
	v, err := deduplicationScript.Run(r, []string{DeduplicationLockKey(k), DeduplicationListKey(k)}, args...).Result()
	if err != nil {
		return false, ConvertError(err)
	}
	res, ok := v.(int64)
	if !ok {
		return false, errStore.New()
	}
	return res == 1, nil
}
