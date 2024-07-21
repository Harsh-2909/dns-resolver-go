package dns

import (
	"bytes"
	"encoding/binary"
)

// HeaderFlag represents the individual flags in the DNS header.
type HeaderFlag struct {
	QR     bool  // QR indicates whether the message is a query (0) or a response (1).
	Opcode uint8 // Opcode specifies the kind of query in the message.
	AA     bool  // AA indicates whether the responding name server is an authority for the domain name in question section.
	TC     bool  // TC indicates whether the message was truncated.
	RD     bool  // RD indicates whether recursion is desired.
	RA     bool  // RA indicates whether recursion is available in the name server.
	Z      uint8 // Z is reserved for future use.
	RCode  uint8 // RCode specifies the response code.
}

// NewHeaderFlag creates a new HeaderFlag instance with the given values.
func NewHeaderFlag(qr bool, opcode uint8, aa bool, tc bool, rd bool, ra bool, z uint8, rcode uint8) *HeaderFlag {
	return &HeaderFlag{
		QR:     qr,
		Opcode: opcode,
		AA:     aa,
		TC:     tc,
		RD:     rd,
		RA:     ra,
		Z:      z,
		RCode:  rcode,
	}
}

// GenerateFlag generates the 16-bit flag value from the individual flag components.
func (hf *HeaderFlag) GenerateFlag() uint16 {
	qr := uint16(boolToInt(hf.QR))
	opcode := uint16(hf.Opcode)
	aa := uint16(boolToInt(hf.AA))
	tc := uint16(boolToInt(hf.TC))
	rd := uint16(boolToInt(hf.RD))
	ra := uint16(boolToInt(hf.RA))
	z := uint16(hf.Z)
	rcode := uint16(hf.RCode)
	return uint16(qr<<15 | opcode<<11 | aa<<10 | tc<<9 | rd<<8 | ra<<7 | z<<4 | rcode)

}

// HeaderFlagFromUint16 creates a HeaderFlag instance from the 16-bit flag value.
func HeaderFlagFromUint16(flag uint16) *HeaderFlag {
	return &HeaderFlag{
		QR:     flag>>15 == 1,
		Opcode: uint8((flag >> 11) & 0b1111),
		AA:     flag>>10 == 1,
		TC:     flag>>9 == 1,
		RD:     flag>>8 == 1,
		RA:     flag>>7 == 1,
		Z:      uint8((flag >> 4) & 0b111),
		RCode:  uint8(flag & 0b1111),
	}
}

// boolToInt converts a boolean value to an integer (0 or 1).
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// ToBytes converts the HeaderFlag to its byte representation.
func (hf *HeaderFlag) ToBytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, hf.GenerateFlag())

	return buf.Bytes()
}

// HeaderFlagFromBytes creates a HeaderFlag instance from its byte representation.
func HeaderFlagFromBytes(b []byte) *HeaderFlag {
	buf := bytes.NewReader(b)

	var flag uint16
	binary.Read(buf, binary.BigEndian, &flag)

	return HeaderFlagFromUint16(flag)
}

// HasError returns whether the HeaderFlag has an error.
// It checks the value of the RCode field.
func (hf *HeaderFlag) HasError() bool {
	return hf.RCode != RCodeNoError
}

// IsQuery returns whether the HeaderFlag is a query.
// It checks the value of the QR field.
func (hf *HeaderFlag) IsQuery() bool {
	return !hf.QR
}

// IsResponse returns whether the HeaderFlag is a response.
// It checks the value of the QR field.
func (hf *HeaderFlag) IsResponse() bool {
	return hf.QR
}
