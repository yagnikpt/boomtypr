package wordlist

import (
	_ "embed"
	"encoding/json"
	"math/rand"
)

//go:embed repo/en.json
var repo []byte

type WordList struct {
	Name  string   `json:"name"`
	Words []string `json:"words"`
}

func New() *WordList {
	var wl WordList

	if err := json.Unmarshal(repo, &wl); err != nil {
		panic(err)
	}

	return &wl
}

func (wl WordList) GetRandomWords(n int) []string {
	if n <= 0 || len(wl.Words) == 0 {
		return nil
	}

	result := make([]string, n)
	for i := range n {
		result[i] = wl.Words[rand.Intn(len(wl.Words))]
	}
	return result
}

func (wl WordList) GetAllWords() []string {
	words := make([]string, len(wl.Words))
	copy(words, wl.Words)
	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})
	return words
}
