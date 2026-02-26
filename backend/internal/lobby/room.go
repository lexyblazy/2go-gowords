package lobby

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lexyblazy/gowords/internal/config"
	"github.com/lexyblazy/gowords/internal/dictionary"
	"github.com/lexyblazy/gowords/internal/game"
)

type Room struct {
	id      int
	players map[string]*Player
	d       *dictionary.Dictionary
	c       *config.Config
	// game state
	gs *game.GameState

	registerChan chan *Player

	unregisterChan chan *Player

	removeFromLobbyFunc func(moniker string)
}

type OutgoingMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type JoinRoomPayload struct {
	Moniker string `json:"moniker,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewRoom(c *config.Config, id int, d *dictionary.Dictionary, removeFromLobbyFunc func(moniker string)) *Room {

	return &Room{
		id:                  id,
		players:             make(map[string]*Player),
		d:                   d,
		c:                   c,
		registerChan:        make(chan *Player),
		unregisterChan:      make(chan *Player),
		removeFromLobbyFunc: removeFromLobbyFunc,
	}
}

func (r *Room) handlePlayerSubmission(playerId string, message []byte) {

	r.gs.SubmitWord(playerId, string(message))

}

func (r *Room) AddPlayer(player *Player) {
	r.players[player.id] = player
	go player.readPump()
	go player.writePump()
}

func (r *Room) RemovePlayer(player *Player) {
	delete(r.players, player.id)
	close(player.sendMsgCh)
	r.removeFromLobbyFunc(player.moniker)
}

func (r *Room) GetPlayerCount() int {
	return len(r.players)
}

func (r *Room) Run() {

	go r.PrintPlayers()
	r.gs = game.NewGameState(r.c, r.d, func(eventType game.EventType, payload any) {

		switch eventType {
		case game.EventTypeGeneral, game.EventTypeRoundWinner:
			r.BroadcastToAllPlayers(eventType, payload)
		case game.EventTypePlayerWordAccepted, game.EventTypePlayerWordRejected, game.EventTypePlayerRoundScores:
			r.BroadcastToPlayer(payload.(game.BasicPayload))
		case game.EventTypeGeneralExcludingPlayer:
			r.BroadcastToOtherPlayers(payload.(game.BasicPayload))
		default:
			log.Println("Unknown event type:", eventType)
		}

	})

	go func() {
		for {
			select {
			case c := <-r.registerChan:
				r.AddPlayer(c)
			case c := <-r.unregisterChan:
				r.RemovePlayer(c)
			}

		}
	}()

	r.gs.Run()

}

func (s *Room) PrintPlayers() {
	for {
		if len(s.players) > 0 {
			// log.Println("Players:", len(s.players))
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *Room) GetPlayerMoniker(playerId string) string {
	return s.players[playerId].moniker
}

func (s *Room) BroadcastToAllPlayers(eventType game.EventType, payload any) {

	outgoingMessage := OutgoingMessage{
		Type: "info",
	}

	switch eventType {
	case game.EventTypeGeneral:
		outgoingMessage.Message = payload.(game.BasicPayload).Message
	case game.EventTypeRoundWinner:
		playerId := payload.(game.RoundWinnerPayload).PlayerId
		score := payload.(game.RoundWinnerPayload).Score
		playerName := s.GetPlayerMoniker(playerId)
		outgoingMessage.Message = fmt.Sprintf("🏆 Kudos to %s for winning the round with %d points", playerName, score)

	}

	messageBytes, err := json.Marshal(outgoingMessage)

	if err != nil {
		log.Println("BroadcastToAllPlayers: Error marshalling:", err)
		return
	}

	for _, p := range s.players {
		p.SendMessage(messageBytes)
	}
}

func (s *Room) BroadcastToPlayer(payload game.BasicPayload) {

	outgoingMessage := OutgoingMessage{
		Type:    "info",
		Message: payload.Message,
	}

	messageBytes, err := json.Marshal(outgoingMessage)

	if err != nil {
		log.Println("BroadcastToPlayer: Error marshalling:", err)
		return
	}

	for _, p := range s.players {
		if p.id == payload.PlayerId {
			p.SendMessage(messageBytes)
		}
	}
}

func (s *Room) BroadcastToOtherPlayers(payload game.BasicPayload) {

	playerName := s.GetPlayerMoniker(payload.PlayerId)

	event := OutgoingMessage{
		Type:    "general",
		Message: fmt.Sprintf("%s: %s", playerName, payload.Message),
	}

	messageBytes, err := json.Marshal(event)
	if err != nil {
		log.Println("BroadcastToOtherPlayers: Error marshalling:", err)
		return
	}

	for _, p := range s.players {
		if p.id != payload.PlayerId {
			p.SendMessage(messageBytes)
		}
	}
}
