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
	"fmt"
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
	"go.thethings.network/lorawan-stack/pkg/redis"
	"go.thethings.network/lorawan-stack/pkg/rpcmetadata"
	"go.thethings.network/lorawan-stack/pkg/util/test"
)

func TestPagination(t *testing.T) {
	a := assertions.New(t)

	for _, tc := range []struct {
		md     rpcmetadata.MD
		limit  int64
		offset int64
	}{
		{
			md:     rpcmetadata.MD{Limit: 10, Page: 1},
			limit:  10,
			offset: 0,
		},
		{
			md:     rpcmetadata.MD{Limit: 10, Page: 2},
			limit:  10,
			offset: 10,
		},
		{
			md:     rpcmetadata.MD{Limit: 10, Page: 3},
			limit:  10,
			offset: 20,
		},
		{
			md:     rpcmetadata.MD{Limit: 0, Page: 1},
			limit:  0,
			offset: 0,
		},
		{
			md:     rpcmetadata.MD{Limit: 0, Page: 2},
			limit:  0,
			offset: 0,
		},
	} {
		t.Run(fmt.Sprintf("limitAndOffsetFromContext, limit:%v, offset:%v", tc.md.Limit, tc.md.Page),
			func(t *testing.T) {
				ctx := tc.md.ToIncomingContext(test.Context())

				ctx = redis.NewContextWithPagination(ctx, 0, 0, nil)

				limit, offset := redis.PaginationLimitAndOffsetFromContext(ctx)

				a.So(limit, should.Equal, tc.limit)
				a.So(offset, should.Equal, tc.offset)

				ctx = redis.NewContextWithPagination(test.Context(), int64(tc.md.Limit), int64(tc.md.Page), nil)

				limit, offset = redis.PaginationLimitAndOffsetFromContext(ctx)

				a.So(limit, should.Equal, tc.limit)
				a.So(offset, should.Equal, tc.offset)
			},
		)
	}

	t.Run("SetTotalCount", func(t *testing.T) {
		var totalCount int64
		ctx := test.Context()
		total := int64(10)

		redis.SetPaginationTotal(ctx, total)
		a.So(totalCount, should.BeZeroValue)

		ctx = redis.NewContextWithPagination(ctx, 5, 1, &totalCount)
		a.So(totalCount, should.BeZeroValue)

		redis.SetPaginationTotal(ctx, total)
		a.So(totalCount, should.Equal, total)
	})
}
