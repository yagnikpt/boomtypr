package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yagnikpt/boomtypr/internal/ui"
	"github.com/yagnikpt/boomtypr/internal/wordlist"
)

func main() {
	dir, _ := os.Getwd()
	logFile := filepath.Join(dir, "debug.log")
	fLog, _ := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	log.SetOutput(fLog)
	wl := wordlist.New()
	p := tea.NewProgram(ui.NewModel(wl), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
