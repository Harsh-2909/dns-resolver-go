package main

import (
	"dns-resolver-go/network"
	"dns-resolver-go/types"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		// fmt.Println("Usage: go run main.go <port>")
		fmt.Println("Usage: go run main.go <domain>")
		os.Exit(1)
	}
	domain := os.Args[1]

	question := types.NewQuestion(domain, 1, 1)
	recursionFlag := types.GenerateFlag(0, 0, 0, 0, 1, 0, 0, 0)
	header := types.NewHeader(22, recursionFlag, 1, 0, 0, 0)
	DNSMessage := types.DNSMessage{
		Header: *header,
		Questions: []types.Question{
			*question,
		},
	}

	fmt.Printf("DNS Message in Bytes: %+v\n", DNSMessage.ToBytes())
	fmt.Printf("DNS Message in Hex: %s\n", hex.EncodeToString(DNSMessage.ToBytes()))

	dnsServer := "8.8.8.8:53"

	// Create a UDP connection
	conn, err := net.Dial("udp", dnsServer)
	if err != nil {
		fmt.Printf("Failed to connect to the DNS server: %v\n", err)
		return
	}
	defer conn.Close()

	// Set a timeout for the connection
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// Send a message
	message := DNSMessage.ToBytes()
	_, err = conn.Write(message)
	if err != nil {
		fmt.Printf("Failed to send the DNS message: %v\n", err)
		return
	}

	// Receive the response
	buf := make([]byte, 1024)
	// Read the response
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Failed to read the response: %v\n", err)
		return
	}
	response := buf[:n]

	fmt.Printf("Response: %v\n", response)
	fmt.Printf("Response in Hex: %v\n", hex.EncodeToString(response))

	// Check if the response ID matches the request ID
	if !network.IDMatcher(message[:2], response[:2]) {
		fmt.Println("The response ID does not match the request ID")
		return
	} else {
		fmt.Println("The response ID matches the request ID")
	}

	// addr := &net.UDPAddr{
	// 	IP:   net.IPv4(127, 0, 0, 1),
	// 	Port: 3000,
	// }
	// conn, err := net.ListenUDP("udp", addr)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer conn.Close()
	// fmt.Printf("Listening on %s\n", addr)

	// for {
	// 	buffer := make([]byte, 1024)
	// 	n, addr, err := conn.ReadFromUDP(buffer)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	fmt.Println("Received", string(buffer[:n]), "from", addr)
	// 	_, err = conn.WriteToUDP(buffer[1:n], addr)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// }
}
