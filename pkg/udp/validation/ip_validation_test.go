// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package validation

import (
	"fmt"
	"net"

	"github.com/TheThingsNetwork/ttn/pkg/types"
	"github.com/TheThingsNetwork/ttn/pkg/udp"
)

func ExampleAlwaysValid() {
	udpPacket := udp.Packet{}

	validator := AlwaysValid()
	fmt.Println(validator.Valid(udpPacket))
	// Output: true
}
func ExampleInMemoryValidator() {
	var eui = new(types.EUI64)

	ip1 := &net.UDPAddr{IP: net.IP("8.8.8.8")}
	udpPacket1 := udp.Packet{
		GatewayAddr: ip1,
		GatewayEUI:  eui,
	}

	ip2 := &net.UDPAddr{IP: net.IP("8.8.4.4")}
	udpPacket2 := udp.Packet{
		GatewayAddr: ip2,
		GatewayEUI:  eui,
	}

	validator := InMemoryValidator(DefaultWaitDuration)
	fmt.Println("Validity of first packet transmitted with `eui`:", validator.Valid(udpPacket1))
	fmt.Println("Validity of second packet transmitted with `eui`:", validator.Valid(udpPacket2))
	// Output: Validity of first packet transmitted with `eui`: true
	// Validity of second packet transmitted with `eui`: false
}

func Example() {
	var eui = new(types.EUI64)

	ip1 := &net.UDPAddr{IP: net.IP("8.8.8.8")}
	udpPacket1 := udp.Packet{
		GatewayAddr: ip1,
		GatewayEUI:  eui,
	}

	// Packet from a different IP, but with the same GatewayEUI
	ip2 := &net.UDPAddr{IP: net.IP("8.8.4.4")}
	udpPacket2 := udp.Packet{
		GatewayAddr: ip2,
		GatewayEUI:  eui,
	}

	// Packet from the same IP as the first one, and with the same GatewayEUI
	ip3 := &net.UDPAddr{IP: net.IP("8.8.8.8")}
	udpPacket3 := udp.Packet{
		GatewayAddr: ip3,
		GatewayEUI:  eui,
	}

	validator := InMemoryValidator(DefaultWaitDuration)
	fmt.Println("Validity of first packet transmitted with `eui`:", validator.Valid(udpPacket1))
	fmt.Println("Validity of second packet transmitted with `eui`:", validator.Valid(udpPacket2))
	fmt.Println("Validity of third packet transmitted with `eui`:", validator.Valid(udpPacket3))
	// Output: Validity of first packet transmitted with `eui`: true
	// Validity of second packet transmitted with `eui`: false
	// Validity of third packet transmitted with `eui`: true
}
