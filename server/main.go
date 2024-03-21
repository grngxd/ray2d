package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"shared/networking/messages"
	"syscall"
	"time"
)

var Connections map[string]*Connection

type Connection struct {
	addr      *net.UDPAddr
	heartbeat time.Time
}

func main() {
	Connections = make(map[string]*Connection)
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
		for _, c := range Connections {
			message := messages.New("closed", messages.Closed)
			data := message.Serialize()
			_, err := conn.WriteToUDP(data, c.addr)
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
			Connections[addr.String()] = &Connection{
				addr:      addr,
				heartbeat: time.Now(),
			}

			log.Println("Client connected:", addr)
			for _, client := range Connections {
				fmt.Printf("%v: %v\n", client, addr)
			}
		case messages.Disconnect:
			delete(Connections, addr.String())
			log.Println("Client disconnected:", addr)
		case messages.Heartbeat:
			log.Println("Received heartbeat from client:", addr)
			Connections[addr.String()].heartbeat = time.Now()
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

		go func() {
			ticker := time.NewTicker(1 * time.Second)
			for range ticker.C {
				for addr, c := range Connections {
					if time.Since(c.heartbeat) > 5*time.Second {
						log.Println("Client timed out:", addr)
						delete(Connections, addr)
						continue
					}
				}
			}
		}()
	}
}
