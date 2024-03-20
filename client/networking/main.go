package networking

import (
	"fmt"
	"log"
	"net"
	"ray2d/networking/messages"
)

type Connection struct {
	connection *net.UDPConn
}

func (c *Connection) Connect() {
	// Handle the connection logic here
	connect := messages.New("connect", messages.Connect)
	data := connect.Serialize()
	_, err := c.connection.Write(data)
	if err != nil {
		fmt.Println("Failed to send packet to server:", err)
		return
	}

	for {
		// Create a buffer to hold incoming data.
		buf := make([]byte, 1024)
		n, addr, err := c.connection.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Failed to read from server:", err)
			return
		}
		fmt.Println("Received packet from server:", string(buf[:n]), " from ", addr)
	}
}

func (c *Connection) Disconnect() {
	if c.connection == nil {
		fmt.Println("No connection to disconnect")
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
