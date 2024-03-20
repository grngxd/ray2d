package main

import (
	"log"
	"net"
	"server/messages"
)

func main() {
	// Create a UDP address to listen on
	address, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		log.Fatalf("Error resolving address: %v", err)
	}

	// Create a UDP connection
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Error closing connection: %v", err)
		}
	}()

	log.Println("UDP server is listening on", address.String())

	// Create a buffer to store incoming data
	buffer := make([]byte, 1024)

	for {
		// Read binary from the connection into the buffer
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatalf("Error reading from connection: %v", err)
		}

		// serialise the data to a Message struct
		message := messages.Deserialize(buffer[:n])
		if message == nil {
			log.Println("Failed to deserialise message")
			continue
		}

		log.Printf("Received message from %v...\nType: %v\nContent: %v\n", addr, message.Type, message.Content)
	}
}
