package network

import (
	"dns-resolver-go/dns"
	"fmt"
	"net"
	"time"
)

// Client represents a UDP client for sending DNS queries.
type Client struct {
	ipAddress string
	port      int
}

// NewClient creates a new Client instance.
func NewClient(addr string, port int) *Client {
	return &Client{
		ipAddress: addr,
		port:      port,
	}
}

// ipType returns the IP type of the client's IP address.
func (c *Client) ipType() (string, error) {
	ip := net.ParseIP(c.ipAddress)
	if ip.To4() != nil {
		return "ipv4", nil
		// addr := fmt.Sprintf("%s:%d", c.ipAddress, c.port)
	} else if ip.To16() != nil {
		return "ipv6", nil
		// addr := fmt.Sprintf("[%s]:%d", c.ipAddress, c.port)
	}
	return "", fmt.Errorf("invalid IP address: %s", c.ipAddress)
}

// Query sends a message to the given ip address and port and returns the response.
func (c *Client) Query(message []byte) ([]byte, error) {
	// Create a UDP connection
	ipType, err := c.ipType()
	var addr string
	if err != nil {
		return nil, fmt.Errorf("failed to get the IP type: %v", err)
	}

	if ipType == "ipv4" {
		addr = fmt.Sprintf("%s:%d", c.ipAddress, c.port)
	} else if ipType == "ipv6" {
		addr = fmt.Sprintf("[%s]:%d", c.ipAddress, c.port)
	}
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the DNS server: %v", err)
	}
	defer conn.Close()

	// Set a timeout for the connection
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// Send a message
	_, err = conn.Write(message)
	if err != nil {
		return nil, fmt.Errorf("failed to send the DNS message: %v", err)
	}

	// Receive the response
	buf := make([]byte, 1024)
	// Read the response
	n, err := conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %v", err)
	}
	response := buf[:n]

	// Check if the response ID matches the request ID
	if !IDMatcher(message[:2], response[:2]) {
		return nil, fmt.Errorf("the response ID does not match the request ID")
	}

	return response, nil
}

// IDMatcher checks if the two given IDs match.
func IDMatcher(m1, m2 []byte) bool {
	m1ID := m1[0:2]
	m2ID := m2[0:2]

	return m1ID[0] == m2ID[0] && m1ID[1] == m2ID[1]
}

// Resolve sends a DNS query to the DNS server and returns the 1st answer of the query in parsed format.
// This function recursively queries the DNS server until it finds the Answer of the given type.
// It also prints all the non-authoritative answers in stdout.
func Resolve(domain string, questionType uint16) string {
	question := dns.NewQuestion(domain, questionType, dns.ClassIN)
	flag := dns.NewHeaderFlag(false, 0, false, false, false, false, 0, 0).GenerateFlag()
	header := dns.NewHeader(22, flag, 1, 0, 0, 0)
	DNSMessage := dns.NewDNSMessage(*header, []dns.Question{*question})
	var response []byte
	var parsedResponse *dns.DNSMessage
	var err error
	dnsServerIP := dns.RootDNS
	dnsServerPort := dns.RootDNSPort

	for {
		fmt.Printf("Querying %s for %s\n", dnsServerIP, domain)
		// fmt.Printf("DNS Message:\n %+v\n\n", DNSMessage)
		client := NewClient(dnsServerIP, dnsServerPort)
		response, err = client.Query(DNSMessage.ToBytes())
		if err != nil {
			fmt.Printf("Failed to query the DNS server: %v\n", err)
			return ""
		}
		parsedResponse = dns.DNSMessageFromBytes(response)

		if parsedResponse.Header.ANCount > 0 {
			fmt.Printf("\nNon-authoritative answer:\n")
			if parsedResponse.Answers[0].Type == dns.TypeCNAME {
				fmt.Printf("%s	canonical name = %s.\n", parsedResponse.Answers[0].Name, parsedResponse.Answers[0].RDataParsed)
				Resolve(parsedResponse.Answers[0].RDataParsed, dns.TypeA)
			} else {
				for _, answer := range parsedResponse.Answers {
					fmt.Printf("Name: %s\n", answer.Name)
					fmt.Printf("Address: %s\n", answer.RDataParsed)
				}
			}
			break
		}

		if parsedResponse.Header.ARCount > 0 {
			if ip := getRecord(parsedResponse.AdditionalRRs); ip != "" {
				dnsServerIP = ip
			}
			continue
		}

		if parsedResponse.Header.NSCount > 0 {
			if nsDomain := getRecord(parsedResponse.AuthorityRRs); nsDomain != "" {
				dnsServerIP = Resolve(nsDomain, dns.TypeA)
			}
		}
	}
	return parsedResponse.Answers[0].RDataParsed
}

// getRecord returns the first record of the given type from the given records.
// It is used to get the parsed address of the whitelisted record type.
//
// It returns an empty string if no record of the given type is found.
func getRecord(records []dns.ResourceRecord) string {
	for _, record := range records {
		switch record.Type {
		case dns.TypeA, dns.TypeNS, dns.TypeCNAME:
			return record.RDataParsed
		}
	}
	return ""
}
