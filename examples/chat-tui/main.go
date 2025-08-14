// FSK Chat TUI - Terminal User Interface for FSK-based chat
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gleicon/go-fsk/fsk"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	channelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575"))

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	ownMessageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)

// Message represents a chat message
type Message struct {
	Channel   int
	Username  string
	Content   string
	Timestamp time.Time
	IsOwn     bool
}

// Global message channel for FSK callbacks to communicate with Bubble Tea
var messageChan = make(chan messageReceivedMsg, 100)

// Model holds the application state
type Model struct {
	username    string
	chat        *fsk.MultiChannelChat
	channels    []fsk.ChannelConfig
	activeChans []int
	currentChan int
	messages    []Message
	input       string
	mode        string // "channel", "chat", "help"
	width       int
	height      int
	ready       bool
}

// Messages for Bubble Tea
type tickMsg struct{}
type messageReceivedMsg struct {
	channelID int
	username  string
	message   string
}

func tickEvery() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func listenForMessages() tea.Cmd {
	return func() tea.Msg {
		select {
		case msg := <-messageChan:
			return msg
		default:
			return nil
		}
	}
}

func initialModel() Model {
	channels := fsk.PredefinedChannels()
	username := "User"
	if len(os.Args) > 1 {
		username = os.Args[1]
	}

	return Model{
		username:    username,
		channels:    channels,
		currentChan: 1,
		mode:        "channel",
		messages:    []Message{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickEvery(),
		listenForMessages(),
		tea.WindowSize(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		// Initialize chat system once we have the window
		if m.chat == nil {
			m.chat = fsk.NewMultiChannelChat(m.username, func(channelID int, username, message string) {
				// Send received message to the global channel
				select {
				case messageChan <- messageReceivedMsg{
					channelID: channelID,
					username:  username,
					message:   message,
				}:
				default:
					// Channel full, drop message
				}
			})
		}
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case "channel":
			return m.updateChannelSelect(msg)
		case "chat":
			return m.updateChat(msg)
		case "help":
			return m.updateHelp(msg)
		}

	case tickMsg:
		return m, tea.Batch(tickEvery(), listenForMessages())

	case messageReceivedMsg:
		if msg.channelID != 0 { // Only add non-nil messages
			m.messages = append(m.messages, Message{
				Channel:   msg.channelID,
				Username:  msg.username,
				Content:   msg.message,
				Timestamp: time.Now(),
				IsOwn:     false,
			})
		}
		return m, listenForMessages()
	}

	return m, nil
}

func (m Model) updateChannelSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.currentChan > 1 {
			m.currentChan--
		}
	case "down", "j":
		if m.currentChan < len(m.channels) {
			m.currentChan++
		}
	case "enter", " ":
		// Join selected channel
		if m.chat != nil {
			channelConfig := m.channels[m.currentChan-1]
			err := m.chat.JoinChannel(channelConfig, 2, 100)
			if err != nil {
				m.messages = append(m.messages, Message{
					Channel:   0,
					Username:  "System",
					Content:   fmt.Sprintf("Failed to join %s: %v", channelConfig.Name, err),
					Timestamp: time.Now(),
					IsOwn:     false,
				})
			} else {
				m.activeChans = append(m.activeChans, channelConfig.ID)
				m.messages = append(m.messages, Message{
					Channel:   0,
					Username:  "System",
					Content:   fmt.Sprintf("Joined %s", channelConfig.Name),
					Timestamp: time.Now(),
					IsOwn:     false,
				})
			}
		}
	case "l":
		// Leave current channel
		if len(m.activeChans) > 0 && m.chat != nil {
			channelID := m.channels[m.currentChan-1].ID
			for i, activeID := range m.activeChans {
				if activeID == channelID {
					err := m.chat.LeaveChannel(channelID)
					if err == nil {
						// Remove from active channels
						m.activeChans = append(m.activeChans[:i], m.activeChans[i+1:]...)
						m.messages = append(m.messages, Message{
							Channel:   0,
							Username:  "System",
							Content:   fmt.Sprintf("Left channel %d", channelID),
							Timestamp: time.Now(),
							IsOwn:     false,
						})
					}
					break
				}
			}
		}
	case "c":
		if len(m.activeChans) > 0 {
			m.mode = "chat"
		} else {
			m.messages = append(m.messages, Message{
				Channel:   0,
				Username:  "System",
				Content:   "Join at least one channel first!",
				Timestamp: time.Now(),
				IsOwn:     false,
			})
		}
	case "h", "?":
		m.mode = "help"
	}
	return m, nil
}

func (m Model) updateChat(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = "channel"
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		if m.input != "" && m.chat != nil {
			// Send message to all active channels
			err := m.chat.BroadcastMessage(m.input)
			if err == nil {
				m.messages = append(m.messages, Message{
					Channel:   0,
					Username:  m.username,
					Content:   m.input,
					Timestamp: time.Now(),
					IsOwn:     true,
				})
			} else {
				m.messages = append(m.messages, Message{
					Channel:   0,
					Username:  "System",
					Content:   fmt.Sprintf("Send error: %v", err),
					Timestamp: time.Now(),
					IsOwn:     false,
				})
			}
			m.input = ""
		}
	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m Model) updateHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.mode = "channel"
	case "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) View() string {
	if !m.ready {
		return "Loading FSK Chat TUI..."
	}

	switch m.mode {
	case "channel":
		return m.viewChannelSelect()
	case "chat":
		return m.viewChat()
	case "help":
		return m.viewHelp()
	}
	return ""
}

