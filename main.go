package main

import (
	"dns-resolver-go/dns"
	"dns-resolver-go/network"
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
	network.Resolve(domain, dns.TypeA)
}
