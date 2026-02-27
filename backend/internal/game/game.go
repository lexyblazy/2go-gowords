package game

import (
	"context"
	"math/rand"
	"time"

	"github.com/lexyblazy/gowords/internal/config"
	"github.com/lexyblazy/gowords/internal/dictionary"
	"github.com/lexyblazy/gowords/internal/events"
)

type GameState struct {
	box          *Gamebox
	rounds       chan *GameRound
	rng          *rand.Rand
	currentRound *GameRound
	c            *config.Config
	emitEvent    func(event events.EnrichableEvent)
}

func NewGameState(
	c *config.Config, d *dictionary.Dictionary,
	emitEvent func(event events.EnrichableEvent),
) *GameState {
	rng := rand.New(rand.NewSource(NewSeed()))
	box := NewGamebox(d, rng)

	return &GameState{
		box:          box,
		rng:          rng,
		currentRound: nil,
		rounds:       make(chan *GameRound, 500),
		emitEvent:    emitEvent,
		c:            c,
	}
}

func (gs *GameState) NewRound() *GameRound {
	words := gs.box.GetWordsForRound(gs.c.Game.WordLength, gs.c.Game.WordCount)
	validWords := gs.box.GetValidWordCombinations(words...)
	expansionWord := gs.box.GetExpansionWord(words...)
	validWordsWithExpansion := gs.box.GetValidWordCombinationsWithExpansion(expansionWord, words...)

	return &GameRound{
		words:                   words,
		expansionWord:           expansionWord,
		validWords:              validWords,
		validWordsWithExpansion: validWordsWithExpansion,
		seenWords:               make(map[string]struct{}),
		submissionChan:          make(chan *Submission),
		scores:                  make(map[string]int),
		emitEvent:               gs.emitEvent,
	}
}

func (gs *GameState) RefillRounds() {
	for {

		if len(gs.rounds) >= 500 {

			time.Sleep(10 * time.Second)
			continue
		}

		r := gs.NewRound()

		if gs.box.GetDistinctCharacterCount(r.words...) >= gs.c.Game.DistinctCharacterCount {
			gs.rounds <- r
		}
		// add a delay to avoid too many rounds being created at once
		time.Sleep(100 * time.Millisecond)

	}
}

func (gs *GameState) PlayCurrentRound(timeLimit time.Duration) {

	ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
	defer cancel()

	go gs.PrintRound(ctx)

	// during the last 30 seconds of the round, the expansion word will be added to the words list
	// and the valid words map will be updated to include the expansion word
	go func(c context.Context) {
		time.Sleep(timeLimit - 30*time.Second)
		select {
		case <-c.Done():
			return
		default:
			gs.currentRound.ExpandWords()
		}
	}(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case s := <-gs.currentRound.submissionChan:
			gs.currentRound.handleSubmission(ctx, s)
		}
	}
}

func (gs *GameState) SubmitWord(playerId string, word string) {
	if gs.currentRound == nil {
		return
	}

	gs.currentRound.submissionChan <- &Submission{playerId: playerId, word: word}
}

func (gs *GameState) Run() {

	go gs.RefillRounds()

	// pick one round at random
	timeLimit := 60 * time.Second

	for {

		gs.PrintRules()
		gs.currentRound = <-gs.rounds
		gs.PlayCurrentRound(timeLimit)

		var event events.RoundOverEvent
		event.Payload.Message = "The round is over"
		gs.emitEvent(&event)

		gs.currentRound.ReportScores()
		gs.currentRound = nil

		var nxtRoundCntDwn events.NextRoundCountdownEvent
		nxtRoundCntDwn.Payload.RoundIntervalSeconds = gs.c.Game.RoundIntervalSeconds
		gs.emitEvent(&nxtRoundCntDwn)
		// Sleep for the round interval before starting a new round
		time.Sleep(time.Duration(gs.c.Game.RoundIntervalSeconds) * time.Second)

	}
}

func (gs *GameState) PrintRound(ctx context.Context) {

	if gs.currentRound == nil {
		return
	}

	for {

		select {
		case <-ctx.Done():
			return
		default:
		}

		words := []string{}
		for _, word := range gs.currentRound.words {
			words = append(words, word.Text)
		}

		var event events.RoundInfoEvent
		event.Payload.Words = words
		event.Payload.ValidWordsCount = len(gs.currentRound.validWords)
		gs.emitEvent(&event)
		time.Sleep(5 * time.Second)

	}
}

func (gs *GameState) PrintRules() {
	rules := `
Rules:
You are given a set of words.
You can use the words to form other words.
You can only use each word once.
The minimum word length is 3 letters.
During the last 30 seconds of the round, a two letter word will be added to expand the possible words.
Happy Guessing!`

	var event events.GameRulesEvent
	event.Payload.Message = rules

	gs.emitEvent(&event)

	time.Sleep(5 * time.Second)
}
