package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/yagnikpt/boomtypr/internal/typing"
	"github.com/yagnikpt/boomtypr/internal/ui"
	"github.com/yagnikpt/boomtypr/internal/wordlist"
)

var version = "dev"

func main() {
	var showVersion bool
	var showHelp bool

	flag.BoolVar(&showVersion, "version", false, "show version number")
	flag.BoolVar(&showVersion, "v", false, "show version number (shorthand)")

	flag.BoolVar(&showHelp, "help", false, "show help text")
	flag.BoolVar(&showHelp, "h", false, "show help text (shorthand)")

	flag.Usage = func() {
		fmt.Println("A terminal-based typing speed test.\n\nUsage: boomtypr\n\nTUI Keybinds:\n- tab: Toggle Modes\n- up/right: Increase time/words\n- down/left: Decrease time/words\n- esc: Quit")
	}

	flag.Parse()

	if showVersion {
		fmt.Println(version)
		return
	}

	if showHelp {
		flag.Usage()
		return
	}

	if version == "dev" {
		dir, _ := os.Getwd()
		logFile := filepath.Join(dir, "debug.log")
		fLog, _ := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		log.SetOutput(fLog)
	}

	wl := wordlist.New()
	p := tea.NewProgram(ui.NewModel(wl, typing.ModeTime, 30*time.Second, 50))

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
