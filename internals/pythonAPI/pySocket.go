package pythonAPI

import (
	"bufio"
	"fmt"
	"github.com/Nevoral/quadrupot/internals/Robot"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// BSDSocket holds the configuration and connected clients.
type BSDSocket struct {
	listener net.Listener
	clients  map[string]net.Conn // Track connected clients by their address
	mu       sync.Mutex          // Ensure thread-safe access to clients map
	robot    *Robot.Robot
}

// NewBSDSocket creates and returns a new BSDSocket instance.
func NewBSDSocket(r *Robot.Robot) *BSDSocket {
	return &BSDSocket{
		clients: make(map[string]net.Conn),
		robot:   r,
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
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("Client connected [%s]\n", remoteAddr)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("Received from [%s]: %s\n", remoteAddr, line)
		m := s.ParseMessage(line)
		m.ActionsCall(s.robot)
		if _, err := conn.Write(m.Response); err != nil {
			fmt.Printf("Error sending echo to [%s]: %s\n", remoteAddr, err)
			break
		}
	}

	s.mu.Lock()
	delete(s.clients, remoteAddr) // Remove client from tracking on disconnect
	s.mu.Unlock()

	fmt.Printf("Client disconnected [%s]\n", remoteAddr)
}

// SendMessage sends a message to specified clients or all connected clients if no specific clients are provided.
// Waits for a specified duration for connections to be established if none exist.
func (s *BSDSocket) SendMessage(msg string, specificClients []string, waitTime time.Duration) {
	// Wait for at least one connection to be established if none exist.
	start := time.Now()
	for {
		s.mu.Lock()
		if len(s.clients) > 0 || time.Since(start) > waitTime {
			s.mu.Unlock()
			break
		}
		s.mu.Unlock()
		time.Sleep(time.Second) // Check every second.
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Convert specificClients slice to a map for faster lookup, if needed.
	targetClients := make(map[string]bool)
	for _, addr := range specificClients {
		targetClients[addr] = true
	}

	// Send the message as before.
	for addr, conn := range s.clients {
		if len(specificClients) == 0 || targetClients[addr] {
			_, err := conn.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Printf("Error sending message to %s: %s\n", addr, err)
				continue
			}
		}
	}
}

func (s *BSDSocket) ParseMessage(m string) *Message {
	indexMethod := strings.Index(m, "method:") + len("method:")
	indexAction := strings.Index(m, "actions:") + len("actions:")
	indexEndHead := strings.Index(m, ";>")
	indexEndMethod := strings.Index(m, ";actions:")
	indexEndMessage := strings.Index(m, ">>>")

	method := m[indexMethod:indexEndMethod]
	action := strings.Split(m[indexAction:indexEndHead], ";")

	if method == "POST" {
		indexStartBody := strings.Index(m, "><") + len("><")
		body := m[indexStartBody:indexEndMessage]
		return &Message{
			Method:   method,
			Actions:  action,
			Body:     body,
			Response: []byte{byte('<'), byte('<')},
		}
	}
	return &Message{
		Method:   method,
		Actions:  action,
		Body:     "",
		Response: []byte{byte('<'), byte('<')},
	}
}
