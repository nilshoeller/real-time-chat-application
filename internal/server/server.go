package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/net/websocket"
)

var borderStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("30")).
	Padding(1, 2).
	Width(70).
	Height(10).
	Align(lipgloss.Left)

type Server struct {
	// Registered clients.
	connections map[*websocket.Conn]bool
	mutex       sync.Mutex
	messageChan chan string
}

type MessageData struct {
	ClientName string `json:"clientName"`
	Message    string `json:"message"`
}

func NewServer() *Server {
	return &Server{
		connections: make(map[*websocket.Conn]bool),
		messageChan: make(chan string),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	// fmt.Println("new incoming connection from client:", ws.RemoteAddr())

	s.mutex.Lock()
	s.connections[ws] = true
	s.mutex.Unlock()

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	defer func() {
		s.mutex.Lock()
		delete(s.connections, ws)
		s.mutex.Unlock()
		ws.Close()
	}()

	for {
		buff := make([]byte, 1024)

		n, err := ws.Read(buff)
		if err != nil {
			// Handle EOF and closed connection errors
			if err == io.EOF {
				fmt.Println("client disconnected")
				break
			}
			// Specific handling for closed network connections
			if strings.Contains(err.Error(), "use of closed network connection") {
				fmt.Println("Attempted to read from closed connection")
				break
			}
			fmt.Println("read error:", err) // Log the error for debugging
			break                           // Exit the loop on any other error
		}

		msg := buff[:n]

		var msgData MessageData
		err = json.Unmarshal(msg, &msgData)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			continue
		}

		response := msgData.ClientName + ": " + msgData.Message
		s.messageChan <- response           // Send the message to the Bubble Tea app
		_, err = ws.Write([]byte(response)) // Return the message to the client
		if err != nil {
			fmt.Println("write error:", err)
			return // Exit the loop on write error
		}
	}
}

func (s *Server) Run() {
	fmt.Println("Running on :3000")

	http.Handle("/ws", websocket.Handler(s.handleWS))
	http.ListenAndServe(":3000", nil)

}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

// //////////////////////////////////////////////////////////////////
type (
	errMsg error
)

// //////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////
// BUBBLETEA ////////////////////////////////////////////////////////
type model struct {
	err      error
	server   *Server
	step     int
	messages []string
}

func initialModel() model {

	server := NewServer()
	go server.Run()

	return model{
		server:   server,
		step:     0,
		messages: []string{},
	}
}

// //////////////////////////////////////////////////////////////////
// INIT /////////////////////////////////////////////////////////////
func (m model) Init() tea.Cmd {

	return tea.Batch(
		tea.EnterAltScreen,
		m.listenForMessages(),
	)
}

// //////////////////////////////////////////////////////////////////
// UPDATE ///////////////////////////////////////////////////////////
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case string:
		// Add the received message to the list
		m.messages = append(m.messages, msg)

		// Continue listening for messages
		return m, m.listenForMessages()
	case tea.KeyMsg:
		switch msg.Type {
		// case tea.KeyEnter:
		// 	// Step 1: Starting server
		// 	if m.step == 0 {
		// 		m.step = 1
		// 		return m, nil
		// 	}

		// 	// Step 2: Receiving messages
		// 	if m.step == 1 {

		// 		return m, nil
		// 	}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, cmd
}

// //////////////////////////////////////////////////////////////////
// VIEW /////////////////////////////////////////////////////////////
func (m model) View() string {
	if m.err != nil {
		return borderStyle.Render(
			fmt.Sprintf("An error occurred: %v\n", m.err),
		)
	}

	var viewContent string
	switch m.step {
	case 0:
		messagesView := strings.Join(m.messages, "\n")

		viewContent = fmt.Sprintf(
			"Server is running\n\n%s",
			messagesView,
		)
	case 1:
		messagesView := strings.Join(m.messages, "\n")
		viewContent = messagesView
	default:
		viewContent = "Something went wrong."
	}
	return borderStyle.Render(viewContent)
}

func (m model) View2() string {
	if m.err != nil {
		return borderStyle.Render(fmt.Sprintf("An error occurred: %v\n", m.err))
	}

	var messagesView string
	if len(m.messages) > 0 {
		messagesView = strings.Join(m.messages, "\n")
	} else {
		messagesView = "No messages yet."
	}

	return borderStyle.Render(fmt.Sprintf("Server is running\n\nMessages:\n%s", messagesView))
}

func (m model) listenForMessages() tea.Cmd {
	return func() tea.Msg {
		// Block and wait for a message from the server
		msg := <-m.server.messageChan
		return msg
	}
}
