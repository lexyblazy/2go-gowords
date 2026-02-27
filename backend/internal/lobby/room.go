package lobby

import (
	"encoding/json"
	"log"
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

	r.gs.SubmitWord(playerId, string(message))

}

func (r *Room) AddPlayer(player *Player) {
	r.players[player.id] = player
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
	var moniker string

	playerId := event.GetPlayerID()
	if playerId == "" {
		moniker = "System 🤖🤖🤖"
	} else {
		moniker = r.players[playerId].moniker
	}

	event.Enrich(moniker)

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

func (s *Room) GetPlayerMoniker(playerId string) string {
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
