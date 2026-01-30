package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/muesli/reflow/wordwrap"
	"github.com/yagnikpt/boomtypr/internal/typing"
	"github.com/yagnikpt/boomtypr/internal/utils"
	"github.com/yagnikpt/boomtypr/internal/wordlist"
	"golang.org/x/term"
)

func AssignWords(wordlist *wordlist.WordList, mode typing.Mode, duration time.Duration, wordCount int) (text []rune, lines []Line, lineBreaks []int) {
	timeModeWords := wordCount
	var words []string
	if mode == typing.ModeTime {
		timeModeWords = int(duration.Seconds()) * 3
	}
	if mode == typing.ModeZen {
		words = wordlist.GetAllWords()
	} else {
		words = wordlist.GetRandomWords(timeModeWords)
	}
	joinedWords := strings.Join(words, " ")
	termWidth, _, _ := GetTermDimensions()
	frameX := frameStyles.GetHorizontalFrameSize()
	wrappedPara := wordwrap.String(joinedWords, termWidth-frameX)
	lineBreaks = utils.LineBreakIndexes(wrappedPara)

	return []rune(joinedWords), GetLinesFromWrappedText(wrappedPara), lineBreaks
}

func GetLinesFromWrappedText(text string) []Line {
	lineBreaks := utils.LineBreakIndexes(text)
	lines := make([]Line, utils.CountLines(text))
	linesFromPara := utils.SplitIntoLines(text)
	for i, line := range linesFromPara {
		startIdx := 0
		if i > 0 {
			startIdx = lineBreaks[i-1] + 1
		}

		lines[i] = Line{
			Text:  []rune(line),
			Start: startIdx,
		}
	}
	return lines
}

func GetTermDimensions() (int, int, error) {
	fd := int(os.Stdout.Fd())

	width, height, err := term.GetSize(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting terminal size: %v\n", err)
		return 0, 0, err
	}

	return width, height, nil
}

func CalcHorizontalPadding() int {
	width, _, err := GetTermDimensions()
	if err != nil {
		panic(err)
	}
	if width < 60 {
		return 4
	}
	if width < 80 {
		return 8
	}
	return (width - 80) / 2
}

func CalcPaddingToCenterVertically(height, lines, offset int) string {
	var b strings.Builder
	halfHeight := height / 2
	if height%2 == 0 {
		l := (lines + 1) / 2
		padLength := halfHeight - l - offset
		for range padLength {
			b.WriteString("\n")
		}
	} else {
		l := lines / 2
		padLength := halfHeight - l - offset
		for range padLength {
			b.WriteString("\n")
		}
	}

	return b.String()
}
