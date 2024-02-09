package pythonAPI

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

// BSDSocket holds the configuration and connected clients.
type BSDSocket struct {
	listener net.Listener
	clients  map[string]net.Conn // Track connected clients by their address
	mu       sync.Mutex          // Ensure thread-safe access to clients map
}

// NewBSDSocket creates and returns a new BSDSocket instance.
func NewBSDSocket() *BSDSocket {
	return &BSDSocket{
		clients: make(map[string]net.Conn),
	}
}

// OpenSocket listens on the specified port and handles incoming connections.
func (s *BSDSocket) OpenSocket(port string) {
	var err error
	s.listener, err = net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Listening on :%s\n", port)

	go s.acceptConnections() // Use a goroutine to accept connections concurrently
}

// acceptConnections handles incoming client connections.
func (s *BSDSocket) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		s.mu.Lock()
		s.clients[conn.RemoteAddr().String()] = conn
		s.mu.Unlock()

		go s.handleConnection(conn)
	}
}

// handleConnection reads messages from a single client and echoes them back.
func (s *BSDSocket) handleConnection(conn net.Conn) {
	fmt.Printf("Client connected [%s]\n", conn.RemoteAddr().String())
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("Received: %s\n", line)
		_, err := conn.Write([]byte("Echo: " + line + "\n"))
		if err != nil {
			fmt.Println("Error sending echo:", err)
			break
		}
	}

	s.mu.Lock()
	delete(s.clients, conn.RemoteAddr().String()) // Remove client from tracking on disconnect
	s.mu.Unlock()

	fmt.Printf("Client disconnected [%s]\n", conn.RemoteAddr().String())
}

// SendMessage sends a message to all connected clients.
func (s *BSDSocket) SendMessage(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for addr, conn := range s.clients {
		_, err := conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Printf("Error sending message to %s: %s\n", addr, err)
			continue
		}
	}
}
