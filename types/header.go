package types

import (
	"bytes"
	"encoding/binary"
)

// See https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1 for more information
type Header struct {
	ID    uint16
	Flags uint16
	// QR      bool
	// Opcode  uint8
	// AA      bool
	// TC      bool
	// RD      bool
	// RA      bool
	// Z       uint8
	// RCode   uint8
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

func NewHeader(id, flags, qdCount, anCount, nsCount, arCount uint16) *Header {
	return &Header{
		ID: id,
		// QR:      false,
		// Opcode:  0,
		// AA:      false,
		// TC:      false,
		// RD:      false,
		// RA:      false,
		// Z:       0,
		// RCode:   0,
		Flags:   flags,
		QDCount: qdCount,
		ANCount: anCount,
		NSCount: nsCount,
		ARCount: arCount,
	}
}

func GenerateFlag(qr, opcode, aa, tc, rd, ra, z, rcode uint16) uint16 {
	return uint16(qr<<15 | opcode<<11 | aa<<10 | tc<<9 | rd<<8 | ra<<7 | z<<4 | rcode)
}

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
