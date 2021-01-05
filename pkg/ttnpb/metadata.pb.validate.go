// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gogo/protobuf/types"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = types.DynamicAny{}
)

// define the regex for a UUID once up-front
var _metadata_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// ValidateFields checks the field values on RxMetadata with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *RxMetadata) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = RxMetadataFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "gateway_ids":

			if v, ok := interface{}(&m.GatewayIdentifiers).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return RxMetadataValidationError{
						field:  "gateway_ids",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "packet_broker":

			if v, ok := interface{}(m.GetPacketBroker()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return RxMetadataValidationError{
						field:  "packet_broker",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "antenna_index":
			// no validation rules for AntennaIndex
		case "time":

			if v, ok := interface{}(m.GetTime()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return RxMetadataValidationError{
						field:  "time",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "timestamp":
			// no validation rules for Timestamp
		case "fine_timestamp":
			// no validation rules for FineTimestamp
		case "encrypted_fine_timestamp":
			// no validation rules for EncryptedFineTimestamp
		case "encrypted_fine_timestamp_key_id":
			// no validation rules for EncryptedFineTimestampKeyID
		case "rssi":
			// no validation rules for RSSI
		case "signal_rssi":

			if v, ok := interface{}(m.GetSignalRSSI()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return RxMetadataValidationError{
						field:  "signal_rssi",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "channel_rssi":
			// no validation rules for ChannelRSSI
		case "rssi_standard_deviation":
			// no validation rules for RSSIStandardDeviation
		case "snr":
			// no validation rules for SNR
		case "frequency_offset":
			// no validation rules for FrequencyOffset
		case "location":

			if v, ok := interface{}(m.GetLocation()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return RxMetadataValidationError{
						field:  "location",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "downlink_path_constraint":

			if _, ok := DownlinkPathConstraint_name[int32(m.GetDownlinkPathConstraint())]; !ok {
				return RxMetadataValidationError{
					field:  "downlink_path_constraint",
					reason: "value must be one of the defined enum values",
				}
			}

		case "uplink_token":
			// no validation rules for UplinkToken
		case "channel_index":

			if m.GetChannelIndex() > 255 {
				return RxMetadataValidationError{
					field:  "channel_index",
					reason: "value must be less than or equal to 255",
				}
			}

		case "advanced":

			if v, ok := interface{}(m.GetAdvanced()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return RxMetadataValidationError{
						field:  "advanced",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		default:
			return RxMetadataValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// RxMetadataValidationError is the validation error returned by
// RxMetadata.ValidateFields if the designated constraints aren't met.
type RxMetadataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RxMetadataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RxMetadataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RxMetadataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RxMetadataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RxMetadataValidationError) ErrorName() string { return "RxMetadataValidationError" }

// Error satisfies the builtin error interface
func (e RxMetadataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRxMetadata.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RxMetadataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RxMetadataValidationError{}

// ValidateFields checks the field values on Location with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *Location) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = LocationFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "latitude":

			if val := m.GetLatitude(); val < -90 || val > 90 {
				return LocationValidationError{
					field:  "latitude",
					reason: "value must be inside range [-90, 90]",
				}
			}

		case "longitude":

			if val := m.GetLongitude(); val < -180 || val > 180 {
				return LocationValidationError{
					field:  "longitude",
					reason: "value must be inside range [-180, 180]",
				}
			}

		case "altitude":
			// no validation rules for Altitude
		case "accuracy":
			// no validation rules for Accuracy
		case "source":

			if _, ok := LocationSource_name[int32(m.GetSource())]; !ok {
				return LocationValidationError{
					field:  "source",
					reason: "value must be one of the defined enum values",
				}
			}

		default:
			return LocationValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// LocationValidationError is the validation error returned by
// Location.ValidateFields if the designated constraints aren't met.
type LocationValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LocationValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LocationValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LocationValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LocationValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LocationValidationError) ErrorName() string { return "LocationValidationError" }

// Error satisfies the builtin error interface
func (e LocationValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLocation.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LocationValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LocationValidationError{}

// ValidateFields checks the field values on PacketBrokerMetadata with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *PacketBrokerMetadata) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = PacketBrokerMetadataFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "message_id":
			// no validation rules for MessageID
		case "forwarder_net_id":
			// no validation rules for ForwarderNetID
		case "forwarder_tenant_id":
			// no validation rules for ForwarderTenantID
		case "forwarder_cluster_id":
			// no validation rules for ForwarderClusterID
		case "home_network_net_id":
			// no validation rules for HomeNetworkNetID
		case "home_network_tenant_id":
			// no validation rules for HomeNetworkTenantID
		case "home_network_cluster_id":
			// no validation rules for HomeNetworkClusterID
		case "hops":

			for idx, item := range m.GetHops() {
				_, _ = idx, item

				if v, ok := interface{}(item).(interface{ ValidateFields(...string) error }); ok {
					if err := v.ValidateFields(subs...); err != nil {
						return PacketBrokerMetadataValidationError{
							field:  fmt.Sprintf("hops[%v]", idx),
							reason: "embedded message failed validation",
							cause:  err,
						}
					}
				}

			}

		default:
			return PacketBrokerMetadataValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// PacketBrokerMetadataValidationError is the validation error returned by
// PacketBrokerMetadata.ValidateFields if the designated constraints aren't met.
type PacketBrokerMetadataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PacketBrokerMetadataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PacketBrokerMetadataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PacketBrokerMetadataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PacketBrokerMetadataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PacketBrokerMetadataValidationError) ErrorName() string {
	return "PacketBrokerMetadataValidationError"
}

// Error satisfies the builtin error interface
func (e PacketBrokerMetadataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPacketBrokerMetadata.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PacketBrokerMetadataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PacketBrokerMetadataValidationError{}

// ValidateFields checks the field values on PacketBrokerRouteHop with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *PacketBrokerRouteHop) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = PacketBrokerRouteHopFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "received_at":

			if v, ok := interface{}(&m.ReceivedAt).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return PacketBrokerRouteHopValidationError{
						field:  "received_at",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "sender_name":
			// no validation rules for SenderName
		case "sender_address":
			// no validation rules for SenderAddress
		case "receiver_name":
			// no validation rules for ReceiverName
		case "receiver_agent":
			// no validation rules for ReceiverAgent
		default:
			return PacketBrokerRouteHopValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// PacketBrokerRouteHopValidationError is the validation error returned by
// PacketBrokerRouteHop.ValidateFields if the designated constraints aren't met.
type PacketBrokerRouteHopValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PacketBrokerRouteHopValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PacketBrokerRouteHopValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PacketBrokerRouteHopValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PacketBrokerRouteHopValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PacketBrokerRouteHopValidationError) ErrorName() string {
	return "PacketBrokerRouteHopValidationError"
}

// Error satisfies the builtin error interface
func (e PacketBrokerRouteHopValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPacketBrokerRouteHop.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PacketBrokerRouteHopValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PacketBrokerRouteHopValidationError{}
