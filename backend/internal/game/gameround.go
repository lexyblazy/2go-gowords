package game

import (
	"context"
	"fmt"
	"strings"

	"github.com/lexyblazy/gowords/internal/dictionary"
)

type Submission struct {
	playerId string
	word     string
}

type GameRound struct {
	words                   []dictionary.Word
	expansionWord           dictionary.Word
	validWords              map[string]struct{}
	validWordsWithExpansion map[string]struct{}
	seenWords               map[string]struct{}
	scores                  map[string]int

	submissionChan chan *Submission
	emitCb         func(eventType EventType, payload any)
}

func (gr *GameRound) handleSubmission(ctx context.Context, s *Submission) {
	// input is only one word, multiple words are not allowed

	select {
	case <-ctx.Done():
		return
	default:
	}

	word := s.word
	playerId := s.playerId

	length := strings.Split(word, " ")
	if len(length) > 1 {
		gr.emitEvent(EventTypePlayerWordRejected, BasicPayload{Message: "Multiple words are not allowed", PlayerId: playerId})
		return
	}
	word = strings.ToLower(strings.TrimSpace(word))

	if len(word) < 3 {
		gr.emitEvent(EventTypePlayerWordRejected, BasicPayload{Message: "Words must be at least 3 letters long", PlayerId: playerId})
		return
	}

	// check if the word is in the valid words
	if _, ok := gr.validWords[word]; !ok {
		gr.emitEvent(EventTypePlayerWordRejected, BasicPayload{Message: fmt.Sprintf("%s is not a valid submission", word), PlayerId: playerId})
		return
	}

	// check if the word has already been seen
	if _, ok := gr.seenWords[word]; ok {
		gr.emitEvent(EventTypePlayerWordRejected, BasicPayload{Message: fmt.Sprintf("%s has already been used", word), PlayerId: playerId})
		return
	}

	gr.seenWords[word] = struct{}{}

	gr.AwardPoints(word, playerId)

}

func (gr *GameRound) emitEvent(eventType EventType, payload any) {

	gr.emitCb(eventType, payload)
}

func (gr *GameRound) AwardPoints(word string, playerId string) {
	points := len(word) - 2
	gr.scores[playerId] += points
	content := fmt.Sprintf("You earned %d points for %s", points, word)

	gr.emitEvent(EventTypePlayerWordAccepted, BasicPayload{Message: content, PlayerId: playerId})

	// send the players submission to the general channel excluding the player
	gr.emitEvent(EventTypeGeneralExcludingPlayer, BasicPayload{Message: word, PlayerId: playerId})

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
			content := fmt.Sprintf("🎉 You scored a total of %d points in this round", score)
			gr.emitEvent(EventTypePlayerRoundScores, BasicPayload{Message: content, PlayerId: playerId})
		}
	}

	if winningScore > 0 {
		gr.emitEvent(EventTypeRoundWinner, RoundWinnerPayload{PlayerId: winningPlayerId, Score: winningScore})
	}

}

func (gr *GameRound) ExpandWords() {
	gr.words = append(gr.words, gr.expansionWord)
	gr.validWords = gr.validWordsWithExpansion
}
