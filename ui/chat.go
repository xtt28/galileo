package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type message struct {
	Sender string
	Content string
}

type ChatModel struct {
	Messages []message
}

func InitialModel() ChatModel {
	return ChatModel {
		Messages: []message{{"User", "hi"}},
	}
}

func (m ChatModel) Init() tea.Cmd {
	return nil
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.Messages = append(m.Messages, message{"AI", "hello"})
	return m, nil
}

func (m ChatModel) View() string {
	s := "Here are our chat messages."

	for _, msg := range m.Messages {
		s += fmt.Sprintf("%s: %s\n", msg.Sender, msg.Content)
	}
	
	return s
}

