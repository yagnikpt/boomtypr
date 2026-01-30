package typing

import (
	"sort"
)

type CharState int

const (
	CharPending CharState = iota
	CharCorrect
	CharIncorrect
)

type Engine struct {
	Text        []rune
	CurrentChar int
	CurrentLine int
	LineBreaks  []int
	Track       []CharState
	Finished    bool
}

func NewEngine(text []rune, lineBreaks []int) *Engine {
	track := make([]CharState, len(text))
	lbs := make([]int, len(lineBreaks))
	copy(lbs, lineBreaks)
	return &Engine{
		Text:       text,
		LineBreaks: lbs,
		Track:      track,
	}
}

func (e *Engine) TypeChar(char rune) {
	if e.Finished {
		return
	}
	if e.CurrentLine < len(e.LineBreaks) {
		if e.LineBreaks[e.CurrentLine] == e.CurrentChar && string(char) != " " {
			return
		}
		if e.CurrentChar < len(e.Text) && e.CurrentChar == e.LineBreaks[e.CurrentLine] {
			e.CurrentLine++
		}
	}
	if string(char) == string(e.Text[e.CurrentChar]) {
		e.Track[e.CurrentChar] = CharCorrect
	} else {
		e.Track[e.CurrentChar] = CharIncorrect
	}
	e.CurrentChar++
	if e.CurrentChar >= len(e.Text) {
		e.Finished = true
	}
}

func (e *Engine) Backspace() {
	if e.CurrentChar == 0 {
		return
	}
	e.Finished = false
	// if e.CurrentChar < len(e.Text) {
	if e.CurrentLine > 0 && e.CurrentChar == e.LineBreaks[e.CurrentLine-1]+1 {
		e.CurrentLine--
	}
	// }
	e.CurrentChar--
	e.Track[e.CurrentChar] = CharPending
}

func (e *Engine) UpdateLines(newLines []int) {
	e.LineBreaks = newLines

	e.CurrentLine = sort.Search(len(newLines), func(i int) bool {
		return newLines[i] > e.CurrentChar
	})
}

// func (e *Engine) Reset() {
// 	e.CurrentChar = 0
// 	e.CurrentLine = 0
// 	e.Finished = false

// 	for i := range e.Track {
// 		e.Track[i] = CharPending
// 	}
// }
