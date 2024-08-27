package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/net/websocket"
)

var serverURL string = "ws://localhost:3000/ws"

type Client struct {
	client_name string
	server_url  string
	ws          *websocket.Conn
}

func NewClient(client_name string, server_url string) *Client {
	return &Client{
		client_name: client_name,
		server_url:  server_url,
	}
}

func (c *Client) Connect() error {
	// Connect to the WebSocket server
	var err error
	c.ws, err = websocket.Dial(c.server_url, "", c.client_name)
	if err != nil {
		return err
	}

	fmt.Println("Client connected to the server.")
	return nil
}

func (c *Client) SendMessage(message string) error {
	_, err := c.ws.Write([]byte(message))
	if err != nil {
		return err
	}
	fmt.Println("Message sent to the server:", message)
	return nil
}

func (c *Client) ReceiveMessage() (string, error) {
	buff := make([]byte, 1024)
	n, err := c.ws.Read(buff)
	if err != nil {
		return "", err
	}
	return string(buff[:n]), nil
}

func (c *Client) Close() {
	c.ws.Close()
}

func (c *Client) Run() {
	c.Connect()
	c.SendMessage("Hello Server, this is Nils.")
	c.ReceiveMessage()
	c.Close()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	// newClient := NewClient("http://this-is-a-new-client:8000/", serverURL)
	// newClient.Run()
}

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	err       error
	client    *Client
	step      int
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter your client name..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
		client:    nil,
		step:      0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Step 1: Entering client name
			if m.step == 0 {
				clientName := m.textInput.Value()
				clientURI := "http://" + clientName + ":8000/"
				m.client = NewClient(clientURI, serverURL)

				if err := m.client.Connect(); err != nil {
					m.err = err
					return m, tea.Quit
				}

				m.step = 1
				m.textInput.SetValue("")
				m.textInput.Placeholder = "Enter a message to send..."
				return m, nil
			}

			// Step 2: Sending message
			if m.step == 1 {
				message := m.textInput.Value()

				if message == "" {
					return m, tea.Quit // safety: todo -> return better and continue asking for a message
				}

				err := m.client.SendMessage(message)
				if err != nil {
					m.err = err
				}

				// Optional: receive a message from the server after sending
				receivedMsg, err := m.client.ReceiveMessage()
				if err != nil {
					m.err = err
				} else {
					fmt.Println("Message received from server:", receivedMsg)
				}

				return m, tea.Quit
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("An error occurred: %v\n", m.err)
	}

	switch m.step {
	case 0:
		return fmt.Sprintf(
			"Enter your client name:\n\n%s\n\n%s",
			m.textInput.View(),
			"(esc to quit)",
		)
	case 1:
		return fmt.Sprintf(
			"Send a message to the server!\n\n%s\n\n%s",
			m.textInput.View(),
			"(esc to quit)",
		)
	default:
		return "Something went wrong."
	}
}
