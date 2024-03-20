package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"shared/networking/messages"
	"syscall"
)

var Connections map[string]*net.UDPAddr

func main() {
	Connections = make(map[string]*net.UDPAddr)
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

	// Create a channel to receive SIGINT signal
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// Send "closed" message to all connected clients
		for _, addr := range Connections {
			message := messages.New("closed", messages.Closed)
			data := message.Serialize()
			_, err := conn.WriteToUDP(data, addr)
			if err != nil {
				log.Println("Failed to send closed message to client:", err)
			}
		}
		os.Exit(0)
	}()

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

		//log.Printf("Received message from %v...\nType: %v\nContent: %v\n", addr, message.Type, message.Content)
		switch message.Type {
		case messages.Connect:
			Connections[addr.String()] = addr
			log.Println("Client connected:", addr)
		case messages.Disconnect:
			delete(Connections, addr.String())
			log.Println("Client disconnected:", addr)
		case messages.Heartbeat:
			log.Println("Received heartbeat from client:", addr)
			// Send a response back to the client
			response := messages.New("heartbeat response", messages.Heartbeat)
			data := response.Serialize()
			_, err := conn.WriteToUDP(data, addr)
			if err != nil {
				log.Println("Failed to send heartbeat response to client:", err)
			}
		default:
			log.Println("Unknown message type received from client:", addr, message.Type)
		}
	}
}
