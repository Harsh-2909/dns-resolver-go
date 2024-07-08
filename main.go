package main

import (
	"fmt"
	"net"
)

func main() {
	addr := &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 3000,
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Printf("Listening on %s\n", addr)

	for {
		buffer := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Received", string(buffer[:n]), "from", addr)
		_, err = conn.WriteToUDP(buffer[1:n], addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
