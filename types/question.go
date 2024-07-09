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

func (q *Question) SetName(name string) {
	q.Name = name
	// Update the QName value based on some algorithm
	// Replace this with your own algorithm
	q.QName = convertToQName(name)
}

func convertToQName(name string) string {
	// Implement your algorithm to convert the name to QName
	// Example algorithm: replace '.' with '-'
	// return strings.ReplaceAll(name, ".", "-")
	domain_parts := strings.Split(name, ".")
	qname := ""
	for _, part := range domain_parts {
		new_domain_part := fmt.Sprintf("%d%s", len(part), part)
		qname += new_domain_part
	}
	return qname + "0"
}

func (q *Question) ToBytes() []byte {
	buf := new(bytes.Buffer)

	buf.Write([]byte(q.QName))
	binary.Write(buf, binary.BigEndian, q.QType)
	binary.Write(buf, binary.BigEndian, q.QClass)

	return buf.Bytes()
}
