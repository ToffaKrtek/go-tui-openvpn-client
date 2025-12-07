package main

import (
	"fmt"
	"os"

	"github.com/ToffaKrtek/go-tui-openvpn-client/internal/cmd"
	"github.com/ToffaKrtek/go-tui-openvpn-client/types"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			// case "i":
			// TODO:: import
			// case "d":
			// TODO:: delete
		}
		// switch msg.Type {
		// case tea.KeySpace:
		// 	index := m.list.Index()
		// 	if index < len(configs) && index >= 0 {
		// 		name := configs[index].Name
		// 		activeDeactive(name)
		// 	}
		// }
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(activeSession, m.list.View())
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

var configs []types.Item
var configsList []list.Item
var activeSession string

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
	m := model{list: list.New(configsList, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = ""
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
