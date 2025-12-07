package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ToffaKrtek/go-tui-openvpn-client/internal/cmd"
	"github.com/ToffaKrtek/go-tui-openvpn-client/types"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle          = lipgloss.NewStyle().Margin(1, 2)
	inputStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	inputStyleBlurred = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle           = lipgloss.NewStyle()
	configs           []types.Item
	configsList       []list.Item
	activeSession     string
	blurredButton     = fmt.Sprintf("[ %s ]", inputStyleBlurred.Render("Submit"))
	focusedButton     = inputStyle.Render(" [ Submit ]")
)

type modal struct {
	focusIndex int
	fields     []textinput.Model
	cursorMode cursor.Mode
}

type model struct {
	list      list.Model
	showModal bool
	modal     modal
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.showModal {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case " ":
				index := m.list.Index()
				if index < len(configs) && index >= 0 {
					name := configs[index].Name
					activeDeactive(name)
				}
			case "d":
				index := m.list.Index()
				if index < len(configs) && index >= 0 {
					name := configs[index].Name
					cmd.DeleteConfig(name)
					setListConfigs()
					return m, m.list.SetItems(configsList)
				}
			case "i":
				m.showModal = true
				// case "d":
				// TODO:: delete
			case "r":
				setListConfigs()
				setActiveSessions()
				return m, m.list.SetItems(configsList)
			}
		case tea.WindowSizeMsg:
			h, v := docStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)

		}
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "ctrl+r":
			m.modal.cursorMode++
			if m.modal.cursorMode > cursor.CursorHide {
				m.modal.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.modal.fields))
			for i := range m.modal.fields {
				cmds[i] = m.modal.fields[i].Cursor.SetMode(m.modal.cursorMode)
			}
			return m, tea.Batch(cmds...)

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.modal.focusIndex == len(m.modal.fields) {
				m.showModal = false
				path := m.modal.fields[0].Value()
				name := m.modal.fields[1].Value()
				fmt.Println(path, name)
				if len(path) > 0 && len(name) > 0 {
					err := cmd.ImportConfig(name, path)
					if err == nil {
						m.modal.fields[0].SetValue("")
						m.modal.fields[1].SetValue("")
						setListConfigs()
						return m, m.list.SetItems(configsList)
					}
				}
			}

			if s == "up" || s == "shift+tab" {
				m.modal.focusIndex--
			} else {
				m.modal.focusIndex++
			}

			if m.modal.focusIndex > len(m.modal.fields) {
				m.modal.focusIndex = 0
			} else if m.modal.focusIndex < 0 {
				m.modal.focusIndex = len(m.modal.fields)
			}

			cmds := make([]tea.Cmd, len(m.modal.fields))
			for i := 0; i <= len(m.modal.fields)-1; i++ {
				if i == m.modal.focusIndex {
					// Set focused state
					cmds[i] = m.modal.fields[i].Focus()
					m.modal.fields[i].PromptStyle = inputStyle
					m.modal.fields[i].TextStyle = inputStyle
					continue
				}
				// Remove focused state
				m.modal.fields[i].Blur()
				m.modal.fields[i].PromptStyle = noStyle
				m.modal.fields[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.modal.fields))
	for i := range m.modal.fields {
		m.modal.fields[i], cmds[i] = m.modal.fields[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.showModal {
		return docStyle.Render(activeSession, m.list.View())
	}
	var b strings.Builder

	for i := range m.modal.fields {
		b.WriteString(m.modal.fields[i].View())
		if i < len(m.modal.fields)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.modal.focusIndex == len(m.modal.fields) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}

func activeDeactive(name string) {
	activeSessions, err := cmd.GetSession()
	if err != nil {
		return
	}
	for _, s := range activeSessions {
		if s.Name == name {
			cmd.DisconnectSession(name)
			setActiveSessions()
			return
		}
	}
	cmd.ActiveConfig(name)
	setActiveSessions()
}

func setListConfigs() {
	var err error
	configsList = []list.Item{}
	configs, err = cmd.GetConfigs()
	if err != nil {
		return
	}
	listItems := []list.Item{}
	for _, cfg := range configs {
		listItems = append(listItems, cfg)
	}
	configsList = listItems
}

func setActiveSessions() {
	sessions, err := cmd.GetSession()
	if err != nil {
		activeSession = "ERROR\n"
		return
	}
	if len(sessions) < 1 {
		activeSession = "NO ACTIVE\n"
		return
	}
	activeSession = fmt.Sprintf("%s: %s\n", sessions[0].Name, sessions[0].Text)
}

func main() {
	setActiveSessions()
	setListConfigs()

	fields := make([]textinput.Model, 2)
	pathInput := textinput.New()
	pathInput.Cursor.Style = inputStyle
	pathInput.PromptStyle = inputStyle
	pathInput.TextStyle = inputStyle
	pathInput.CharLimit = 64
	pathInput.Placeholder = "Path to file"
	pathInput.Focus()
	fields[0] = pathInput
	nameInput := textinput.New()
	nameInput.Cursor.Style = inputStyle
	pathInput.CharLimit = 64
	pathInput.Placeholder = "Config name"
	fields[1] = nameInput
	m := model{
		list:      list.New(configsList, list.NewDefaultDelegate(), 0, 0),
		showModal: false,
		modal: modal{
			focusIndex: 0,
			fields:     fields,
		},
	}
	m.list.Title = ""
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
