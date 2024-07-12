package types

import (
	"bytes"
	"encoding/binary"
)

// See https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.3 for more information
type ResourceRecord struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    []byte
}

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

func TrimResourceRecordBytes(buf *bytes.Buffer) []byte {
	rrBytes := appendFromBufferUntilNull(buf)
	rrBytes = append(rrBytes, buf.Next(7)...) // appending until ttl
	rdLength := buf.Next(2)
	rrBytes = append(rrBytes, rdLength...) // appending rdLength
	rdLengthCasted := binary.BigEndian.Uint16(rdLength)
	rrBytes = append(rrBytes, buf.Next(int(rdLengthCasted))...) // appending rdata
	return rrBytes
}

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

func ResourceRecordFromBytes(data []byte) *ResourceRecord {
	buf := bytes.NewBuffer(data)

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
