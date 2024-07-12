package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

// See https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.2 for more information
type Question struct {
	Name   string // This is a domain name
	QName  string // This is the converted domain name based on the RFC 1035 document
	QType  uint16
	QClass uint16
}

func NewQuestion(name string, qType, qClass uint16) *Question {
	q := &Question{
		Name:   name,
		QType:  qType,
		QClass: qClass,
	}
	q.QName = encodeToQName(name)
	return q
}

func (q *Question) SetName(name string) {
	q.Name = name
	q.QName = encodeToQName(name)
}

func encodeToQName(name string) string {
	domainParts := strings.Split(name, ".")
	qname := ""
	for _, part := range domainParts {
		newDomainPart := string(byte(len(part))) + part
		qname += newDomainPart
	}
	return qname + "\x00"
}

func decodeFromQName(qname string) (string, error) {
	encoded := []byte(qname)
	var result bytes.Buffer

	for i := 0; i < len(encoded); {
		length := int(encoded[i])
		if length == 0 {
			break
		}
		i++
		if i+length > len(encoded) {
			return "", fmt.Errorf("invalid encoded domain name")
		}
		if result.Len() > 0 {
			result.WriteByte('.')
		}
		result.Write(encoded[i : i+length])
		i += length
	}

	return result.String(), nil
}

func (q *Question) ToBytes() []byte {
	buf := new(bytes.Buffer)

	buf.Write([]byte(q.QName))
	binary.Write(buf, binary.BigEndian, q.QType)
	binary.Write(buf, binary.BigEndian, q.QClass)

	return buf.Bytes()
}

func QuestionFromBytes(b []byte) *Question {
	length := len(b)
	qname := string(b[:length-4])

	name, _ := decodeFromQName(qname)

	return &Question{
		Name:   name,
		QName:  qname,
		QType:  binary.BigEndian.Uint16(b[length-4 : length-2]),
		QClass: binary.BigEndian.Uint16(b[length-2:]),
	}
}
