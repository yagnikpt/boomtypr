package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type blinkMsg struct{}
type resumeBlinkMsg struct {
	id int
}

func blinkCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg {
		return blinkMsg{}
	})
}

func resumeBlinkCmd(id int) tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg {
		return resumeBlinkMsg{id: id}
	})
}
