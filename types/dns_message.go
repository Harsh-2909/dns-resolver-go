package types

import (
	"bytes"
	"encoding/binary"
)

// See https://datatracker.ietf.org/doc/html/rfc1035#section-4.1 for more information
type DNSMessage struct {
	Header        Header
	Questions     []Question
	Answers       []ResourceRecord
	AuthorityRRs  []ResourceRecord
	AdditionalRRs []ResourceRecord
}

func (m *DNSMessage) ToBytes() []byte {
	// Create a buffer to store the bytes
	buf := new(bytes.Buffer)

	// Write the header to the buffer
	binary.Write(buf, binary.BigEndian, m.Header)

	// Write the questions to the buffer
	for _, q := range m.Questions {
		buf.Write(q.ToBytes())
	}

	// Write the answers to the buffer
	for _, a := range m.Answers {
		buf.Write(a.ToBytes())
	}

	// Write the authority RRs to the buffer
	for _, rr := range m.AuthorityRRs {
		buf.Write(rr.ToBytes())
	}

	// Write the additional RRs to the buffer
	for _, rr := range m.AdditionalRRs {
		buf.Write(rr.ToBytes())
	}

	// Return the bytes from the buffer
	return buf.Bytes()
}
