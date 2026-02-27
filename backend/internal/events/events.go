package events

import (
	"time"
)

// import "encoding/json"

type EventType string
type EventDestination string

const (
	EventDestinationAll          EventDestination = "ALL"
	EventDestinationPlayer       EventDestination = "PLAYER"
	EventDestinationOtherPlayers EventDestination = "OTHER_PLAYERS"
)

const (
	EventTypeJoinRoomRequest EventType = "JOIN_ROOM_REQUEST"
	EventTypeJoinRoomOK      EventType = "JOIN_ROOM_OK"
	EventTypeJoinRoomError   EventType = "JOIN_ROOM_ERROR"

	GameRules          EventType = "GAME_RULES"
	RoundOver          EventType = "ROUND_OVER"
	RoundInfo          EventType = "ROUND_INFO"
	RoundWinner        EventType = "ROUND_WINNER"
	PlayerRoundScores  EventType = "PLAYER_ROUND_SCORES"
	NextRoundCountdown EventType = "NEXT_ROUND_COUNTDOWN"

	PlayerWordAccepted EventType = "PLAYER_WORD_ACCEPTED"
	PlayerWordRejected EventType = "PLAYER_WORD_REJECTED"

	PlayerSubmissionBroadcast EventType = "PLAYER_SUBMISSION_BROADCAST"
)

type JoinRoomRequest struct {
	Type    EventType `json:"type"`
	Payload struct {
		Moniker   string `json:"moniker"`
		Timestamp int64  `json:"timestamp"`
	} `json:"payload"`
}

type JoinRoomOK struct {
	Type        EventType        `json:"type"`
	Destination EventDestination `json:"destination"`
	Payload     struct {
		Moniker   string `json:"moniker"`
		PlayerId  string `json:"playerId"`
		Timestamp int64  `json:"timestamp"`
		RoomId    int    `json:"roomId"`
	} `json:"payload"`
}

type JoinRoomError struct {
	Type    EventType `json:"type"`
	Payload struct {
		Message   string `json:"message"`
		Timestamp int64  `json:"timestamp"`
	} `json:"payload"`
}

type EnrichableEvent interface {
	GetType() EventType
	GetDestination() EventDestination
	GetPlayerID() string
	Enrich(moniker string)
}

type GameRulesEvent struct {
	Type    EventType `json:"type"`
	Payload struct {
		Moniker   string `json:"moniker"`
		Message   string `json:"message"`
		Timestamp int64  `json:"timestamp"`
	} `json:"payload"`
}

func (e *GameRulesEvent) GetType() EventType {
	return GameRules
}

func (e *GameRulesEvent) GetDestination() EventDestination {
	return EventDestinationAll
}

func (e *GameRulesEvent) GetPlayerID() string {
	return ""
}

func (e *GameRulesEvent) Enrich(moniker string) {
	e.Payload.Timestamp = time.Now().UnixMilli()

	e.Payload.Moniker = moniker
}

type RoundInfoEvent struct {
	Type    EventType `json:"type"`
	Payload struct {
		Words           []string `json:"words"`
		ValidWordsCount int      `json:"validWordsCount"`
		Timestamp       int64    `json:"timestamp"`
		Moniker         string   `json:"moniker"`
	} `json:"payload"`
}

func (e *RoundInfoEvent) GetType() EventType {
	return RoundInfo
}

func (e *RoundInfoEvent) GetDestination() EventDestination {
	return EventDestinationAll
}

func (e *RoundInfoEvent) GetPlayerID() string {
	return ""
}

func (e *RoundInfoEvent) Enrich(moniker string) {
	e.Payload.Timestamp = time.Now().UnixMilli()
	e.Payload.Moniker = moniker
}

type RoundOverEvent struct {
	Type    EventType `json:"type"`
	Payload struct {
		Message   string `json:"message"`
		Timestamp int64  `json:"timestamp"`
		Moniker   string `json:"moniker"`
	} `json:"payload"`
}

func (e *RoundOverEvent) GetType() EventType {
	return RoundOver
}

func (e *RoundOverEvent) GetDestination() EventDestination {
	return EventDestinationAll
}

func (e *RoundOverEvent) GetPlayerID() string {
	return ""
}

func (e *RoundOverEvent) Enrich(moniker string) {
	e.Payload.Timestamp = time.Now().UnixMilli()
}

type RoundWinnerEvent struct {
	Type    EventType `json:"type"`
	Payload struct {
		PlayerId  string `json:"playerId"`
		Moniker   string `json:"moniker"`
		Score     int    `json:"score"`
		Timestamp int64  `json:"timestamp"`
	} `json:"payload"`
}

