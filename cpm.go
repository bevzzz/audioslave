package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type CPM struct {
	Text []rune // Text to try

	CurrentIndex    int   // current tipping index
	MistakesIndexes []int // indexes of mistakes happened

	StartTime time.Time     // starting time
	Duration  time.Duration // tipping duration
}

// FromMonkeyType - Creates a CPM with the given language and text length with texts from MonkeyTypeGame
func FromMonkeyType(language string, textLength int) (*CPM, error) {
	if language == "" {
		language = "english"
	}

	// get text from MonkeyTypeGame
	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/monkeytypegame/monkeytype/master/frontend/static/languages/%s.json", language))
	if err != nil {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		log.Println(err)
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while fetching language (code %d)", resp.StatusCode)
	}

	words := struct {
		Words []string `json:"words"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&words)
	if err != nil {
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
	// create new Rand
	r := rand.New(rand.NewSource(time.Now().Unix()))
	r.Shuffle(len(words.Words), func(i, j int) {
		words.Words[i], words.Words[j] = words.Words[j], words.Words[i]
	})

	return &CPM{
		Text:            []rune(strings.Join(words.Words[:textLength], " ")),
		CurrentIndex:    0,
		MistakesIndexes: make([]int, 0),
		StartTime:       time.Now(),
		Duration:        time.Duration(0),
	}, nil
}

// RegisterRune - Registers a user input
func (c *CPM) RegisterRune(character rune) {
	if c.Text[c.CurrentIndex] != character {
		c.MistakesIndexes = append(c.MistakesIndexes, c.CurrentIndex)
	}
	c.CurrentIndex++
}

// Start - starts the timer and resets settings
func (c *CPM) Start() {
	c.StartTime = time.Now()
	c.CurrentIndex = 0
	c.MistakesIndexes = make([]int, 0)
}

// Stop - stops the timer and calculates the Duration
func (c *CPM) Stop() {
	c.Duration = time.Now().Sub(c.StartTime)
}

// GetAverageCPM - returns the average cpm
func (c *CPM) GetAverageCPM() float64 {
	return float64(c.CurrentIndex-len(c.MistakesIndexes)) / c.Duration.Minutes()
}

// GetText - returns the CPM Text
func (c *CPM) GetText() string {
	return string(c.Text)
}
