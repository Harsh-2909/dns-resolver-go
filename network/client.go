package network

import (
	"dns-resolver-go/dns"
	"fmt"
	"net"
	"time"
)

type Client struct {
	ipAddress string
	port      int
}

func NewClient(addr string, port int) *Client {
	return &Client{
		ipAddress: addr,
		port:      port,
	}
}

func (c *Client) Query(message []byte) ([]byte, error) {
	// Create a UDP connection
	addr := fmt.Sprintf("%s:%d", c.ipAddress, c.port)
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

func IDMatcher(m1, m2 []byte) bool {
	m1ID := m1[0:2]
	m2ID := m2[0:2]

	return m1ID[0] == m2ID[0] && m1ID[1] == m2ID[1]
}

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
		fmt.Printf("Querying %s for %s\n\n", dnsServerIP, domain)
		// fmt.Printf("DNS Message:\n %+v\n\n", DNSMessage)
		client := NewClient(dnsServerIP, dnsServerPort)
		response, err = client.Query(DNSMessage.ToBytes())
		if err != nil {
			fmt.Printf("Failed to query the DNS server: %v\n", err)
			return ""
		}
		parsedResponse = dns.DNSMessageFromBytes(response)

		fmt.Printf("Parsed Response:\n %+v\n\n", *parsedResponse)
		if parsedResponse.Header.ANCount > 0 {
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
	fmt.Printf("\n\nFinal Answer:\n %+v\n\n", parsedResponse.Answers)
	return getRecord(parsedResponse.Answers)
}

func getRecord(records []dns.ResourceRecord) string {
	for _, record := range records {
		if record.Type == dns.TypeA {
			ipData := record.RData
			return fmt.Sprintf("%d.%d.%d.%d", ipData[0], ipData[1], ipData[2], ipData[3])
		} else if record.Type == dns.TypeNS {
			domainData := record.RData
			decodedName, _ := dns.DecodeName(string(domainData))
			return decodedName
		}
	}
	return ""
}
