package tui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/archesai/archesai/pkg/llm"
)

// Model represents the TUI application state.
type Model struct {
	chatClient    llm.ChatClient
	session       *llm.ChatSession
	messages      []Message
	input         string
	width         int
	height        int
	isProcessing  bool
	err           error
	selectedAgent int
	personas      []*llm.ChatPersona
	showAgentList bool
	style         *Styles
}

// Message represents a chat message.
type Message struct {
	Role    string
	Content string
	Agent   string
}

// Styles holds all the styling configurations.
type Styles struct {
	Title      lipgloss.Style
	Agent      lipgloss.Style
	User       lipgloss.Style
	Assistant  lipgloss.Style
	System     lipgloss.Style
	Input      lipgloss.Style
	Error      lipgloss.Style
	Processing lipgloss.Style
	Border     lipgloss.Style
	AgentList  lipgloss.Style
	Selected   lipgloss.Style
}

// NewStyles creates default styles.
func NewStyles() *Styles {
	return &Styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			Background(lipgloss.Color("63")).
			Padding(0, 2).
			Margin(1),

		Agent: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true),

		User: lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Bold(true).
			Align(lipgloss.Right),

		Assistant: lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")),

		System: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true),

		Input: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1),

		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true),

		Processing: lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Italic(true),

		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")),

		AgentList: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("99")).
			Padding(1),

		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Background(lipgloss.Color("235")).
			Bold(true),
	}
}

// New creates a new TUI model.
func New(chatClient llm.ChatClient, personas []*llm.ChatPersona) Model {
	var session *llm.ChatSession
	if len(personas) > 0 {
		session = chatClient.NewSession(personas[0])
	}

	return Model{
		chatClient: chatClient,
		session:    session,
		personas:   personas,
		messages:   []Message{},
		style:      NewStyles(),
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.addSystemMessage("Welcome to Arches TUI! Press Tab to switch agents, Ctrl+C to quit."),
	)
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case responseMsg:
		m.isProcessing = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		agentName := ""
		if m.session != nil && m.session.Persona != nil {
			agentName = m.session.Persona.Name
		}
		m.messages = append(m.messages, Message{
			Role:    "assistant",
			Content: msg.content,
			Agent:   agentName,
		})
		return m, nil

	case systemMsg:
		m.messages = append(m.messages, Message{
			Role:    "system",
			Content: msg.content,
		})
		return m, nil
	}

	return m, nil
}

// View renders the TUI.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	var sections []string

	// Title
	title := m.style.Title.Render("🤖 Arches Agent TUI")
	sections = append(sections, lipgloss.PlaceHorizontal(m.width, lipgloss.Center, title))

	// Current persona indicator
	if m.session != nil && m.session.Persona != nil {
		agentInfo := m.style.Agent.Render(fmt.Sprintf("Current Agent: %s", m.session.Persona.Name))
		sections = append(sections, agentInfo)
	}

	// Agent list (if showing)
	if m.showAgentList {
		sections = append(sections, m.renderAgentList())
	}

	// Messages area
	messagesView := m.renderMessages()
	sections = append(sections, messagesView)

	// Error display
	if m.err != nil {
		errorMsg := m.style.Error.Render(fmt.Sprintf("Error: %v", m.err))
		sections = append(sections, errorMsg)
	}

	// Processing indicator
	if m.isProcessing {
		processing := m.style.Processing.Render("🔄 Processing...")
		sections = append(sections, processing)
	}

	// Input area
	inputPrompt := m.style.User.Render("You: ")
	inputBox := m.style.Input.Render(m.input)
	inputArea := lipgloss.JoinHorizontal(lipgloss.Left, inputPrompt, inputBox)
	sections = append(sections, inputArea)

	// Help text
	help := m.style.System.Render("Tab: Switch Agent | Enter: Send | Ctrl+C: Quit")
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// handleKeyPress processes keyboard input.
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit

	case tea.KeyTab:
		m.showAgentList = !m.showAgentList
		return m, nil

	case tea.KeyUp:
		if m.showAgentList && m.selectedAgent > 0 {
			m.selectedAgent--
		}
		return m, nil

	case tea.KeyDown:
		if m.showAgentList && m.selectedAgent < len(m.personas)-1 {
			m.selectedAgent++
		}
		return m, nil

	case tea.KeyEnter:
		if m.showAgentList {
			// Select persona
			if m.selectedAgent < len(m.personas) {
				m.session = m.chatClient.NewSession(m.personas[m.selectedAgent])
				m.showAgentList = false
				return m, m.addSystemMessage(
					fmt.Sprintf("Switched to agent: %s", m.personas[m.selectedAgent].Name),
				)
			}
		} else if m.input != "" && !m.isProcessing {
			// Send message
			m.messages = append(m.messages, Message{
				Role:    "user",
				Content: m.input,
			})
			m.isProcessing = true
			cmd := m.sendMessage(m.input)
			m.input = ""
			return m, cmd
		}
		return m, nil

	case tea.KeyBackspace:
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
		return m, nil

	default:
		if msg.Type == tea.KeyRunes {
			m.input += string(msg.Runes)
		}
	}

	return m, nil
}

// renderMessages renders the message history.
func (m Model) renderMessages() string {
	var messages []string

	availableHeight := m.height - 10 // Reserve space for UI elements
	if availableHeight < 5 {
		availableHeight = 5
	}

	for _, msg := range m.messages {
		var styled string
		switch msg.Role {
		case "user":
			styled = m.style.User.Render("You: ") + msg.Content
		case "assistant":
			agentLabel := ""
			if msg.Agent != "" {
				agentLabel = fmt.Sprintf("[%s] ", msg.Agent)
			}
			styled = m.style.Agent.Render(agentLabel+"Assistant: ") +
				m.style.Assistant.Render(msg.Content)
		case "system":
			styled = m.style.System.Render("System: " + msg.Content)
		}
		messages = append(messages, styled)
	}

	content := strings.Join(messages, "\n")

	// Create a bordered box for messages
	messageBox := m.style.Border.
		Width(m.width - 4).
		Height(availableHeight).
		Render(content)

	return messageBox
}

// renderAgentList renders the list of available personas.
func (m Model) renderAgentList() string {
	var items []string
	for i, persona := range m.personas {
		item := fmt.Sprintf("%d. %s", i+1, persona.Name)
		if i == m.selectedAgent {
			item = m.style.Selected.Render(item)
		}
		items = append(items, item)
	}

	list := strings.Join(items, "\n")
	return m.style.AgentList.Render("Select Agent:\n" + list)
}

// sendMessage sends a message to the current chat session.
func (m Model) sendMessage(content string) tea.Cmd {
	return func() tea.Msg {
		if m.session == nil || m.chatClient == nil {
			return responseMsg{
				err: fmt.Errorf("no chat session or client initialized"),
			}
		}

		// Send message to chat client
		response, err := m.chatClient.SendMessage(context.Background(), m.session, content)
		if err != nil {
			return responseMsg{err: err}
		}

		return responseMsg{content: response.Content}
	}
}

// addSystemMessage adds a system message.
func (m Model) addSystemMessage(content string) tea.Cmd {
	return func() tea.Msg {
		return systemMsg{content: content}
	}
}

// Message types for tea.Cmd.
type responseMsg struct {
	content string
	err     error
}

type systemMsg struct {
	content string
}
