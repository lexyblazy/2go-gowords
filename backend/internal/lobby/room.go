package lobby

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/lexyblazy/gowords/internal/config"
	"github.com/lexyblazy/gowords/internal/dictionary"
	"github.com/lexyblazy/gowords/internal/events"
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
	Type    events.EventType `json:"type"`
	Message string           `json:"message"`
	Status  string           `json:"status"`
	Payload struct {
		Moniker   string `json:"moniker"`
		Timestamp int64  `json:"timestamp"`
	} `json:"payload"`
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

	var event events.PlayerWordSubmissionEvent
	err := json.Unmarshal(message, &event)

	if err != nil {
		log.Println("Error unmarshalling player word submission event:", err)
		return
	}

	if event.Payload.PlayerId != playerId {
		log.Println("Player ID mismatch:", event.Payload.PlayerId, playerId)
		return
	}

	r.gs.SubmitWord(&event)

}

func (r *Room) AddPlayer(player *Player) {
	r.players[player.id] = player

	// broadcast the rules to the player
	var event events.GameRulesEvent
	event.Type = events.GameRules
	event.Payload.Rules = strings.Split(r.c.Game.Rules, "\n")
	event.Payload.Timestamp = time.Now().UnixMilli()
	event.Payload.SystemMoniker = r.c.Lobby.SystemMoniker
	r.BroadcastToPlayer(player.id, &event)

	// broadcast the current round info to the player immediately
	r.gs.BroadcastRoundInfo()

}

func (r *Room) RemovePlayer(player *Player) {
	delete(r.players, player.id)
	close(player.sendMsgCh)
	r.removeFromLobbyFunc(player.moniker)
}

func (r *Room) GetPlayerCount() int {
	return len(r.players)
}

func (r *Room) Broadcast(event events.EnrichableEvent) {
	destination := event.GetDestination()
	var playerName string

	playerId := event.GetPlayerID()

	if len(playerId) > 0 {
		player, exists := r.players[playerId]

		if !exists || player == nil {
			// Player no longer exists (disconnected or not yet joined)
			return
		}

		playerName = player.moniker
	}

	event.Enrich(events.EnrichmentParams{PlayerName: playerName, SystemMoniker: r.c.Lobby.SystemMoniker})

	switch destination {
	case events.EventDestinationAll:
		r.BroadcastToAll(event)
	case events.EventDestinationPlayer:
		r.BroadcastToPlayer(event.GetPlayerID(), event)
	case events.EventDestinationOtherPlayers:
		r.BroadcastToOtherPlayers(event.GetPlayerID(), event)
	}

}

func (r *Room) Run() {

	go r.PrintPlayers()
	r.gs = game.NewGameState(r.c, r.d, r.Broadcast)

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

func (s *Room) GetPlayerName(playerId string) string {
	return s.players[playerId].moniker
}

func (s *Room) BroadcastToAll(event events.EnrichableEvent) {

	messageBytes, err := json.Marshal(event)

	if err != nil {
		log.Println("BroadcastToAllPlayers: Error marshalling:", err)
		return
	}

	for _, p := range s.players {
		p.SendMessage(messageBytes)
	}
}

func (s *Room) BroadcastToPlayer(playerId string, event events.EnrichableEvent) {

	messageBytes, err := json.Marshal(event)
	if err != nil {
		log.Println("BroadcastToPlayer: Error marshalling:", err)
		return
	}

	for _, p := range s.players {
		if p.id == playerId {
			p.SendMessage(messageBytes)
		}
	}
}

func (s *Room) BroadcastToOtherPlayers(playerId string, event events.EnrichableEvent) {

	messageBytes, err := json.Marshal(event)
	if err != nil {
		log.Println("BroadcastToOtherPlayers: Error marshalling:", err)
		return
	}

	for _, p := range s.players {
		if p.id != playerId {
			p.SendMessage(messageBytes)
		}
	}
}
