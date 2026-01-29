package typing

import "time"

type Model struct {
	// test config
	Mode       Mode
	Target     []rune
	Duration   time.Duration
	WordsCount int

	RemainingDuration int
	CurrentWord       int

	// runtime state
	StartedAt time.Time
	EndedAt   time.Time
	Done      bool

	// input tracking
	Keystrokes []Keystroke

	// UI
	FocusMode bool
	Err       error
}

func (m *Model) AddKeystroke(r rune, expected rune, backspace bool) {
	ks := Keystroke{
		Rune:      r,
		Expected:  expected,
		Time:      time.Now(),
		Backspace: backspace,
	}
	if !backspace {
		ks.Correct = r == expected
	}
	m.Keystrokes = append(m.Keystrokes, ks)
}

type Keystroke struct {
	Rune      rune
	Expected  rune
	Time      time.Time
	Correct   bool
	Backspace bool
}

type Mode int

const (
	ModeTime Mode = iota
	ModeWords
	ModeZen
)
