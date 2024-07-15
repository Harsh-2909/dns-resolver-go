package dns

import (
	"bytes"
	"encoding/binary"
)

// ResourceRecord represents a DNS resource record.
//
// See https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.3 for more information
type ResourceRecord struct {
	Name     string // The domain name of the resource record
	Type     uint16 // The type of the resource record
	Class    uint16 // The class of the resource record
	TTL      uint32 // The time to live of the resource record
	RDLength uint16 // The length of the resource data
	RData    []byte // The resource data
}

// NewResourceRecord creates a new ResourceRecord instance.
func NewResourceRecord(name string, rType uint16, class uint16, ttl uint32, rdLength uint16, rData []byte) *ResourceRecord {
	return &ResourceRecord{
		Name:     name,
		Type:     rType,
		Class:    class,
		TTL:      ttl,
		RDLength: rdLength,
		RData:    rData,
	}
}

// TrimResourceRecordBytes appends bytes from the buffer until it completely parses all the bytes of a resource record.
// It is useful to trim the bytes of a resource record from a buffer.
func TrimResourceRecordBytes(buf *bytes.Buffer) []byte {
	rrBytes := appendFromBufferUntilNull(buf)
	rrBytes = append(rrBytes, buf.Next(7)...) // appending until ttl
	rdLength := buf.Next(2)
	rrBytes = append(rrBytes, rdLength...) // appending rdLength
	rdLengthCasted := binary.BigEndian.Uint16(rdLength)
	rrBytes = append(rrBytes, buf.Next(int(rdLengthCasted))...) // appending rdata
	return rrBytes
}

// ToBytes converts the ResourceRecord to a byte slice.
func (rr *ResourceRecord) ToBytes() []byte {
	buf := new(bytes.Buffer)

	buf.Write([]byte(rr.Name))
	binary.Write(buf, binary.BigEndian, rr.Type)
	binary.Write(buf, binary.BigEndian, rr.Class)
	binary.Write(buf, binary.BigEndian, rr.TTL)
	binary.Write(buf, binary.BigEndian, rr.RDLength)
	buf.Write(rr.RData)

	return buf.Bytes()
}

// ResourceRecordFromBytes creates a ResourceRecord from a byte slice.
func ResourceRecordFromBytes(data []byte, messageBufs ...*bytes.Buffer) *ResourceRecord {
	buf := bytes.NewBuffer(data)
	var messageBuf *bytes.Buffer
	if messageBufs != nil {
		messageBuf = messageBufs[0]
	}

	name := make([]byte, 0)
	nameLength := 0
	for {
		b, _ := buf.ReadByte()
		// If the byte is 0, then we have reached the end of the name and we can break the loop
		// After name, type is present which is of 2 bytes but 1 of the byte is always 0 so we can break the loop
		if b == 0 {
			break
		}
		name = append(name, b)
		nameLength += 1
	}

	// Check if the name is a pointer. Parse the pointer, get the offset and parse the name from the offset.
	// See https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.4 for more information
	if len(name) == 2 && name[0]>>6 == 0b11 {
		offset := int(name[1])
		if messageBuf != nil {
			messageBytes := messageBuf.Bytes()
			messageBytes = messageBytes[offset:]
			name = appendFromBufferUntilNull(bytes.NewBuffer(messageBytes))
			n, _ := decodeFromQName(string(name))
			name = []byte(n)
		}
	}

	typ := binary.BigEndian.Uint16(data[nameLength : nameLength+2])
	class := binary.BigEndian.Uint16(data[nameLength+2 : nameLength+4])
	ttl := binary.BigEndian.Uint32(data[nameLength+4 : nameLength+8])
	rdLength := binary.BigEndian.Uint16(data[nameLength+8 : nameLength+10])

	return &ResourceRecord{
		Name:     string(name),
		Type:     typ,
		Class:    class,
		TTL:      ttl,
		RDLength: rdLength,
		RData:    data[nameLength+10 : nameLength+10+int(rdLength)], // 10 is the length of the fields before RData
	}
}