func (m Model) viewChannelSelect() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("FSK Ultrasonic Chat"))
	b.WriteString("\n\n")

	b.WriteString("Select Frequency Channels:\n\n")

	for i, channel := range m.channels {
		cursor := " "
		if i+1 == m.currentChan {
			cursor = ">"
		}

		status := ""
		for _, activeID := range m.activeChans {
			if activeID == channel.ID {
				status = " [JOINED]"
				break
			}
		}

		line := fmt.Sprintf("%s %d. %s (%.0fkHz)%s",
			cursor, i+1, channel.Name, channel.BaseFreq/1000, status)

		if i+1 == m.currentChan {
			b.WriteString(channelStyle.Render(line))
		} else {
			b.WriteString(line)
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓: navigate • enter: join • l: leave • c: chat mode • h: help • q: quit"))

	// Show recent messages
	if len(m.messages) > 0 {
		b.WriteString("\n\nRecent Messages:\n")
		start := len(m.messages) - 5
		if start < 0 {
			start = 0
		}
		for i := start; i < len(m.messages); i++ {
			msg := m.messages[i]
			timestamp := msg.Timestamp.Format("15:04:05")
			if msg.IsOwn {
				b.WriteString(ownMessageStyle.Render(fmt.Sprintf("[%s] %s: %s", timestamp, msg.Username, msg.Content)))
			} else if msg.Username == "System" {
				b.WriteString(statusStyle.Render(fmt.Sprintf("[%s] %s", timestamp, msg.Content)))
			} else {
				b.WriteString(messageStyle.Render(fmt.Sprintf("[%s] %s: %s", timestamp, msg.Username, msg.Content)))
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (m Model) viewChat() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(fmt.Sprintf("FSK Chat - %s", m.username)))
	b.WriteString("\n\n")

	// Active channels
	b.WriteString("Active Channels: ")
	if len(m.activeChans) == 0 {
		b.WriteString(statusStyle.Render("None"))
	} else {
		for i, chanID := range m.activeChans {
			if i > 0 {
				b.WriteString(", ")
			}
			for _, ch := range m.channels {
				if ch.ID == chanID {
					b.WriteString(channelStyle.Render(fmt.Sprintf("Ch%d(%.0fkHz)", chanID, ch.BaseFreq/1000)))
					break
				}
			}
		}
	}
	b.WriteString("\n\n")

	// Messages area
	messageHeight := m.height - 10
	if messageHeight < 5 {
		messageHeight = 5
	}

	start := len(m.messages) - messageHeight
	if start < 0 {
		start = 0
	}

	for i := start; i < len(m.messages); i++ {
		msg := m.messages[i]
		timestamp := msg.Timestamp.Format("15:04:05")

		if msg.IsOwn {
			b.WriteString(ownMessageStyle.Render(fmt.Sprintf("[%s] You: %s", timestamp, msg.Content)))
		} else if msg.Username == "System" {
			b.WriteString(statusStyle.Render(fmt.Sprintf("[%s] %s", timestamp, msg.Content)))
		} else {
			channelInfo := ""
			if msg.Channel > 0 {
				channelInfo = fmt.Sprintf(" (Ch%d)", msg.Channel)
			}
			b.WriteString(messageStyle.Render(fmt.Sprintf("[%s] %s%s: %s", timestamp, msg.Username, channelInfo, msg.Content)))
		}
		b.WriteString("\n")
	}

	// Input area
	b.WriteString("\n")
	b.WriteString(inputStyle.Render(fmt.Sprintf("Message: %s", m.input)))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("enter: send • esc: back to channels • ctrl+c: quit"))

	return b.String()
}

func (m Model) viewHelp() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("FSK Ultrasonic Chat - Help"))
	b.WriteString("\n\n")

	b.WriteString("FREQUENCY CHANNELS:\n")
	b.WriteString("This chat uses ultrasonic frequencies (18-26kHz) for communication.\n")
	b.WriteString("Each channel uses a different frequency band to avoid interference.\n\n")

	b.WriteString("CHANNEL TYPES:\n")
	b.WriteString("• Broadcast: All users on same channel can hear each other\n")
	b.WriteString("• Point-to-Point: Use different TX/RX frequencies for private chat\n")
	b.WriteString("• Multi-Channel: Join multiple channels simultaneously\n\n")

	b.WriteString("FREQUENCY MIXING:\n")
	b.WriteString("• Same frequencies: Signals collide and corrupt\n")
	b.WriteString("• Separate frequencies: Clean communication\n")
	b.WriteString("• Overlapping frequencies: Partial interference\n\n")

	b.WriteString("CONTROLS:\n")
	b.WriteString("Channel Selection:\n")
	b.WriteString("  ↑/↓ or j/k: Navigate channels\n")
	b.WriteString("  enter/space: Join selected channel\n")
	b.WriteString("  l: Leave selected channel\n")
	b.WriteString("  c: Switch to chat mode\n")
	b.WriteString("  h or ?: Show this help\n")
	b.WriteString("  q: Quit application\n\n")

	b.WriteString("Chat Mode:\n")
	b.WriteString("  Type message and press enter to send\n")
	b.WriteString("  Messages broadcast to all joined channels\n")
	b.WriteString("  esc: Return to channel selection\n")
	b.WriteString("  ctrl+c: Quit application\n\n")

	b.WriteString("TECHNICAL NOTES:\n")
	b.WriteString("• Ultrasonic signals (>20kHz) are inaudible to humans\n")
	b.WriteString("• Each channel uses 4-FSK modulation (2 bits per symbol)\n")
	b.WriteString("• Baud rate: 100 symbols/sec (200 bits/sec)\n")
	b.WriteString("• Range: Limited by speaker/microphone quality\n\n")

	b.WriteString(helpStyle.Render("esc or q: back to channels"))

	return b.String()
}

func main() {
	model := initialModel()

	// Start the Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}

	// Cleanup
	if model.chat != nil {
		model.chat.Close()
	}
}
