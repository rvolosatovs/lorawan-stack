// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package udp

// Data contains a LoRaWAN packet
type Data struct {
	RxPacket    []*RxPacket  `json:"rxpk,omitempty"`
	Stat        *Stat        `json:"stat,omitempty"`
	TxPacket    *TxPacket    `json:"txpk,omitempty"`
	TxPacketAck *TxPacketAck `json:"txpk_ack,omitempty"`
}

// RxPacket contains a RX message
type RxPacket struct {
	Time *CompactTime `json:"time,omitempty"` // UTC time of pkt RX, us precision, ISO 8601 'compact' format
	Tmms *uint64      `json:"tmms,omitempty"` // GPS time of pkt RX, number of milliseconds since 06.Jan.1980
	Tmst uint32       `json:"tmst"`           // Internal timestamp of "RX finished" event (32b unsigned)
	Freq float64      `json:"freq"`           // RX central frequency in MHz (unsigned float, Hz precision)
	Chan uint8        `json:"chan"`           // Concentrator "IF" channel used for RX (unsigned integer)
	RFCh uint8        `json:"rfch"`           // Concentrator "RF chain" used for RX (unsigned integer)
	Stat int8         `json:"stat"`           // CRC status: 1 = OK, -1 = fail, 0 = no CRC
	Modu string       `json:"modu"`           // Modulation identifier "LORA" or "FSK"
	DatR DataRate     `json:"datr"`           // LoRa datarate identifier (eg. SF12BW500) or FSK datarate (unsigned, in bits per second)
	CodR string       `json:"codr"`           // LoRa ECC coding rate identifier
	RSSI int16        `json:"rssi"`           // RSSI in dBm (signed integer, 1 dB precision)
	LSNR float64      `json:"lsnr"`           // Lora SNR ratio in dB (signed float, 0.1 dB precision)
	Size uint16       `json:"size"`           // RF packet payload size in bytes (unsigned integer)
	Data string       `json:"data"`           // Base64 encoded RF packet payload, padded
	RSig []RSig       `json:"rsig"`           // Received signal information, per antenna (Optional)
	Brd  uint8        `json:"brd"`            // Concentrator board used for RX (unsigned integer)
	Aesk uint8        `json:"aesk"`           // AES key index used for encrypting fine timestamps
}

// RSig contains the metadata associated with the received signal
type RSig struct {
	Ant    uint8   `json:"ant"`    // Antenna number on which signal has been received
	Chan   uint8   `json:"chan"`   // Concentrator "IF" channel used for RX (unsigned integer)
	RSSIC  int16   `json:"rssic"`  // RSSI in dBm of the channel (signed integer, 1 dB precision)
	RSSIS  int16   `json:"rssis"`  // RSSI in dBm of the signal (signed integer, 1 DB precision) (Optional)
	RSSISD uint16  `json:"rssisd"` // Standard deviation of RSSI during preamble (unsigned integer) (Optional)
	LSNR   float64 `json:"lsnr"`   // Lora SNR ratio in dB (signed float, 0.1 dB precision)
	ETime  string  `json:"etime"`  // Encrypted timestamp, ns precision [0..999999999] (Optional)
	FOff   int32   `json:"foff"`   // Frequency offset in Hz [-125kHz..+125Khz] (Optional)
}

// TxPacket contains a TX message
type TxPacket struct {
	Imme bool         `json:"imme"`           // Send packet immediately (will ignore tmst & time)
	Tmst uint32       `json:"tmst,omitempty"` // Send packet on a certain timestamp value (will ignore time)
	Tmms *uint64      `json:"tmms,omitempty"` // Send packet at a certain GPS time (GPS synchronization required)
	Time *CompactTime `json:"time,omitempty"` // Send packet at a certain time (GPS synchronization required)
	Freq float64      `json:"freq"`           // TX central frequency in MHz (unsigned float, Hz precision)
	Brd  uint8        `json:"brd"`            // Concentrator board used for TX (unsigned integer)
	Ant  uint8        `json:"ant"`            // Concentrator antenna used for TX (unsigned integer)
	RFCh uint8        `json:"rfch"`           // Concentrator "RF chain" used for TX (unsigned integer)
	Powe uint8        `json:"powe"`           // TX output power in dBm (unsigned integer, dBm precision)
	Modu string       `json:"modu"`           // Modulation identifier "LORA" or "FSK"
	DatR DataRate     `json:"datr"`           // LoRa datarate identifier (eg. SF12BW500) || FSK datarate (unsigned, in bits per second)
	CodR string       `json:"codr,omitempty"` // LoRa ECC coding rate identifier
	FDev uint16       `json:"fdev,omitempty"` // FSK frequency deviation (unsigned integer, in Hz)
	IPol bool         `json:"ipol"`           // Lora modulation polarization inversion
	Prea uint16       `json:"prea,omitempty"` // RF preamble size (unsigned integer)
	Size uint16       `json:"size"`           // RF packet payload size in bytes (unsigned integer)
	NCRC bool         `json:"ncrc,omitempty"` // If true, disable the CRC of the physical layer (optional)
	Data string       `json:"data"`           // Base64 encoded RF packet payload, padding optional
}

