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