func (e *RoundWinnerEvent) GetType() EventType {
	return RoundWinner
}

func (e *RoundWinnerEvent) GetDestination() EventDestination {
	return EventDestinationAll
}

func (e *RoundWinnerEvent) GetPlayerID() string {
	return e.Payload.PlayerId
}

func (e *RoundWinnerEvent) Enrich(moniker string) {
	e.Payload.Moniker = moniker
	e.Payload.Timestamp = time.Now().UnixMilli()
}

type PlayerRoundScoresEvent struct {
	Type    EventType `json:"type"`
	Payload struct {
		Moniker   string `json:"moniker"`
		Score     int    `json:"score"`
		Timestamp int64  `json:"timestamp"`
		PlayerId  string `json:"playerId"`
	} `json:"payload"`
}

func (e *PlayerRoundScoresEvent) GetType() EventType {
	return PlayerRoundScores
}

func (e *PlayerRoundScoresEvent) GetDestination() EventDestination {
	return EventDestinationPlayer
}

func (e *PlayerRoundScoresEvent) GetPlayerID() string {
	return e.Payload.PlayerId
}

func (e *PlayerRoundScoresEvent) Enrich(moniker string) {
	e.Payload.Timestamp = time.Now().UnixMilli()
	e.Payload.Moniker = moniker
}

type PlayerWordAcceptedEvent struct {
	Type    EventType `json:"type"`
	Payload struct {
		PlayerId  string `json:"playerId"`
		Moniker   string `json:"moniker"`
		Word      string `json:"word"`
		Points    int    `json:"points"`
		Timestamp int64  `json:"timestamp"`
	} `json:"payload"`
}

func (e *PlayerWordAcceptedEvent) GetType() EventType {
	return PlayerWordAccepted
}

func (e *PlayerWordAcceptedEvent) GetDestination() EventDestination {
	return EventDestinationPlayer
}

func (e *PlayerWordAcceptedEvent) GetPlayerID() string {
	return e.Payload.PlayerId
}

func (e *PlayerWordAcceptedEvent) Enrich(moniker string) {
	e.Payload.Moniker = moniker
	e.Payload.Timestamp = time.Now().UnixMilli()
}

type PlayerWordRejectedEvent struct {
	Type    EventType `json:"type"`
	Payload struct {
		PlayerId  string `json:"playerId"`
		Moniker   string `json:"moniker"`
		Word      string `json:"word"`
		Message   string `json:"message"`
		Timestamp int64  `json:"timestamp"`
	} `json:"payload"`
}

func (e *PlayerWordRejectedEvent) GetType() EventType {
	return PlayerWordRejected
}
func (e *PlayerWordRejectedEvent) GetDestination() EventDestination {
	return EventDestinationPlayer
}

func (e *PlayerWordRejectedEvent) GetPlayerID() string {
	return e.Payload.PlayerId
}

func (e *PlayerWordRejectedEvent) Enrich(moniker string) {
	e.Payload.Moniker = moniker
	e.Payload.Timestamp = time.Now().UnixMilli()
}

type PlayerSubmissionBroadcastEvent struct {
	Type    EventType `json:"type"`
	Payload struct {
		PlayerId  string `json:"playerId"`
		Moniker   string `json:"moniker"`
		Word      string `json:"word"`
		Timestamp int64  `json:"timestamp"`
	} `json:"payload"`
}

func (e *PlayerSubmissionBroadcastEvent) GetType() EventType {
	return PlayerSubmissionBroadcast
}

func (e *PlayerSubmissionBroadcastEvent) GetDestination() EventDestination {
	return EventDestinationOtherPlayers
}

func (e *PlayerSubmissionBroadcastEvent) GetPlayerID() string {
	return e.Payload.PlayerId
}

func (e *PlayerSubmissionBroadcastEvent) Enrich(moniker string) {
	e.Payload.Moniker = moniker
	e.Payload.Timestamp = time.Now().UnixMilli()
}

type NextRoundCountdownEvent struct {
	Type    EventType `json:"type"`
	Payload struct {
		Timestamp            int64  `json:"timestamp"`
		RoundIntervalSeconds int    `json:"roundIntervalSeconds"`
		Moniker              string `json:"moniker"`
	} `json:"payload"`
}

func (e *NextRoundCountdownEvent) GetType() EventType {
	return NextRoundCountdown
}

func (e *NextRoundCountdownEvent) GetDestination() EventDestination {
	return EventDestinationAll
}

func (e *NextRoundCountdownEvent) GetPlayerID() string {
	return ""
}

func (e *NextRoundCountdownEvent) Enrich(moniker string) {
	e.Payload.Timestamp = time.Now().UnixMilli()
	e.Payload.Moniker = moniker
}
