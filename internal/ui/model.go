package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/muesli/reflow/wordwrap"
	"github.com/yagnikpt/boomtypr/internal/typing"
	"github.com/yagnikpt/boomtypr/internal/utils"
)

var (
	frameStyles        = lipgloss.NewStyle().Padding(2, CalcHorizontalPadding())
	pendingCharStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#4C4C4C"))
	incorrectCharStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Underline(true)
	cursorStyle        = lipgloss.NewStyle().Background(lipgloss.Color("7")).Foreground(lipgloss.Color("0"))
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
	Text          string
	State         UIState
	Engine        *typing.Engine
	Lines         []Line
	Width         int
	Height        int
	CursorVisible bool
	Blinking      bool
	BlinkID       int
}

func NewModel(words []string) Model {
	joinedWords := strings.Join(words, " ")
	termWidth, _, _ := GetTermDimensions()
	frameX, _ := frameStyles.GetFrameSize()
	wrappedPara := wordwrap.String(joinedWords, termWidth-frameX)
	lineBreaks := utils.LineBreakIndexes(wrappedPara)

	return Model{
		Text:     joinedWords,
		State:    StateMenu,
		Lines:    GetLinesFromWrappedText(wrappedPara),
		Engine:   typing.NewEngine(joinedWords, lineBreaks),
		Blinking: true,
	}
}

func (m Model) Init() tea.Cmd {
	return blinkCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case blinkMsg:
		if m.Blinking {
			m.CursorVisible = !m.CursorVisible
			return m, blinkCmd()
		}
		return m, nil

	case resumeBlinkMsg:
		if msg.id != m.BlinkID {
			return m, nil
		}
		m.Blinking = true
		return m, blinkCmd()

	case tea.WindowSizeMsg:
		frameStyles = frameStyles.Padding(2, CalcHorizontalPadding())
		frameX, _ := frameStyles.GetFrameSize()
		wrappedPara := wordwrap.String(m.Text, msg.Width-frameX)
		m.Lines = GetLinesFromWrappedText(wrappedPara)
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		m.Blinking = false
		m.CursorVisible = true
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "backspace":
			m.Engine.Backspace()
		default:
			if len(msg.Runes) > 0 {
				m.Engine.TypeChar(msg.Runes[0])
			}
		}

		m.BlinkID++
		return m, resumeBlinkCmd(m.BlinkID)
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	padding := CalcPaddingToCenterVertically(m.Height, LinesWindowSize, 2)
	b.WriteString(padding)
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
			}

			if charIndex == m.Engine.CurrentChar && m.CursorVisible {
				rendered = cursorStyle.Render(rendered)
			}

			b.WriteString(rendered)
		}
		if m.Engine.CurrentChar == line.Start+len(line.Text) && m.CursorVisible {
			b.WriteString(cursorStyle.Render(" "))
		}
		b.WriteString("\n")
	}

	return frameStyles.Render(b.String())
}
