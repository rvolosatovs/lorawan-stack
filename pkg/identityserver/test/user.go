// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package test

import (
	"fmt"
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/identityserver/types"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/smartystreets/assertions"
)

func defaultUser(in interface{}) (*ttnpb.User, error) {
	if u, ok := in.(types.User); ok {
		return u.GetUser(), nil
	}

	if u, ok := in.(ttnpb.User); ok {
		return &u, nil
	}

	if ptr, ok := in.(*ttnpb.User); ok {
		return ptr, nil
	}

	return nil, fmt.Errorf("Expected: '%v' to be of type ttnpb.User but it wasn't", in)
}

// ShouldBeUser checks if two users resemble each other.
func ShouldBeUser(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return fmt.Sprintf("Expected: one user to match but got %v", len(expected))
	}

	a, s := defaultUser(actual)
	if s != nil {
		return s.Error()
	}

	b, s := defaultUser(expected[0])
	if s != nil {
		return s.Error()
	}

	return all(
		ShouldBeUserIgnoringAutoFields(a, b),
		assertions.ShouldHappenWithin(a.CreatedAt, time.Millisecond, b.CreatedAt),
	)
}

// ShouldBeUserIgnoringAutoFields checks if two users resemble each other
// without looking at fields that are generated by the database: joined.
func ShouldBeUserIgnoringAutoFields(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return fmt.Sprintf("Expected: one user to match but got %v", len(expected))
	}

	a, s := defaultUser(actual)
	if s != nil {
		return s.Error()
	}

	b, s := defaultUser(expected[0])
	if s != nil {
		return s.Error()
	}

	return all(
		assertions.ShouldEqual(a.UserID, b.UserID),
		assertions.ShouldEqual(a.Email, b.Email),
		assertions.ShouldEqual(a.Name, b.Name),
		assertions.ShouldEqual(a.Password, b.Password),
		assertions.ShouldBeTrue(a.ValidatedAt.Equal(b.ValidatedAt)),
		assertions.ShouldEqual(a.Admin, b.Admin),
		assertions.ShouldBeTrue(a.ArchivedAt.Equal(b.ArchivedAt)),
	)
}
