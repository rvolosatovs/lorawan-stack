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

package identityserver

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"math/rand"
	"runtime/trace"
	"strings"
	"time"

	ulid "github.com/oklog/ulid/v2"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/picture"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
	"go.thethings.network/lorawan-stack/v3/pkg/util/randutil"
)

const maxProfilePictureStoredDimensions = 1024
const maxEndDevicePictureStoredDimensions = 1024

var profilePictureDimensions = []int{64, 128, 256, 512}
var endDevicePictureDimensions = []int{64, 128, 256, 512}

var pictureRand = randutil.NewLockedRand(rand.NewSource(time.Now().UnixNano()))

func fillGravatar(ctx context.Context, usr *ttnpb.User) (err error) {
	if usr == nil || usr.ProfilePicture != nil || usr.PrimaryEmailAddress == "" {
		return nil
	}
	hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(usr.PrimaryEmailAddress))))
	usr.ProfilePicture = &ttnpb.Picture{
		Sizes: map[uint32]string{},
	}
	for _, size := range profilePictureDimensions {
		usr.ProfilePicture.Sizes[uint32(size)] = fmt.Sprintf("https://www.gravatar.com/avatar/%x?s=%d&d=404", hash, size)
	}
	return nil
}

var errProfilePictureUploadDisabled = errors.DefineFailedPrecondition(
	"profile_picture_uploads_disabled",
	"profile picture uploads are disabled",
)

func (is *IdentityServer) processUserProfilePicture(ctx context.Context, usr *ttnpb.User) (err error) {
	config := is.configFromContext(ctx)
	if config.ProfilePicture.DisableUpload {
		return errProfilePictureUploadDisabled.New()
	}
	// External pictures, consider only largest.
	if usr.ProfilePicture.Sizes != nil {
		original := usr.ProfilePicture.Sizes[0]
		if original == "" {
			var max uint32
			for size, url := range usr.ProfilePicture.Sizes {
				if size > max {
					max = size
					original = url
				}
			}
		}
		if original != "" {
			usr.ProfilePicture.Sizes = map[uint32]string{0: original}
		} else {
			usr.ProfilePicture.Sizes = nil
		}
	}

	// Embedded (uploaded) picture. Make square.
	if usr.ProfilePicture.Embedded != nil && len(usr.ProfilePicture.Embedded.Data) > 0 {
		region := trace.StartRegion(ctx, "make profile picture square")
		usr.ProfilePicture, err = picture.MakeSquare(bytes.NewBuffer(usr.ProfilePicture.Embedded.Data), maxProfilePictureStoredDimensions)
		region.End()
		if err != nil {
			return err
		}
	}

	// Store picture to bucket.
	bucket, err := is.Component.GetBaseConfig(ctx).Blob.Bucket(ctx, config.ProfilePicture.Bucket)
	if err != nil {
		return err
	}
	id := fmt.Sprintf("%s.%s", unique.ID(ctx, usr.UserIdentifiers), ulid.MustNew(ulid.Now(), pictureRand).String())

	region := trace.StartRegion(ctx, "store profile picture")
	usr.ProfilePicture, err = picture.Store(ctx, bucket, id, usr.ProfilePicture, profilePictureDimensions...)
	region.End()
	if err != nil {
		return err
	}

	return
}

var errEndDevicePictureUploadDisabled = errors.DefineFailedPrecondition(
	"end_device_picture_uploads_disabled",
	"end device picture uploads are disabled",
)

func (is *IdentityServer) processEndDevicePicture(ctx context.Context, dev *ttnpb.EndDevice) (err error) {
	config := is.configFromContext(ctx)
	if config.EndDevicePicture.DisableUpload {
		return errEndDevicePictureUploadDisabled.New()
	}
	// External pictures, consider only largest.
	if dev.Picture.Sizes != nil {
		original := dev.Picture.Sizes[0]
		if original == "" {
			var max uint32
			for size, url := range dev.Picture.Sizes {
				if size > max {
					max = size
					original = url
				}
			}
		}
		if original != "" {
			dev.Picture.Sizes = map[uint32]string{0: original}
		} else {
			dev.Picture.Sizes = nil
		}
	}

	// Embedded (uploaded) picture. Make square.
	if dev.Picture.Embedded != nil && len(dev.Picture.Embedded.Data) > 0 {
		region := trace.StartRegion(ctx, "make end device picture square")
		dev.Picture, err = picture.MakeSquare(bytes.NewBuffer(dev.Picture.Embedded.Data), maxEndDevicePictureStoredDimensions)
		region.End()
		if err != nil {
			return err
		}
	}

	// Store picture to bucket.
	bucket, err := is.Component.GetBaseConfig(ctx).Blob.Bucket(ctx, config.EndDevicePicture.Bucket)
	if err != nil {
		return err
	}
	id := fmt.Sprintf("%s.%s", unique.ID(ctx, dev.EndDeviceIdentifiers), ulid.MustNew(ulid.Now(), pictureRand).String())

	region := trace.StartRegion(ctx, "store end device picture")
	dev.Picture, err = picture.Store(ctx, bucket, id, dev.Picture, endDevicePictureDimensions...)
	region.End()
	if err != nil {
		return err
	}

	return
}
