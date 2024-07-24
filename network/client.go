package network

import (
	"dns-resolver-go/cache"
	"dns-resolver-go/dns"
	"fmt"
	"net"
	"os"
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
	} else if ip.To16() != nil {
		return "ipv6", nil
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

// Option represents the options for the Resolve function.
type Option struct {
	UseCache bool // Whether to use the cache or not.
}

// DefaultOption returns the default options for the Resolve function.
func DefaultOption() Option {
	return Option{
		UseCache: true,
	}
}

// Resolve sends a DNS query to the DNS server and returns the 1st answer of the query in parsed format.
// This function recursively queries the DNS server until it finds the Answer of the given type.
// It also prints all the non-authoritative answers in stdout.
func Resolve(domain string, questionType uint16, options ...Option) string {
	var option Option
	if len(options) > 0 {
		option = options[0]
	} else {
		option = DefaultOption()
	}

	cacheClient, err := cache.NewClient()
	if err != nil {
		fmt.Printf("Failed to create the cache client: %v\n", err)
		os.Exit(1)
	}
	defer cacheClient.Close()

	// Using cache
	if option.UseCache {
		results, err := cacheClient.Get(domain)
		if err != nil {
			fmt.Printf("Failed to get the cached results: %v\n", err)
			os.Exit(1)
		}
		if len(results) > 0 {
			fmt.Printf("Cache hit for %s\n", domain)
			for _, result := range results {
				fmt.Printf("Name: %s\n", result.Name)
				fmt.Printf("Address: %s\n", result.RDataParsed)
			}
			return results[0].RDataParsed
		}
	}

	question := dns.NewQuestion(domain, questionType, dns.ClassIN)
	flag := dns.NewHeaderFlag(false, 0, false, false, false, false, 0, 0).GenerateFlag()
	header := dns.NewHeader(22, flag, 1, 0, 0, 0)
	DNSMessage := dns.NewDNSMessage(*header, []dns.Question{*question})
	var parsedResponse *dns.DNSMessage
	dnsServerIP := dns.RootDNS
	dnsServerPort := dns.RootDNSPort

	for {
		fmt.Printf("Querying %s for %s\n", dnsServerIP, domain)
		// fmt.Printf("DNS Message:\n %+v\n\n", DNSMessage)
		client := NewClient(dnsServerIP, dnsServerPort)
		response, err := client.Query(DNSMessage.ToBytes())
		if err != nil {
			fmt.Printf("Failed to query the DNS server: %v\n", err)
			return ""
		}
		parsedResponse = dns.DNSMessageFromBytes(response)
		// fmt.Printf("parsedResponse:\n %+v\n\n", parsedResponse)
		flags := dns.HeaderFlagFromUint16(parsedResponse.Header.Flags)

		if flags.HasError() {
			fmt.Printf("The DNS server returned an error: %s\n", parsedResponse.Answers[0].RDataParsed)
			os.Exit(1)
		}

		if flags.IsQuery() {
			fmt.Printf("The returned DNS message is not a response.\n")
			os.Exit(1)
		}

		if parsedResponse.Header.ANCount > 0 {
			fmt.Printf("\nNon-authoritative answer:\n")
			if parsedResponse.Answers[0].Type == dns.TypeCNAME {
				fmt.Printf("%s	canonical name = %s.\n", parsedResponse.Answers[0].Name, parsedResponse.Answers[0].RDataParsed)
				Resolve(parsedResponse.Answers[0].RDataParsed, dns.TypeA)
			} else {
				for _, answer := range parsedResponse.Answers {
					fmt.Printf("Name: %s\n", answer.Name)
					fmt.Printf("Address: %s\n", answer.RDataParsed)
					cacheClient.Insert(domain, dns.TypeA, answer.RDataParsed, int(answer.TTL))
				}
			}
			break
		} else if parsedResponse.Header.ARCount > 0 {
			if ip := getRecord(parsedResponse.AdditionalRRs); ip != "" {
				dnsServerIP = ip
			}
			continue
		} else if parsedResponse.Header.NSCount > 0 {
			if nsDomain := getRecord(parsedResponse.AuthorityRRs); nsDomain != "" {
				dnsServerIP = Resolve(nsDomain, dns.TypeA)
			}
		} else {
			fmt.Printf("No answers found for %s\n", domain)
			os.Exit(1)
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
