package game

import (
	"encoding/json"
	"fmt"
)

type EventType string

const (
	EventTypeGeneral EventType = "general"
	// EventTypeForPlayer              EventType = "forPlayer"
	EventTypeGeneralExcludingPlayer EventType = "generalExcludingPlayer"

	EventTypeJoinRoom EventType = "joinRoom"

	EventTypePlayerWordAccepted EventType = "playerWordAccepted"
	EventTypePlayerWordRejected EventType = "playerWordRejected"
	EventTypePlayerRoundScores  EventType = "playerRoundScores"
	EventTypeRoundWinner        EventType = "roundWinner"
)

type Event struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (e *Event) String() string {
	return fmt.Sprintf("Event{Type: %s, Payload: %s}", e.Type, string(e.Payload))
}

func (e *Event) ToBytes() []byte {
	eventBytes, err := json.Marshal(e)

	if err != nil {
		fmt.Println("Error marshalling event:", err)
		return nil
	}
	return eventBytes
}



type BasicPayload struct {
	Message  string `json:"message"`
	PlayerId string `json:"playerId"`
}

type RoundWinnerPayload struct {
	PlayerId string `json:"playerId"`
	Score    int    `json:"score"`
}

