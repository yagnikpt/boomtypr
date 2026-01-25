package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/yagnikpt/boomtypr/internal/utils"
	"golang.org/x/term"
)

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

func CalcPaddingToCenterVertically(height, lines, padding int) string {
	var b strings.Builder
	halfHeight := height / 2
	halfPadding := padding / 2
	if height%2 == 0 {
		l := (lines + 1) / 2
		padLength := halfHeight - l - halfPadding
		for range padLength {
			b.WriteString("\n")
		}
	} else {
		l := lines / 2
		padLength := halfHeight - l - halfPadding
		for range padLength {
			b.WriteString("\n")
		}
	}

	return b.String()
}
