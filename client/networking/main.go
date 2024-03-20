package networking

import (
	"fmt"
	"log"
	"net"
	"shared/networking/messages"
	"time"
)

type Connection struct {
	connection *net.UDPConn
	connected  bool
}

func (c *Connection) Connect() {
	// Handle the connection logic here
	c.connected = true
	connect := messages.New("connect", messages.Connect)
	data := connect.Serialize()
	_, err := c.connection.Write(data)
	if err != nil {
		fmt.Println("Failed to send packet to server:", err)
		return
	}
	go c.sendHeartbeat()    // Start sending heartbeats to the server
	go c.receiveHeartbeat() // Start receiving heartbeats from the server
	for c.connected {

		// Create a buffer to hold incoming data.
		buffer := make([]byte, 1024)
		n, addr, err := c.connection.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Failed to read from connection:", err)
			continue
		}
		// serialise the data to a Message struct
		message := messages.Deserialize(buffer[:n])
		if message == nil {
			log.Println("Failed to deserialise message")
			continue
		}

		switch message.Type {
		case messages.Closed:
			log.Println("Server closed the connection:", addr)
			c.Disconnect()
			return
		case messages.Heartbeat:
			fmt.Println("Received heartbeat response from server")
		default:
			log.Println("Unknown message type received from client:", addr)
		}

	}
}
func (c *Connection) sendHeartbeat() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		heartbeat := messages.New("heartbeat", messages.Heartbeat)
		data := heartbeat.Serialize()
		_, err := c.connection.Write(data)
		if err != nil {
			fmt.Println("Failed to send heartbeat packet to server:", err)
			c.Disconnect()
			return
		}
	}
}

func (c *Connection) receiveHeartbeat() {
	for c.connected {
		buf := make([]byte, 1024)
		n, _, err := c.connection.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Failed to read heartbeat response from server:", err)
			c.Disconnect()
			return
		}

		response := messages.Deserialize(buf[:n])
		if response == nil {
			fmt.Println("Failed to deserialise heartbeat response")
			continue
		}

		if response.Type == messages.Heartbeat {
			fmt.Println("Received heartbeat response from server")
		} else {
			fmt.Println("Unknown message type received from server", response.Type)
		}
	}
}

func (c *Connection) Disconnect() {
	if c.connection == nil {
		return
	}

	disconnect := messages.New("disconnect", messages.Disconnect)
	data := disconnect.Serialize()
	_, err := c.connection.Write(data)
	if err != nil {
		fmt.Println("Failed to send disconnect packet to server:", err)
		return
	}

	err = c.connection.Close()
	if err != nil {
		fmt.Println("Failed to disconnect:", err)
		return
	}

	c.connected = false
	fmt.Println("Disconnected successfully")
}

func NewConnection(addr string) *Connection {
	c := &Connection{}
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		// fatal
		log.Fatal("Failed to resolve server address:", err)
		return nil
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
		return nil
	}

	c.connection = conn
	return c
}
