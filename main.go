package main

import (
	"dns-resolver-go/dns"
	"dns-resolver-go/network"
	"fmt"
	"os"
	"strings"
)

func main() {
	argLen := len(os.Args)
	if argLen < 2 {
		fmt.Println("Usage: go run main.go <domain> [OPTIONS]")
		fmt.Println("OPTIONS:")
		fmt.Println("  --no-cache: Resolve the domain without using the cache.")
		os.Exit(1)
	}
	domain := os.Args[1]
	userOptions := strings.Join(os.Args[2:], ",")
	option := network.DefaultOption()
	if strings.Contains(userOptions, "--no-cache") {
		option.UseCache = false
	}
	network.Resolve(domain, dns.TypeA, option)
}
