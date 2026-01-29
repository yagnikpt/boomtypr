package typing

import "time"

type Stats struct {
	StartTime    time.Time
	EndTime      time.Time
	TotalChars   int
	CorrectChars int
	Errors       int
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) Calculate(keys []Keystroke) {
	correct := 0
	incorrect := 0

	for _, k := range keys {
		if k.Backspace {
			continue
		}

		if k.Correct {
			correct++
		} else {
			incorrect++
		}
	}

	s.TotalChars = len(keys)
	s.CorrectChars = correct
	s.Errors = incorrect
}

func (s *Stats) Start() {
	s.StartTime = time.Now()
}

func (s *Stats) Stop() {
	s.EndTime = time.Now()
}

func (s Stats) WPM() float64 {
	duration := s.EndTime.Sub(s.StartTime).Minutes()
	return float64(s.CorrectChars) / 5.0 / duration
}

func (s Stats) Accuracy() float64 {
	return float64(s.CorrectChars) / float64(s.TotalChars) * 100
}

// func (s *Stats) Reset() {
// 	*s = Stats{}
// }
