package game

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/lexyblazy/gowords/internal/config"
	"github.com/lexyblazy/gowords/internal/dictionary"
	"github.com/lexyblazy/gowords/internal/events"
	"github.com/lexyblazy/gowords/internal/store"
)

type GameState struct {
	box          *Gamebox
	rounds       chan *GameRound
	rng          *rand.Rand
	currentRound *GameRound
	c            *config.Config
	emitEvent    func(event events.EnrichableEvent)

	rs *store.RedisStore
}

func NewGameState(
	c *config.Config, d *dictionary.Dictionary,
	emitEvent func(event events.EnrichableEvent),
	rs *store.RedisStore,
) *GameState {
	rng := rand.New(rand.NewSource(NewSeed()))
	box := NewGamebox(d, rng)

	return &GameState{
		box:          box,
		rng:          rng,
		currentRound: nil,
		rounds:       make(chan *GameRound, c.Game.RoundCount),
		emitEvent:    emitEvent,
		c:            c,
		rs:           rs,
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
		submissionChan:          make(chan *events.PlayerWordSubmissionEvent),
		scores:                  make(map[string]int),
		emitEvent:               gs.emitEvent,
		updateLeaderBoards:      gs.updateLeaderBoards,
	}
}

func (gs *GameState) RefillRounds() {
	for {

		r := gs.NewRound()

		if gs.box.GetDistinctCharacterCount(r.words...) >= gs.c.Game.DistinctCharacterCount {
			gs.rounds <- r
		}

	}
}

func (gs *GameState) PlayCurrentRound() {

	roundDuration := time.Duration(gs.c.Game.RoundDurationSeconds) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), roundDuration)
	defer cancel()

	go gs.BroadcastRoundInfoPeriodically(ctx)

	// during the last 30 seconds of the round, the expansion word will be added to the words list
	// and the valid words map will be updated to include the expansion word
	go func(c context.Context) {
		time.Sleep(roundDuration - 30*time.Second)
		select {
		case <-c.Done():
			return
		default:
			gs.currentRound.ExpandWords()
			gs.BroadcastRoundInfo()
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

func (gs *GameState) SubmitWord(event *events.PlayerWordSubmissionEvent) {
	if gs.currentRound == nil {
		return
	}

	gs.currentRound.submissionChan <- event
}

func (gs *GameState) Run() {

	go gs.RefillRounds()

	// pick one round at random

	for {

		gs.currentRound = <-gs.rounds
		gs.currentRound.endsAt = time.Now().Add(time.Duration(gs.c.Game.RoundDurationSeconds) * time.Second).UnixMilli()

		// broadcast the round info immediately
		gs.BroadcastRoundInfo()
		gs.PlayCurrentRound()

		var event events.RoundOverEvent
		event.Type = events.RoundOver
		event.Payload.Message = "The round is over"
		gs.emitEvent(&event)

		gs.currentRound.ReportScores()
		gs.currentRound = nil

		var nxtRoundCntDwn events.NextRoundCountdownEvent
		nxtRoundCntDwn.Type = events.NextRoundCountdown
		nxtRoundCntDwn.Payload.EndsAt = time.Now().Add(time.Duration(gs.c.Game.RoundIntervalSeconds) * time.Second).UnixMilli()
		gs.emitEvent(&nxtRoundCntDwn)
		gs.BroadcastRules()
		// Sleep for the round interval before starting a new round
		time.Sleep(time.Duration(gs.c.Game.RoundIntervalSeconds) * time.Second)

	}
}

func (gs *GameState) BroadcastRoundInfo() {

	if gs.currentRound == nil {
		return
	}

	words := []string{}
	for _, word := range gs.currentRound.words {
		words = append(words, word.Text)
	}

	var event events.RoundInfoEvent
	event.Type = events.RoundInfo
	event.Payload.Words = words
	event.Payload.ValidWordsCount = len(gs.currentRound.validWords)
	event.Payload.EndsAt = gs.currentRound.endsAt
	gs.emitEvent(&event)

}

func (gs *GameState) BroadcastRoundInfoPeriodically(ctx context.Context) {

	for {

		select {
		case <-ctx.Done():
			return
		default:
			gs.BroadcastRoundInfo()
		}

		time.Sleep(time.Duration(gs.c.Game.PrintRoundIntervalSeconds) * time.Second)

	}
}

func (gs *GameState) BroadcastRules() {
	var event events.GameRulesEvent
	event.Type = events.GameRules
	event.Payload.Rules = strings.Split(gs.c.Game.Rules, "\n")

	gs.emitEvent(&event)

}

func (gs *GameState) updateLeaderBoards(scoresMap map[string]int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	gs.rs.UpdateLeaderBoards(ctx, scoresMap)
}
