package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type CPM struct {
	Text []rune

	CurrentIndex    int
	MistakesIndexes []int

	StartTime time.Time
	Duration  time.Duration
}

func FromMonkeytype(language string, textLength int) (*CPM, error) {
	if language == "" {
		language = "english"
	}

	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/monkeytypegame/monkeytype/master/frontend/static/languages/%s.json", language))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while fetching language (code %d)", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	words := struct {
		Words []string `json:"words"`
	}{}
	if err := json.Unmarshal(bodyBytes, &words); err != nil {
		return nil, err
	}
	// fill words until desired length is full filled
	wordsLength := len(words.Words)
	if wordsLength < textLength {
		mul, frac := math.Modf(float64(textLength-wordsLength) / float64(wordsLength))
		for i := 0; i < int(mul); i++ {
			words.Words = append(words.Words, words.Words[:wordsLength]...)
		}
		remainingLength := int(float64(wordsLength) * frac)
		words.Words = append(words.Words, words.Words[:remainingLength]...)
		wordsLength = len(words.Words)
		if wordsLength < textLength {
			words.Words = append(words.Words, words.Words[:textLength-wordsLength]...)
		}
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	r.Shuffle(len(words.Words), func(i, j int) {
		words.Words[i], words.Words[j] = words.Words[j], words.Words[i]
	})
	fmt.Println(len(words.Words))
	formatted := strings.Join(words.Words[:textLength], " ")
	return &CPM{Text: []rune(formatted)}, nil
}

func (c *CPM) RegisterRune(character rune) {
	if c.Text[c.CurrentIndex] != character {
		c.MistakesIndexes = append(c.MistakesIndexes, c.CurrentIndex)
	}
	c.CurrentIndex++
}

func (c *CPM) Start() {
	c.StartTime = time.Now()
	c.CurrentIndex = 0
	c.MistakesIndexes = make([]int, 0)
}

func (c *CPM) Stop() {
	c.Duration = time.Now().Sub(c.StartTime)
}

func (c *CPM) GetAverageCPM() float64 {
	return float64(c.CurrentIndex-len(c.MistakesIndexes)) / c.Duration.Minutes()
}

func (c *CPM) GetText() string {
	return string(c.Text)
}
