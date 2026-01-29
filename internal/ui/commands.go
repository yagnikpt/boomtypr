package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type startMsg struct{}
type finishMsg struct{}
type endTickerMsg struct{}
type tickerMsg time.Time

func startCmd() tea.Cmd {
	return func() tea.Msg {
		return startMsg{}
	}
}

func finishCmd() tea.Cmd {
	return func() tea.Msg {
		return finishMsg{}
	}
}

func endTickerCmd(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return endTickerMsg{}
	})
}

func tickerCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickerMsg(t)
	})
}
