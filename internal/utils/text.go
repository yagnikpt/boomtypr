package utils

import (
	"bufio"
	"strings"
)

func SplitIntoLines(para string) []string {
	scanner := bufio.NewScanner(strings.NewReader(para))
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, strings.TrimSpace(line))
	}
	return lines
}

func LineBreakIndexes(text string) []int {
	var indexes []int
	for i, char := range text {
		if char == '\n' {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func CountLines(text string) int {
	scanner := bufio.NewScanner(strings.NewReader(text))
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}
