package ui

import (
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/muesli/reflow/wordwrap"
	"github.com/yagnikpt/boomtypr/internal/typing"
	"github.com/yagnikpt/boomtypr/internal/utils"
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

func NewModel(wordlist *wordlist.WordList) Model {
	words := wordlist.GetRandomWords(50)
	joinedWords := strings.Join(words, " ")
	termWidth, _, _ := GetTermDimensions()
	frameX := frameStyles.GetHorizontalFrameSize()
	wrappedPara := wordwrap.String(joinedWords, termWidth-frameX)
	lineBreaks := utils.LineBreakIndexes(wrappedPara)

	return Model{
		Text:   []rune(joinedWords),
		State:  StateMenu,
		Lines:  GetLinesFromWrappedText(wrappedPara),
		Engine: typing.NewEngine(joinedWords, lineBreaks),
		Model: typing.Model{
			Mode:              typing.ModeTime,
			Duration:          time.Second * 30,
			RemainingDuration: 30,
			Target:            []rune(joinedWords),
			WordsCount:        50,
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
		case "enter":
			if m.State == StateResults {
				wl := wordlist.New()
				prevMode := m.Mode
				prevDuration := m.Duration
				prevWordsCount := m.WordsCount
				prevWidth := m.Width
				prevHeight := m.Height
				m = NewModel(wl)
				m.Mode = prevMode
				m.Duration = prevDuration
				m.RemainingDuration = int(prevDuration.Seconds())
				m.WordsCount = prevWordsCount
				m.Width = prevWidth
				m.Height = prevHeight
			}
		case "backspace":
			if m.State == StateTyping {
				m.Engine.Backspace()
				m.AddKeystroke(32, m.Target[m.Engine.CurrentChar], true)
			}
		default:
			if len(msg.Runes) > 0 && !m.Done {
				if m.State == StateMenu {
					cmds = append(cmds, startCmd())
				}
				m.AddKeystroke(msg.Runes[0], m.Target[m.Engine.CurrentChar], false)
				m.Engine.TypeChar(msg.Runes[0])
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

	padding := CalcPaddingToCenterVertically(m.Height, LinesWindowSize, centerOffset)
	b.WriteString(padding)

	switch m.State {
	case StateMenu:
		switch m.Mode {
		case typing.ModeTime:
			b.WriteString(helperStats.Render(strconv.Itoa(int(m.Duration.Seconds()))+" secs") + "\n\n")
		case typing.ModeWords:
			b.WriteString(helperStats.Render(strconv.Itoa(m.WordsCount)+" words") + "\n\n")
		}
	case StateTyping:
		switch m.Mode {
		case typing.ModeTime:
			b.WriteString(helperStats.Render(strconv.Itoa(m.RemainingDuration)) + "\n\n")
		case typing.ModeWords:
			b.WriteString(helperStats.Render(strconv.Itoa(m.CurrentWord)+" / "+strconv.Itoa(m.WordsCount)) + "\n\n")
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
