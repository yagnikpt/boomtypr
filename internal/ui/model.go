package ui

import (
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/muesli/reflow/wordwrap"
	"github.com/yagnikpt/boomtypr/internal/typing"
	"github.com/yagnikpt/boomtypr/internal/wordlist"
)

var (
	frameStyles        = lipgloss.NewStyle().Padding(2, CalcHorizontalPadding())
	pendingCharStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // grey
	correctCharStyle   = lipgloss.NewStyle().Bold(true)
	incorrectCharStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Underline(true) // red
	cursorStyle        = lipgloss.NewStyle().Background(lipgloss.Color("7")).Foreground(lipgloss.Color("0"))
	helperStats        = lipgloss.NewStyle().Foreground(lipgloss.Color("4")) // blue
)

type UIState int

const (
	StateMenu UIState = iota
	StateTyping
	StateResults
)

var LinesWindowSize = 3

type Line struct {
	Text  []rune
	Start int
}

type Model struct {
	Text   []rune
	State  UIState
	Engine *typing.Engine
	Lines  []Line
	Width  int
	Height int
	Stats  *typing.Stats
	typing.Model
}

func NewModel(wordlist *wordlist.WordList, mode typing.Mode, duration time.Duration, wordCount int) Model {
	text, lines, lineBreaks := AssignWords(wordlist, mode, duration, wordCount)

	return Model{
		Text:   text,
		State:  StateMenu,
		Lines:  lines,
		Engine: typing.NewEngine(text, lineBreaks),
		Model: typing.Model{
			Mode:              mode,
			Duration:          duration,
			RemainingDuration: int(duration.Seconds()),
			Target:            text,
			WordCount:         wordCount,
			CurrentWord:       1,
		},
		Stats: typing.NewStats(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case startMsg:
		m.Stats.Start()
		m.State = StateTyping
		if m.Mode == typing.ModeTime {
			cmds = append(cmds, endTickerCmd(m.Duration), tickerCmd())
		}

	case tickerMsg:
		if m.State == StateTyping && m.Mode == typing.ModeTime {
			m.RemainingDuration--
			if !m.Done {
				cmds = append(cmds, tickerCmd())
			}
		}

	case endTickerMsg:
		if m.State == StateTyping {
			cmds = append(cmds, finishCmd())
		}

	case finishMsg:
		m.Done = true
		m.Stats.Stop()
		m.Stats.Calculate(m.Keystrokes)
		m.State = StateResults

	case tea.WindowSizeMsg:
		frameStyles = frameStyles.Padding(2, CalcHorizontalPadding())
		frameX, _ := frameStyles.GetFrameSize()
		wrappedPara := wordwrap.String(string(m.Text), msg.Width-frameX)
		m.Lines = GetLinesFromWrappedText(wrappedPara)
		newLineBreaks := make([]int, len(m.Lines)-1)
		for i, lines := range m.Lines {
			if i != len(m.Lines)-1 {
				newLineBreaks[i] = lines.Start + len(lines.Text)
			}
		}
		m.Engine.UpdateLines(newLineBreaks)
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab":
			if m.State == StateMenu {
				m.NextMode(m.Mode)
				wl := wordlist.New()
				text, lines, lineBreaks := AssignWords(wl, m.Mode, m.Duration, m.WordCount)
				m.Text = text
				m.Target = text
				m.Lines = lines
				m.Engine.UpdateLines(lineBreaks)
				m.Engine.Text = text
				m.Engine.Track = make([]typing.CharState, len(text))
			}
		case "up":
			if m.State == StateMenu {
				switch m.Mode {
				case typing.ModeTime:
					m.NextDuration()
				case typing.ModeWords:
					m.NextWordCount()
				}
				wl := wordlist.New()
				text, lines, lineBreaks := AssignWords(wl, m.Mode, m.Duration, m.WordCount)
				m.Text = text
				m.Target = text
				m.Lines = lines
				m.Engine.UpdateLines(lineBreaks)
				m.Engine.Text = text
				m.Engine.Track = make([]typing.CharState, len(text))
			}
		case "down":
			if m.State == StateMenu {
				switch m.Mode {
				case typing.ModeTime:
					m.PrevDuration()
				case typing.ModeWords:
					m.PrevWordCount()
				}
				wl := wordlist.New()
				text, lines, lineBreaks := AssignWords(wl, m.Mode, m.Duration, m.WordCount)
				m.Text = text
				m.Target = text
				m.Lines = lines
				m.Engine.UpdateLines(lineBreaks)
				m.Engine.Text = text
				m.Engine.Track = make([]typing.CharState, len(text))
			}
		case "enter":
			if m.State == StateResults {
				wl := wordlist.New()
				prevWidth := m.Width
				prevHeight := m.Height
				m = NewModel(wl, m.Mode, m.Duration, m.WordCount)
				m.Width = prevWidth
				m.Height = prevHeight
			}
			if m.State == StateTyping && m.Mode == typing.ModeZen {
				cmds = append(cmds, finishCmd())
			}
		case "backspace":
			if m.State == StateTyping && m.Engine.CurrentChar > 0 {
				prevChar := m.Engine.Text[m.Engine.CurrentChar-1]

				if prevChar == ' ' && m.CurrentWord > 1 {
					m.CurrentWord--
				}

				m.Engine.Backspace()

				if m.Engine.CurrentChar < len(m.Engine.Text) {
					m.AddKeystroke(' ', m.Engine.Text[m.Engine.CurrentChar], true)
				}
			}
		case " ":
			if m.State == StateTyping && !m.Done && !m.Engine.Finished {
				if m.Engine.CurrentChar > 0 {
					prevChar := m.Engine.Text[m.Engine.CurrentChar-1]

					if prevChar != ' ' && m.Engine.Text[m.Engine.CurrentChar] == 32 {
						m.CurrentWord++
					}
				}

				m.AddKeystroke(' ', m.Engine.Text[m.Engine.CurrentChar], false)
				m.Engine.TypeChar(' ')
			}
		default:
			if len(msg.Runes) > 0 && !m.Done {
				if m.State == StateMenu {
					cmds = append(cmds, startCmd())
				}
				m.AddKeystroke(msg.Runes[0], m.Engine.Text[m.Engine.CurrentChar], false)
				m.Engine.TypeChar(msg.Runes[0])
				if m.Engine.Finished && m.Mode != typing.ModeTime {
					cmds = append(cmds, finishCmd())
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var b strings.Builder

	centerOffset := 4 // padding
	if m.State == StateResults {
		centerOffset -= 2
	}
	if m.State == StateTyping && m.Mode == typing.ModeZen {
		centerOffset -= 2
	}

	padding := CalcPaddingToCenterVertically(m.Height, LinesWindowSize, centerOffset)
	b.WriteString(padding)

	switch m.State {
	case StateMenu:
		switch m.Mode {
		case typing.ModeTime:
			b.WriteString(helperStats.Render(strconv.Itoa(int(m.Duration.Seconds()))+" secs") + "\n\n")
		case typing.ModeWords:
			b.WriteString(helperStats.Render(strconv.Itoa(m.WordCount)+" words") + "\n\n")
		case typing.ModeZen:
			b.WriteString(helperStats.Render("zen") + "\n\n")
		}
	case StateTyping:
		switch m.Mode {
		case typing.ModeTime:
			b.WriteString(helperStats.Render(strconv.Itoa(m.RemainingDuration)) + "\n\n")
		case typing.ModeWords:
			b.WriteString(helperStats.Render(strconv.Itoa(m.CurrentWord)+" / "+strconv.Itoa(m.WordCount)) + "\n\n")
		}
	}

	if m.State == StateMenu || m.State == StateTyping {
		for i := m.Engine.CurrentLine; i < m.Engine.CurrentLine+LinesWindowSize && i < len(m.Lines); i++ {
			line := m.Lines[i]
			for j, char := range line.Text {
				charIndex := line.Start + j

				rendered := string(char)
				switch m.Engine.Track[charIndex] {
				case typing.CharPending:
					rendered = pendingCharStyle.Render(rendered)
				case typing.CharIncorrect:
					rendered = incorrectCharStyle.Render(rendered)
				case typing.CharCorrect:
					rendered = correctCharStyle.Render(rendered)
				}

				if charIndex == m.Engine.CurrentChar {
					rendered = cursorStyle.Render(string(char))
				}
				b.WriteString(rendered)
			}
			if m.Engine.CurrentChar == line.Start+len(line.Text) {
				b.WriteString(cursorStyle.Render(" "))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	centerStyles := lipgloss.NewStyle().Width(m.Width - frameStyles.GetHorizontalFrameSize()).AlignHorizontal(lipgloss.Center)

	if m.Done && m.State == StateResults {
		b.WriteString(centerStyles.Foreground(lipgloss.Color("4")).Render("WPM: "+strconv.Itoa(int(m.Stats.WPM()))+", Accuracy: "+strconv.Itoa(int(m.Stats.Accuracy()))+"%") + "\n\n")
		b.WriteString(centerStyles.Foreground(lipgloss.Color("8")).Render("Press Enter to restart • Esc to quit"))
	}

	return frameStyles.Render(b.String())
}