// Stat contains a status message
type Stat struct {
	Time ExpandedTime  `json:"time"`           // UTC 'system' time of the gateway, ISO 8601 'expanded' format (e.g 2014-01-12 08:59:28 GMT)
	Boot *ExpandedTime `json:"boot,omitempty"` // UTC 'boot' time of the gateway, ISO 8601 'expanded' format (e.g 2014-01-12 08:59:28 GMT)
	Lati *float64      `json:"lati,omitempty"` // GPS latitude of the gateway in degree (float, N is +)
	Long *float64      `json:"long,omitempty"` // GPS latitude of the gateway in degree (float, E is +)
	Alti *int32        `json:"alti,omitempty"` // GPS altitude of the gateway in meter RX (integer)
	RXNb uint32        `json:"rxnb"`           // Number of radio packets received (unsigned integer)
	RXOK uint32        `json:"rxok"`           // Number of radio packets received with a valid PHY CRC
	RXFW uint32        `json:"rxfw"`           // Number of radio packets forwarded (unsigned integer)
	ACKR float64       `json:"ackr"`           // Percentage of upstream datagrams that were acknowledged
	DWNb uint32        `json:"dwnb"`           // Number of downlink datagrams received (unsigned integer)
	TXNb uint32        `json:"txnb"`           // Number of packets emitted (unsigned integer)
	LMOK *uint32       `json:"lmok,omitempty"` // Number of packets received from link testing mote, with CRC OK (unsigned inteter)
	LMST *uint32       `json:"lmst,omitempty"` // Sequence number of the first packet received from link testing mote (unsigned integer)
	LMNW *uint32       `json:"lmnw,omitempty"` // Sequence number of the last packet received from link testing mote (unsigned integer)
	LPPS *uint32       `json:"lpps,omitempty"` // Number of lost PPS pulses (unsigned integer)
	Temp *int32        `json:"temp,omitempty"` // Temperature of the Gateway (signed integer)
	FPGA *uint32       `json:"fpga,omitempty"` // Version of Gateway FPGA (unsigned integer)
	DSP  *uint32       `json:"dsp,omitempty"`  // Version of Gateway DSP software (unsigned interger)
	HAL  *string       `json:"hal,omitempty"`  // Version of Gateway driver (format X.X.X)
}

// TxError is returned in the TxPacketAck
type TxError string

var (
	// TxErrNone is returned if packet has been programmed for downlink
	TxErrNone TxError = "NONE"
	// TxErrTooLate is returned if packet rejected because it was already too late to program this packet for downlink
	TxErrTooLate TxError = "TOO_LATE"
	// TxErrTooEarly is returned if packet rejected because downlink packet timestamp is too much in advance
	TxErrTooEarly TxError = "TOO_EARLY"
	// TxErrCollisionPacket is returned if packet rejected because there was already a packet programmed in requested timeframe
	TxErrCollisionPacket TxError = "COLLISION_PACKET"
	// TxErrCollisionBeacon is returned if packet rejected because there was already a beacon planned in requested timeframe
	TxErrCollisionBeacon TxError = "COLLISION_BEACON"
	// TxErrTxFreq is returned if packet rejected because requested frequency is not supported by TX RF chain
	TxErrTxFreq TxError = "TX_FREQ"
	// TxErrTxPower is returned if packet rejected because requested power is not supported by gateway
	TxErrTxPower TxError = "TX_POWER"
	// TxErrGPSUnlocked is returned if packet rejected because GPS is unlocked, so GPS timestamp cannot be used
	TxErrGPSUnlocked TxError = "GPS_UNLOCKED"
)

// TxPacketAck contains a TX acknowledgment packet
type TxPacketAck struct {
	Error TxError `json:"error"`
}
