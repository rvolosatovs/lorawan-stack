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
	"context"
	"time"

	"go.thethings.network/lorawan-stack/pkg/config"
	"go.thethings.network/lorawan-stack/pkg/email"
	"go.thethings.network/lorawan-stack/pkg/email/sendgrid"
	"go.thethings.network/lorawan-stack/pkg/email/smtp"
	"go.thethings.network/lorawan-stack/pkg/fetch"
	"go.thethings.network/lorawan-stack/pkg/oauth"
)

// Config for the Identity Server
type Config struct {
	DatabaseURI      string `name:"database-uri" description:"Database connection URI"`
	UserRegistration struct {
		Invitation struct {
			Required bool          `name:"required" description:"Require invitations for new users"`
			TokenTTL time.Duration `name:"token-ttl" description:"TTL of user invitation tokens"`
		} `name:"invitation"`
		ContactInfoValidation struct {
			Required bool `name:"required" description:"Require contact info validation for new users"`
		} `name:"contact-info-validation"`
		AdminApproval struct {
			Required bool `name:"required" description:"Require admin approval for new users"`
		} `name:"admin-approval"`
		PasswordRequirements struct {
			MinLength    int `name:"min-length" description:"Minimum password length"`
			MaxLength    int `name:"max-length" description:"Maximum password length"`
			MinUppercase int `name:"min-uppercase" description:"Minimum number of uppercase letters"`
			MinDigits    int `name:"min-digits" description:"Minimum number of digits"`
			MinSpecial   int `name:"min-special" description:"Minimum number of special characters"`
		} `name:"password-requirements"`
	} `name:"user-registration"`
	AuthCache struct {
		MembershipTTL time.Duration `name:"membership-ttl" description:"TTL of membership caches"`
	} `name:"auth-cache"`
	OAuth          oauth.Config `name:"oauth"`
	ProfilePicture struct {
		UseGravatar bool   `name:"use-gravatar" description:"Use Gravatar fallback for users without profile picture"`
		Bucket      string `name:"bucket" description:"Bucket used for storing profile pictures"`
		BucketURL   string `name:"bucket-url" description:"Base URL for public bucket access"`
	} `name:"profile-picture"`
	EndDevicePicture struct {
		Bucket    string `name:"bucket" description:"Bucket used for storing end device pictures"`
		BucketURL string `name:"bucket-url" description:"Base URL for public bucket access"`
	} `name:"end-device-picture"`
	Email struct {
		email.Config `name:",squash"`
		SendGrid     sendgrid.Config      `name:"sendgrid"`
		SMTP         smtp.Config          `name:"smtp"`
		Templates    emailTemplatesConfig `name:"templates"`
	} `name:"email"`
}

type emailTemplatesConfig struct {
	Source    string                `name:"source" description:"Source of the email template files (static, directory, url, blob)"`
	Static    map[string][]byte     `name:"-"`
	Directory string                `name:"directory" description:"Retrieve the email templates from the filesystem"`
	URL       string                `name:"url" description:"Retrieve the email templates from a web server"`
	Blob      config.BlobPathConfig `name:"blob"`

	Includes []string `name:"includes" description:"The email templates that will be preloaded on startup"`
}

// Fetcher returns a fetch.Interface based on the configuration.
// If no configuration source is set, this method returns nil, nil.
func (c emailTemplatesConfig) Fetcher(ctx context.Context, blobConf config.BlobConfig) (fetch.Interface, error) {
	// TODO: Remove detection mechanism (https://github.com/TheThingsNetwork/lorawan-stack/issues/1450)
	if c.Source == "" {
		switch {
		case c.Static != nil:
			c.Source = "static"
		case c.Directory != "":
			c.Source = "directory"
		case c.URL != "":
			c.Source = "url"
		case !c.Blob.IsZero():
			c.Source = "blob"
		}
	}
	switch c.Source {
	case "static":
		return fetch.NewMemFetcher(c.Static), nil
	case "directory":
		return fetch.FromFilesystem(c.Directory), nil
	case "url":
		return fetch.FromHTTP(c.URL, true)
	case "blob":
		b, err := blobConf.Bucket(ctx, c.Blob.Bucket)
		if err != nil {
			return nil, err
		}
		return fetch.FromBucket(ctx, b, c.Blob.Path), nil
	default:
		return nil, nil
	}
}
