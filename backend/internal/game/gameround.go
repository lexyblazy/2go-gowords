package game

import (
	"context"
	"strings"

	"github.com/lexyblazy/gowords/internal/dictionary"
	"github.com/lexyblazy/gowords/internal/events"
)

type GameRound struct {
	words                   []dictionary.Word
	expansionWord           dictionary.Word
	validWords              map[string]struct{}
	validWordsWithExpansion map[string]struct{}
	seenWords               map[string]struct{}
	scores                  map[string]int

	submissionChan chan *events.PlayerWordSubmissionEvent
	emitEvent      func(event events.EnrichableEvent)
}

func (gr *GameRound) makeWordRejectedEvent(message string, playerId string, word string) {
	var event events.PlayerWordRejectedEvent
	event.Type = events.PlayerWordRejected
	event.Payload.Message = message
	event.Payload.PlayerId = playerId
	event.Payload.Word = word
	gr.emitEvent(&event)
}

func (gr *GameRound) handleSubmission(ctx context.Context, event *events.PlayerWordSubmissionEvent) {
	// input is only one word, multiple words are not allowed

	select {
	case <-ctx.Done():
		return
	default:
	}

	word := event.Payload.Word
	playerId := event.Payload.PlayerId

	length := strings.Split(word, " ")
	if len(length) > 1 {
		gr.makeWordRejectedEvent("multiple words are not allowed", playerId, word)
		return

	}
	word = strings.ToLower(strings.TrimSpace(word))

	if len(word) < 3 {
		gr.makeWordRejectedEvent("word must be at least 3 letters long", playerId, word)
		return
	}

	// check if the word is in the valid words
	if _, ok := gr.validWords[word]; !ok {
		gr.makeWordRejectedEvent("not a valid submission", playerId, word)
		return
	}

	// check if the word has already been seen
	if _, ok := gr.seenWords[word]; ok {
		gr.makeWordRejectedEvent("it's already been used", playerId, word)
		return
	}

	gr.seenWords[word] = struct{}{}

	gr.AwardPoints(word, playerId)

}

func (gr *GameRound) AwardPoints(word string, playerId string) {
	points := len(word) - 2
	gr.scores[playerId] += points

	var event events.PlayerWordAcceptedEvent
	event.Type = events.PlayerWordAccepted
	event.Payload.Word = word
	event.Payload.Points = points
	event.Payload.PlayerId = playerId
	gr.emitEvent(&event)

	// send the players submission to the general channel excluding the player
	var playerSubmissionEvent events.PlayerSubmissionBroadcastEvent
	playerSubmissionEvent.Type = events.PlayerSubmissionBroadcast
	playerSubmissionEvent.Payload.Word = word
	playerSubmissionEvent.Payload.PlayerId = playerId
	gr.emitEvent(&playerSubmissionEvent)

}

func (gr *GameRound) ReportScores() {
	winningPlayerId := ""
	winningScore := 0

	for playerId, score := range gr.scores {

		if score > winningScore {
			winningScore = score
			winningPlayerId = playerId
		}

		// send each player their score
		if score > 0 {
			var event events.PlayerRoundScoresEvent
			event.Type = events.PlayerRoundScores
			event.Payload.Score = score
			event.Payload.PlayerId = playerId
			gr.emitEvent(&event)
		}
	}

	if winningScore > 0 {
		var event events.RoundWinnerEvent
		event.Type = events.RoundWinner
		event.Payload.WinningPlayerId = winningPlayerId
		event.Payload.Score = winningScore
		gr.emitEvent(&event)
	}

}

func (gr *GameRound) ExpandWords() {
	gr.words = append(gr.words, gr.expansionWord)
	gr.validWords = gr.validWordsWithExpansion
}
