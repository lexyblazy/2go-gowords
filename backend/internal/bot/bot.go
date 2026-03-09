package bot

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/lexyblazy/gowords/internal/dictionary"
	"github.com/lexyblazy/gowords/internal/events"
)

type IncomingMessage struct {
	Type events.EventType `json:"type"`
}

type CurrentRoundInfo struct {
	WordsList []string            `json:"words"`
	EndsAt    int64               `json:"endsAt"`
	Solutions map[string]struct{} `json:"solutions"`
	UsedWords map[string]struct{} `json:"usedWords"`
}

type Bot struct {
	d            *dictionary.Dictionary
	cr           *CurrentRoundInfo
	roundOver    chan bool
	submitWordCb func(text string)
	mu           *sync.Mutex
}

func NewBot(d *dictionary.Dictionary, submitWordCb func(text string)) *Bot {

	return &Bot{
		d:            d,
		cr:           nil,
		roundOver:    make(chan bool),
		submitWordCb: submitWordCb,
		mu:           &sync.Mutex{},
	}
}

func (b *Bot) getMessageType(message []byte) events.EventType {
	var incomingMessage IncomingMessage
	err := json.Unmarshal(message, &incomingMessage)
	if err != nil {
		return ""
	}
	return incomingMessage.Type
}

func makeWords(words []string) []dictionary.Word {
	var w []dictionary.Word
	for _, word := range words {
		w = append(w, dictionary.NewWord(word))
	}
	return w
}

func (b *Bot) handleRoundInfo(event events.RoundInfoEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.cr == nil {
		b.cr = &CurrentRoundInfo{
			WordsList: event.Payload.Words,
			EndsAt:    event.Payload.EndsAt,
			Solutions: make(map[string]struct{}),
			UsedWords: make(map[string]struct{}),
		}

		words := makeWords(event.Payload.Words)

		// get solutions for the words
		solutions := b.d.GenerateValidWords(words...)
		b.cr.Solutions = solutions
	} else {
		// we do an update only if the words are different and we know this because the length of the words is different
		if len(b.cr.WordsList) != len(event.Payload.Words) {
			words := makeWords(event.Payload.Words)
			solutions := b.d.GenerateValidWords(words...)
			b.cr.Solutions = solutions
		}
	}
}

func (b *Bot) HandleMessage(message []byte) error {
	messageType := b.getMessageType(message)
	switch messageType {
	case events.RoundInfo:
		var roundInfoEvent events.RoundInfoEvent
		err := json.Unmarshal(message, &roundInfoEvent)
		if err != nil {
			return err
		}
		b.handleRoundInfo(roundInfoEvent)

	case events.RoundOver:
		var roundOverEvent events.RoundOverEvent
		err := json.Unmarshal(message, &roundOverEvent)
		if err != nil {
			return err
		}
		b.mu.Lock()
		b.cr = nil
		b.mu.Unlock()
		b.roundOver <- true
		time.Sleep(20 * time.Second)
		log.Println("Round over")
	default:
		log.Println(messageType, string(message))
	}

	return nil
}

func (b *Bot) submitSolutions() {

	for word := range b.cr.Solutions {
		select {
		case <-b.roundOver:
			return
		default:
		}

		if b.cr == nil {
			return
		}

		if _, ok := b.cr.UsedWords[word]; ok {
			continue
		}
		b.cr.UsedWords[word] = struct{}{}
		b.submitWordCb(word)
		time.Sleep(time.Duration(500+rand.Intn(500)) * time.Millisecond)

	}
}

func (b *Bot) PlayRound() {

	for {

		if b.cr == nil {
			time.Sleep(10 * time.Second)
			continue
		}

		b.submitSolutions()

	}

}
