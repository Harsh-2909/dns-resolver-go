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
