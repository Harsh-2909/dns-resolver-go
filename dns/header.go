package dns

import (
	"bytes"
	"encoding/binary"
)

// Header represents the DNS header section.
// It contains the ID, flags, and various count fields.
//
// See https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1 for more information
type Header struct {
	ID      uint16 // ID is a 16-bit identifier assigned by the program that generates any kind of query.
	Flags   uint16 // Flags contains various control flags for the DNS message.
	QDCount uint16 // QDCount specifies the number of entries in the question section.
	ANCount uint16 // ANCount specifies the number of resource records in the answer section.
	NSCount uint16 // NSCount specifies the number of name server resource records in the authority section.
	ARCount uint16 // ARCount specifies the number of resource records in the additional records section.
}

// NewHeader creates a new Header instance with the given values.
func NewHeader(id, flags, qdCount, anCount, nsCount, arCount uint16) *Header {
	return &Header{
		ID:      id,
		Flags:   flags,
		QDCount: qdCount,
		ANCount: anCount,
		NSCount: nsCount,
		ARCount: arCount,
	}
}

// ToBytes converts the Header to its byte representation.
func (h *Header) ToBytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, h.ID)
	binary.Write(buf, binary.BigEndian, h.Flags)
	binary.Write(buf, binary.BigEndian, h.QDCount)
	binary.Write(buf, binary.BigEndian, h.ANCount)
	binary.Write(buf, binary.BigEndian, h.NSCount)
	binary.Write(buf, binary.BigEndian, h.ARCount)

	return buf.Bytes()
}

// HeaderFromBytes creates a Header instance from its byte representation.
func HeaderFromBytes(b []byte) *Header {
	buf := bytes.NewReader(b)

	var id, flags, qdCount, anCount, nsCount, arCount uint16
	binary.Read(buf, binary.BigEndian, &id)
	binary.Read(buf, binary.BigEndian, &flags)
	binary.Read(buf, binary.BigEndian, &qdCount)
	binary.Read(buf, binary.BigEndian, &anCount)
	binary.Read(buf, binary.BigEndian, &nsCount)
	binary.Read(buf, binary.BigEndian, &arCount)

	return &Header{
		ID:      id,
		Flags:   flags,
		QDCount: qdCount,
		ANCount: anCount,
		NSCount: nsCount,
		ARCount: arCount,
	}
}

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

// GenerateFlag generates the 16-bit flag value from the individual flag components.
func GenerateFlag(qr, opcode, aa, tc, rd, ra, z, rcode uint16) uint16 {
	return uint16(qr<<15 | opcode<<11 | aa<<10 | tc<<9 | rd<<8 | ra<<7 | z<<4 | rcode)
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

	binary.Write(buf, binary.BigEndian, GenerateFlag(
		uint16(boolToInt(hf.QR)),
		uint16(hf.Opcode),
		uint16(boolToInt(hf.AA)),
		uint16(boolToInt(hf.TC)),
		uint16(boolToInt(hf.RD)),
		uint16(boolToInt(hf.RA)),
		uint16(hf.Z),
		uint16(hf.RCode),
	))

	return buf.Bytes()
}

// HeaderFlagFromBytes creates a HeaderFlag instance from its byte representation.
func HeaderFlagFromBytes(b []byte) *HeaderFlag {
	buf := bytes.NewReader(b)

	var flags uint16
	binary.Read(buf, binary.BigEndian, &flags)

	return &HeaderFlag{
		QR:     flags>>15 == 1,
		Opcode: uint8((flags >> 11) & 0b1111),
		AA:     flags>>10 == 1,
		TC:     flags>>9 == 1,
		RD:     flags>>8 == 1,
		RA:     flags>>7 == 1,
		Z:      uint8((flags >> 4) & 0b111),
		RCode:  uint8(flags & 0b1111),
	}
}
