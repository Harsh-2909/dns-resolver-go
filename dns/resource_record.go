package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

// ResourceRecord represents a DNS resource record.
//
// See https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.3 for more information
type ResourceRecord struct {
	Name        string // The domain name of the resource record
	Type        uint16 // The type of the resource record
	Class       uint16 // The class of the resource record
	TTL         uint32 // The time to live of the resource record
	RDLength    uint16 // The length of the resource data
	RData       []byte // The resource data
	RDataParsed string // The parsed resource data
}

// NewResourceRecord creates a new ResourceRecord instance.
func NewResourceRecord(name string, rType uint16, class uint16, ttl uint32, rdLength uint16, rData []byte) *ResourceRecord {
	rDataParsed, _ := parseRData(rType, rData)
	return &ResourceRecord{
		Name:        name,
		Type:        rType,
		Class:       class,
		TTL:         ttl,
		RDLength:    rdLength,
		RData:       rData,
		RDataParsed: rDataParsed,
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

		// if b>>6 == 0b11 && messageBuf != nil {
		// 	// Check if the name is a pointer. Parse the pointer, get the offset and parse the name from the offset.
		// 	// See https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.4 for more information
		// 	b, _ = buf.ReadByte()
		// 	offset := int(b & 0b11111111)
		// 	messageBytes := messageBuf.Bytes()
		// 	messageBytes = messageBytes[offset:]
		// 	name = appendFromBufferUntilNull(bytes.NewBuffer(messageBytes))
		// 	n, _ := DecodeName(string(name))
		// 	// name = []byte(n)
		// 	name = append(name, []byte(n)...)
		// 	nameLength += 2
		// 	continue
		// }

		name = append(name, b)
		nameLength += 1
	}

	// Check if the name is a pointer. Parse the pointer, get the offset and parse the name from the offset.
	// See https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.4 for more information
	if name[0]>>6 == 0b11 {
		offset := int(name[1])
		if messageBuf != nil {
			messageBytes := messageBuf.Bytes()
			messageBytes = messageBytes[offset:]
			name = appendFromBufferUntilNull(bytes.NewBuffer(messageBytes))
			n, _ := DecodeName(string(name))
			name = []byte(n)
		}
	}

	typ := binary.BigEndian.Uint16(data[nameLength : nameLength+2])
	class := binary.BigEndian.Uint16(data[nameLength+2 : nameLength+4])
	ttl := binary.BigEndian.Uint32(data[nameLength+4 : nameLength+8])
	rdLength := binary.BigEndian.Uint16(data[nameLength+8 : nameLength+10])
	rData := data[nameLength+10 : nameLength+10+int(rdLength)] // 10 is the length of the fields before RData
	rDataParsed, _ := parseRData(typ, rData, messageBuf)

	return &ResourceRecord{
		Name:        string(name),
		Type:        typ,
		Class:       class,
		TTL:         ttl,
		RDLength:    rdLength,
		RData:       rData,
		RDataParsed: rDataParsed,
	}
}

func parseRData(rType uint16, rData []byte, messageBufs ...*bytes.Buffer) (string, error) {
	switch rType {
	case TypeA:
		return parseA(rData)
	case TypeAAAA:
		return parseAAAA(rData)
	case TypeCNAME:
		return parseCNAME(rData, messageBufs...)
	case TypeMX:
		return parseMX(rData)
	case TypeNS:
		return parseNS(rData)
	case TypePTR:
		return parsePTR(rData)
	case TypeSOA:
		return parseSOA(rData)
	case TypeSRV:
		return parseSRV(rData)
	case TypeTXT:
		return parseTXT(rData)
	default:
		return "", fmt.Errorf("unknown resource record type: %d", rType)
	}
}

// parseA parses the A resource record.
func parseA(rData []byte) (string, error) {
	if len(rData) != 4 {
		return "", fmt.Errorf("invalid A record length: %d", len(rData))
	}

	ip := net.IP(rData)
	return ip.String(), nil
}

// parseAAAA parses the AAAA resource record.
func parseAAAA(rData []byte) (string, error) {
	if len(rData) != 16 {
		return "", fmt.Errorf("invalid AAAA record length: %d", len(rData))
	}

	ip := net.IP(rData)
	return ip.String(), nil
}

// parseCNAME parses the CNAME resource record.
func parseCNAME(rData []byte, messageBufs ...*bytes.Buffer) (string, error) {
	if len(rData) == 0 {
		return "", fmt.Errorf("invalid CNAME record length: %d", len(rData))
	}

	name, err := DecodeName(string(rData), messageBufs...)

	return name, err
}

// parseMX parses the MX resource record.
func parseMX(rData []byte) (string, error) {
	if len(rData) < 2 {
		return "", fmt.Errorf("invalid MX record length: %d", len(rData))
	}

	priority := binary.BigEndian.Uint16(rData[0:2])
	name := string(rData[2:])
	return fmt.Sprintf("%d %s", priority, name), nil
}

// parseNS parses the NS resource record.
func parseNS(rData []byte) (string, error) {
	if len(rData) == 0 {
		return "", fmt.Errorf("invalid NS record length: %d", len(rData))
	}

	name := string(rData)
	return name, nil
}

// parsePTR parses the PTR resource record.
func parsePTR(rData []byte) (string, error) {
	if len(rData) == 0 {
		return "", fmt.Errorf("invalid PTR record length: %d", len(rData))
	}

	name := string(rData)
	return name, nil
}

// parseSOA parses the SOA resource record.
func parseSOA(rData []byte) (string, error) {
	if len(rData) < 20 {
		return "", fmt.Errorf("invalid SOA record length: %d", len(rData))
	}

	mname := string(rData[0:rData[0]])
	rname := string(rData[rData[0]:])
	return fmt.Sprintf("%s %s %d %d %d %d %d", mname, rname, binary.BigEndian.Uint32(rData[rData[0]+1:rData[0]+5]), binary.BigEndian.Uint32(rData[rData[0]+5:rData[0]+9]), binary.BigEndian.Uint32(rData[rData[0]+9:rData[0]+13]), binary.BigEndian.Uint32(rData[rData[0]+13:rData[0]+17]), binary.BigEndian.Uint32(rData[rData[0]+17:rData[0]+21])), nil
}

// parseSRV parses the SRV resource record.
func parseSRV(rData []byte) (string, error) {
	if len(rData) < 6 {
		return "", fmt.Errorf("invalid SRV record length: %d", len(rData))
	}

	priority := binary.BigEndian.Uint16(rData[0:2])
	weight := binary.BigEndian.Uint16(rData[2:4])
	port := binary.BigEndian.Uint16(rData[4:6])
	name := string(rData[6:])
	return fmt.Sprintf("%d %d %d %s", priority, weight, port, name), nil
}

// parseTXT parses the TXT resource record.
func parseTXT(rData []byte) (string, error) {
	if len(rData) == 0 {
		return "", fmt.Errorf("invalid TXT record length: %d", len(rData))
	}

	txt := string(rData)
	return txt, nil
}
