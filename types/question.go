package types

import (
	"bytes"
	"encoding/binary"
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
	q.QName = convertToQName(name)
	return q
}

func (q *Question) SetName(name string) {
	q.Name = name
	q.QName = convertToQName(name)
}

func convertToQName(name string) string {
	domain_parts := strings.Split(name, ".")
	qname := ""
	for _, part := range domain_parts {
		new_domain_part := string(byte(len(part))) + part
		qname += new_domain_part
	}
	return qname + "\x00"
}

func (q *Question) ToBytes() []byte {
	buf := new(bytes.Buffer)

	buf.Write([]byte(q.QName))
	binary.Write(buf, binary.BigEndian, q.QType)
	binary.Write(buf, binary.BigEndian, q.QClass)

	return buf.Bytes()
}
