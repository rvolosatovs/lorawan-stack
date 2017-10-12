// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package test

import (
	"fmt"
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/identityserver/types"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/smartystreets/assertions"
)

func defaultClient(in interface{}) (*ttnpb.Client, error) {
	if cli, ok := in.(types.Client); ok {
		return cli.GetClient(), nil
	}

	if cli, ok := in.(ttnpb.Client); ok {
		return &cli, nil
	}

	if ptr, ok := in.(*ttnpb.Client); ok {
		return ptr, nil
	}

	return nil, fmt.Errorf("Expected: '%v' to be of type ttnpb.Client but it was not", in)
}

// ShouldBeClient checks if two clients resemble each other.
func ShouldBeClient(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return fmt.Sprintf("Expected: one client to match but got %v", len(expected))
	}

	a, s := defaultClient(actual)
	if s != nil {
		return s.Error()
	}

	b, s := defaultClient(expected[0])
	if s != nil {
		return s.Error()
	}

	return all(
		ShouldBeClientIgnoringAutoFields(a, b),
		assertions.ShouldHappenWithin(a.CreatedAt, time.Millisecond, b.CreatedAt),
	)
}

// ShouldBeClientIgnoringAutoFields checks if two clients resemble each other
// without looking at fields that are generated by the database: created.
func ShouldBeClientIgnoringAutoFields(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return fmt.Sprintf("Expected: one client to match but got %v", len(expected))
	}

	a, s := defaultClient(actual)
	if s != nil {
		return s.Error()
	}

	b, s := defaultClient(expected[0])
	if s != nil {
		return s.Error()
	}

	return all(
		assertions.ShouldEqual(a.ClientID, b.ClientID),
		assertions.ShouldResemble(a.Description, b.Description),
		assertions.ShouldEqual(a.Secret, b.Secret),
		assertions.ShouldEqual(a.RedirectURI, b.RedirectURI),
		assertions.ShouldEqual(a.State, b.State),
		assertions.ShouldEqual(a.OfficialLabeled, b.OfficialLabeled),
		assertions.ShouldResemble(a.Grants, b.Grants),
		assertions.ShouldResemble(a.Rights, b.Rights),
		assertions.ShouldBeTrue(a.ArchivedAt.Equal(b.ArchivedAt)),
	)
}
