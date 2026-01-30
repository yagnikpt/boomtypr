package typing

import "time"

type Model struct {
	// test config
	Mode      Mode
	Target    []rune
	Duration  time.Duration
	WordCount int

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

func (m *Model) NextMode(currentMode Mode) {
	switch currentMode {
	case ModeTime:
		m.Mode = ModeWords
	case ModeWords:
		m.Mode = ModeZen
	case ModeZen:
		m.Mode = ModeTime
	}
}

func (m *Model) NextDuration() {
	durations := []time.Duration{15 * time.Second, 30 * time.Second, 60 * time.Second, 120 * time.Second}
	current := m.Duration
	var next time.Duration
	for i, d := range durations {
		if d == current {
			next = durations[(i+1)%len(durations)]
			break
		}
	}
	m.Duration = next
	m.RemainingDuration = int(next.Seconds())
}

func (m *Model) NextWordCount() {
	counts := []int{10, 25, 50, 100, 200}
	current := m.WordCount
	var next int
	for i, c := range counts {
		if c == current {
			next = counts[(i+1)%len(counts)]
			break
		}
	}
	m.WordCount = next
}

func (m *Model) PrevDuration() {
	durations := []time.Duration{15 * time.Second, 30 * time.Second, 60 * time.Second, 120 * time.Second}
	current := m.Duration
	var prev time.Duration
	for i, d := range durations {
		if d == current {
			if i == 0 {
				prev = durations[len(durations)-1]
			} else {
				prev = durations[i-1]
			}
			break
		}
	}
	m.Duration = prev
	m.RemainingDuration = int(prev.Seconds())
}

func (m *Model) PrevWordCount() {
	counts := []int{10, 25, 50, 100, 200}
	current := m.WordCount
	var prev int
	for i, c := range counts {
		if c == current {
			if i == 0 {
				prev = counts[len(counts)-1]
			} else {
				prev = counts[i-1]
			}
			break
		}
	}
	m.WordCount = prev
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
