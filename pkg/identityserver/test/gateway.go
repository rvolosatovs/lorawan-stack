// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package test

import (
	"fmt"
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/identityserver/types"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/smartystreets/assertions"
)

func defaultGateway(in interface{}) (*ttnpb.Gateway, error) {
	if gtw, ok := in.(types.Gateway); ok {
		return gtw.GetGateway(), nil
	}

	if gtw, ok := in.(ttnpb.Gateway); ok {
		return &gtw, nil
	}

	if ptr, ok := in.(*ttnpb.Gateway); ok {
		return ptr, nil
	}

	return nil, fmt.Errorf("Expected: '%v' to be of type ttnpb.Gateway but it was not", in)
}

// ShouldBeGateway checks if two Gateways resemble each other.
func ShouldBeGateway(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return fmt.Sprintf("Expected: one gateway to match but got %v", len(expected))
	}

	a, s := defaultGateway(actual)
	if s != nil {
		return s.Error()
	}

	b, s := defaultGateway(expected[0])
	if s != nil {
		return s.Error()
	}

	return all(
		ShouldBeGatewayIgnoringAutoFields(a, b),
		assertions.ShouldHappenWithin(a.UpdatedAt, time.Millisecond, b.UpdatedAt),
		assertions.ShouldHappenWithin(a.CreatedAt, time.Millisecond, b.CreatedAt),
	)
}

// ShouldBeGatewayIgnoringAutoFields checks if two Gateways resemble each other
// without looking at fields that are generated by the database: created.
func ShouldBeGatewayIgnoringAutoFields(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return fmt.Sprintf("Expected: one gateway to match but got %v", len(expected))
	}

	a, s := defaultGateway(actual)
	if s != nil {
		return s.Error()
	}

	b, s := defaultGateway(expected[0])
	if s != nil {
		return s.Error()
	}

	return all(
		assertions.ShouldEqual(a.GatewayID, b.GatewayID),
		assertions.ShouldEqual(a.Description, b.Description),
		assertions.ShouldEqual(a.FrequencyPlanID, b.FrequencyPlanID),
		assertions.ShouldBeTrue(a.ActivatedAt.Equal(b.ActivatedAt)),
		assertions.ShouldResemble(a.PrivacySettings, b.PrivacySettings),
		assertions.ShouldEqual(a.AutoUpdate, b.AutoUpdate),
		assertions.ShouldResemble(a.Platform, b.Platform),
		assertions.ShouldResemble(a.Antennas, b.Antennas),
		assertions.ShouldResemble(a.Attributes, b.Attributes),
		assertions.ShouldResemble(a.ClusterAddress, b.ClusterAddress),
		assertions.ShouldBeTrue(a.ArchivedAt.Equal(b.ArchivedAt)),
	)
}

func gatewayAntenna(in interface{}) (*ttnpb.GatewayAntenna, error) {
	if antenna, ok := in.(*ttnpb.GatewayAntenna); ok {
		return antenna, nil
	}

	if antenna, ok := in.(ttnpb.GatewayAntenna); ok {
		return &antenna, nil
	}

	return nil, fmt.Errorf("Expected: '%v' to be of type ttnpb.GatewayAntenna but it was not", in)
}

// ShouldBeGatewayAntenna checks if two Gateway Antennas resemble each other.
func ShouldBeGatewayAntenna(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return fmt.Sprintf("Expected: one gateway antenna to match but got %v", len(expected))
	}

	a, s := gatewayAntenna(actual)
	if s != nil {
		return s.Error()
	}

	b, s := gatewayAntenna(expected[0])
	if s != nil {
		return s.Error()
	}

	return all(
		assertions.ShouldEqual(a.Gain, b.Gain),
		assertions.ShouldResemble(a.Location, b.Location),
		assertions.ShouldEqual(a.Type, b.Type),
		assertions.ShouldEqual(a.Model, b.Model),
		assertions.ShouldEqual(a.Placement, b.Placement),
	)
}
