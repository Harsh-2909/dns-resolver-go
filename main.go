package main

import (
	"dns-resolver-go/types"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		// fmt.Println("Usage: go run main.go <port>")
		fmt.Println("Usage: go run main.go <domain>")
		os.Exit(1)
	}
	domain := os.Args[1]

	question := types.NewQuestion(domain, 1, 1)
	generateFlag := types.GenerateFlag(0, 0, 0, 0, 1, 0, 0, 0)
	header := types.NewHeader(22, generateFlag, 1, 0, 0, 0)
	DNSMessage := types.DNSMessage{
		Header: *header,
		Questions: []types.Question{
			*question,
		},
	}

	fmt.Printf("Domain: %s\n", domain)
	fmt.Printf("DNS Message: %+v\n\n", DNSMessage)
	fmt.Printf("DNS Message in Bytes: %+v\n", DNSMessage.ToBytes())
	fmt.Printf("DNS Message in Hex: %s\n", hex.EncodeToString(DNSMessage.ToBytes()))
	fmt.Printf("Header in Hex: %s\n", hex.EncodeToString(header.ToBytes()))
	fmt.Printf("Header in Hex: %s\n", hex.EncodeToString(question.ToBytes()))

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
