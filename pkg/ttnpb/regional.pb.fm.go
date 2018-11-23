// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import (
	fmt "fmt"
	time "time"
)

var ConcentratorConfigFieldPathsNested = []string{
	"channels",
	"clock_source",
	"fsk_channel",
	"fsk_channel.bandwidth",
	"fsk_channel.bit_rate",
	"fsk_channel.channel",
	"fsk_channel.channel.frequency",
	"fsk_channel.channel.radio",
	"lbt",
	"lbt.rssi_offset",
	"lbt.rssi_target",
	"lbt.scan_time",
	"lora_standard_channel",
	"lora_standard_channel.bandwidth",
	"lora_standard_channel.channel",
	"lora_standard_channel.channel.frequency",
	"lora_standard_channel.channel.radio",
	"lora_standard_channel.spreading_factor",
	"ping_slot",
	"ping_slot.frequency",
	"ping_slot.radio",
	"radios",
}

var ConcentratorConfigFieldPathsTopLevel = []string{
	"channels",
	"clock_source",
	"fsk_channel",
	"lbt",
	"lora_standard_channel",
	"ping_slot",
	"radios",
}

func (dst *ConcentratorConfig) SetFields(src *ConcentratorConfig, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "channels":
			if len(subs) > 0 {
				return fmt.Errorf("'channels' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Channels = src.Channels
			} else {
				dst.Channels = nil
			}
		case "lora_standard_channel":
			if len(subs) > 0 {
				newDst := dst.LoRaStandardChannel
				if newDst == nil {
					newDst = &ConcentratorConfig_LoRaStandardChannel{}
					dst.LoRaStandardChannel = newDst
				}
				var newSrc *ConcentratorConfig_LoRaStandardChannel
				if src != nil {
					newSrc = src.LoRaStandardChannel
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.LoRaStandardChannel = src.LoRaStandardChannel
				} else {
					dst.LoRaStandardChannel = nil
				}
			}
		case "fsk_channel":
			if len(subs) > 0 {
				newDst := dst.FSKChannel
				if newDst == nil {
					newDst = &ConcentratorConfig_FSKChannel{}
					dst.FSKChannel = newDst
				}
				var newSrc *ConcentratorConfig_FSKChannel
				if src != nil {
					newSrc = src.FSKChannel
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.FSKChannel = src.FSKChannel
				} else {
					dst.FSKChannel = nil
				}
			}
		case "lbt":
			if len(subs) > 0 {
				newDst := dst.LBT
				if newDst == nil {
					newDst = &ConcentratorConfig_LBTConfiguration{}
					dst.LBT = newDst
				}
				var newSrc *ConcentratorConfig_LBTConfiguration
				if src != nil {
					newSrc = src.LBT
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.LBT = src.LBT
				} else {
					dst.LBT = nil
				}
			}
		case "ping_slot":
			if len(subs) > 0 {
				newDst := dst.PingSlot
				if newDst == nil {
					newDst = &ConcentratorConfig_Channel{}
					dst.PingSlot = newDst
				}
				var newSrc *ConcentratorConfig_Channel
				if src != nil {
					newSrc = src.PingSlot
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.PingSlot = src.PingSlot
				} else {
					dst.PingSlot = nil
				}
			}
		case "radios":
			if len(subs) > 0 {
				return fmt.Errorf("'radios' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Radios = src.Radios
			} else {
				dst.Radios = nil
			}
		case "clock_source":
			if len(subs) > 0 {
				return fmt.Errorf("'clock_source' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ClockSource = src.ClockSource
			} else {
				var zero uint32
				dst.ClockSource = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var ConcentratorConfig_ChannelFieldPathsNested = []string{
	"frequency",
	"radio",
}

var ConcentratorConfig_ChannelFieldPathsTopLevel = []string{
	"frequency",
	"radio",
}

func (dst *ConcentratorConfig_Channel) SetFields(src *ConcentratorConfig_Channel, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "frequency":
			if len(subs) > 0 {
				return fmt.Errorf("'frequency' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Frequency = src.Frequency
			} else {
				var zero uint64
				dst.Frequency = zero
			}
		case "radio":
			if len(subs) > 0 {
				return fmt.Errorf("'radio' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Radio = src.Radio
			} else {
				var zero uint32
				dst.Radio = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var ConcentratorConfig_LoRaStandardChannelFieldPathsNested = []string{
	"bandwidth",
	"channel",
	"channel.frequency",
	"channel.radio",
	"spreading_factor",
}

var ConcentratorConfig_LoRaStandardChannelFieldPathsTopLevel = []string{
	"bandwidth",
	"channel",
	"spreading_factor",
}

func (dst *ConcentratorConfig_LoRaStandardChannel) SetFields(src *ConcentratorConfig_LoRaStandardChannel, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "channel":
			if len(subs) > 0 {
				newDst := &dst.ConcentratorConfig_Channel
				var newSrc *ConcentratorConfig_Channel
				if src != nil {
					newSrc = &src.ConcentratorConfig_Channel
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ConcentratorConfig_Channel = src.ConcentratorConfig_Channel
				} else {
					var zero ConcentratorConfig_Channel
					dst.ConcentratorConfig_Channel = zero
				}
			}
		case "bandwidth":
			if len(subs) > 0 {
				return fmt.Errorf("'bandwidth' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Bandwidth = src.Bandwidth
			} else {
				var zero uint32
				dst.Bandwidth = zero
			}
		case "spreading_factor":
			if len(subs) > 0 {
				return fmt.Errorf("'spreading_factor' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.SpreadingFactor = src.SpreadingFactor
			} else {
				var zero uint32
				dst.SpreadingFactor = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var ConcentratorConfig_FSKChannelFieldPathsNested = []string{
	"bandwidth",
	"bit_rate",
	"channel",
	"channel.frequency",
	"channel.radio",
}

var ConcentratorConfig_FSKChannelFieldPathsTopLevel = []string{
	"bandwidth",
	"bit_rate",
	"channel",
}

func (dst *ConcentratorConfig_FSKChannel) SetFields(src *ConcentratorConfig_FSKChannel, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "channel":
			if len(subs) > 0 {
				newDst := &dst.ConcentratorConfig_Channel
				var newSrc *ConcentratorConfig_Channel
				if src != nil {
					newSrc = &src.ConcentratorConfig_Channel
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ConcentratorConfig_Channel = src.ConcentratorConfig_Channel
				} else {
					var zero ConcentratorConfig_Channel
					dst.ConcentratorConfig_Channel = zero
				}
			}
		case "bandwidth":
			if len(subs) > 0 {
				return fmt.Errorf("'bandwidth' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Bandwidth = src.Bandwidth
			} else {
				var zero uint32
				dst.Bandwidth = zero
			}
		case "bit_rate":
			if len(subs) > 0 {
				return fmt.Errorf("'bit_rate' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.BitRate = src.BitRate
			} else {
				var zero uint32
				dst.BitRate = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

var ConcentratorConfig_LBTConfigurationFieldPathsNested = []string{
	"rssi_offset",
	"rssi_target",
	"scan_time",
}

var ConcentratorConfig_LBTConfigurationFieldPathsTopLevel = []string{
	"rssi_offset",
	"rssi_target",
	"scan_time",
}

func (dst *ConcentratorConfig_LBTConfiguration) SetFields(src *ConcentratorConfig_LBTConfiguration, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "rssi_target":
			if len(subs) > 0 {
				return fmt.Errorf("'rssi_target' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.RSSITarget = src.RSSITarget
			} else {
				var zero float32
				dst.RSSITarget = zero
			}
		case "rssi_offset":
			if len(subs) > 0 {
				return fmt.Errorf("'rssi_offset' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.RSSIOffset = src.RSSIOffset
			} else {
				var zero float32
				dst.RSSIOffset = zero
			}
		case "scan_time":
			if len(subs) > 0 {
				return fmt.Errorf("'scan_time' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ScanTime = src.ScanTime
			} else {
				var zero time.Duration
				dst.ScanTime = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}
