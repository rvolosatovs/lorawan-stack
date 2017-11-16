// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package auth

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/tokenkey"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestClaims(t *testing.T) {
	a := assertions.New(t)

	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 1800,
			IssuedAt:  time.Now().Unix() - 1800,
			Subject:   ApplicationSubject("foo"),
			Issuer:    "account.thethingsnetwork.org",
		},
		Creator: "john-doe",
		Rights: []ttnpb.Right{
			ttnpb.RIGHT_APPLICATION_INFO,
			ttnpb.RIGHT_APPLICATION_TRAFFIC_READ,
		},
	}

	a.So(claims.Valid(), should.BeNil)
	a.So(claims.ApplicationID(), should.Equal, "foo")
	a.So(claims.HasRights(ttnpb.RIGHT_APPLICATION_INFO), should.BeTrue)
	a.So(claims.HasRights(ttnpb.RIGHT_APPLICATION_DEVICES_WRITE), should.BeFalse)
	a.So(claims.HasRights(ttnpb.RIGHT_APPLICATION_INFO, ttnpb.RIGHT_APPLICATION_TRAFFIC_READ), should.BeTrue)
	a.So(claims.HasRights(ttnpb.RIGHT_APPLICATION_INFO, ttnpb.RIGHT_APPLICATION_DELETE), should.BeFalse)
}

func TestSign(t *testing.T) {
	a := assertions.New(t)

	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 1800,
			IssuedAt:  time.Now().Unix() - 1800,
			Subject:   ApplicationSubject("foo-app"),
			Issuer:    "account.thethingsnetwork.org",
		},
		Creator: "john-doe",
		Rights: []ttnpb.Right{
			ttnpb.RIGHT_APPLICATION_INFO,
			ttnpb.RIGHT_APPLICATION_TRAFFIC_READ,
		},
	}

	// ECDSA512
	{
		key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		a.So(err, should.BeNil)

		token, err := claims.Sign(key)
		a.So(err, should.BeNil)
		a.So(token, should.NotBeEmpty)

		provider := &tokenkey.ConstProvider{
			Tokens: map[string]map[string]crypto.PublicKey{
				claims.Issuer: {
					"": &key.PublicKey,
				},
			},
		}

		parsed, err := ClaimsFromToken(provider, token)
		a.So(err, should.BeNil)
		a.So(parsed, should.Resemble, claims)

		header, err := JOSEHeader(token)
		a.So(err, should.BeNil)
		a.So(header, should.Resemble, &Header{
			Type:      "JWT",
			Algorithm: "ES512",
		})
	}

	// RSA512
	{
		key, err := rsa.GenerateKey(rand.Reader, 2014)
		a.So(err, should.BeNil)

		token, err := claims.Sign(key)
		a.So(err, should.BeNil)
		a.So(token, should.NotBeEmpty)

		provider := &tokenkey.ConstProvider{
			Tokens: map[string]map[string]crypto.PublicKey{
				claims.Issuer: {
					"": &key.PublicKey,
				},
			},
		}

		parsed, err := ClaimsFromToken(provider, token)
		a.So(err, should.BeNil)
		a.So(parsed, should.Resemble, claims)

		header, err := JOSEHeader(token)
		a.So(err, should.BeNil)
		a.So(header, should.Resemble, &Header{
			Type:      "JWT",
			Algorithm: "RS512",
		})
	}

	// ECDSA512 with kid
	{
		key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		a.So(err, should.BeNil)

		kid := "kid-123"
		token, err := claims.Sign(tokenkey.WithKID(key, kid))
		a.So(err, should.BeNil)
		a.So(token, should.NotBeEmpty)

		provider := &tokenkey.ConstProvider{
			Tokens: map[string]map[string]crypto.PublicKey{
				claims.Issuer: {
					kid: &key.PublicKey,
				},
			},
		}

		parsed, err := ClaimsFromToken(provider, token)
		a.So(err, should.BeNil)
		a.So(parsed, should.Resemble, claims)

		header, err := JOSEHeader(token)
		a.So(err, should.BeNil)
		a.So(header, should.Resemble, &Header{
			Type:      "JWT",
			Algorithm: "ES512",
		})
	}

	// ECDSA512 with wrong kid
	{
		key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		a.So(err, should.BeNil)

		token, err := claims.Sign(tokenkey.WithKID(key, "kid"))
		a.So(err, should.BeNil)
		a.So(token, should.NotBeEmpty)

		provider := &tokenkey.ConstProvider{
			Tokens: map[string]map[string]crypto.PublicKey{
				claims.Issuer: {
					"otherkid": &key.PublicKey,
				},
			},
		}

		_, err = ClaimsFromToken(provider, token)
		a.So(err, should.NotBeNil)

		header, err := JOSEHeader(token)
		a.So(err, should.BeNil)
		a.So(header, should.Resemble, &Header{
			Type:      "JWT",
			Algorithm: "ES512",
		})
	}
}
