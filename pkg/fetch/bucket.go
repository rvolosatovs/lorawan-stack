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

package fetch

import (
	"context"
	"path"
	"time"

	"gocloud.dev/blob"
	"gocloud.dev/gcerrors"
)

type bucketFetcher struct {
	baseFetcher
	context context.Context
	bucket  *blob.Bucket
	root    string
}

func (f *bucketFetcher) File(pathElements ...string) ([]byte, error) {
	if len(pathElements) == 0 {
		return nil, errFilenameNotSpecified
	}

	start := time.Now()

	p := path.Join(pathElements...)
	rp, err := realPath(f.root, p)
	if err != nil {
		return nil, err
	}
	content, err := f.bucket.ReadAll(f.context, rp)
	if err == nil {
		f.observeLatency(time.Since(start))
		return content, nil
	}

	if gcerrors.Code(err) == gcerrors.NotFound {
		return nil, errFileNotFound.WithAttributes("filename", p)
	}
	return nil, errCouldNotReadFile.WithCause(err).WithAttributes("filename", p)
}

// FromBucket returns an interface that fetches files from the given blob bucket.
func FromBucket(ctx context.Context, bucket *blob.Bucket, root string) Interface {
	root = path.Clean(root)
	return &bucketFetcher{
		baseFetcher: baseFetcher{
			latency: fetchLatency.WithLabelValues("bucket", root),
		},
		context: ctx,
		bucket:  bucket,
		root:    root,
	}
}
