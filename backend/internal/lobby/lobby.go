package lobby

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/lexyblazy/gowords/internal/config"
	"github.com/lexyblazy/gowords/internal/dictionary"
	"github.com/lexyblazy/gowords/internal/events"
	"github.com/lexyblazy/gowords/internal/helpers"
	"github.com/lexyblazy/gowords/internal/store"
)

type Lobby struct {
	c     *config.Config
	rooms map[int]*Room
	d     *dictionary.Dictionary
	mu    *sync.RWMutex

	monikers map[string]string

	db *store.SqlDb
	rs *store.RedisStore
}

func New(c *config.Config, db *store.SqlDb, rs *store.RedisStore) *Lobby {
	d := dictionary.NewDictionary(c.Dictionary.FileName)

	return &Lobby{
		c:        c,
		rooms:    make(map[int]*Room),
		d:        d,
		mu:       &sync.RWMutex{},
		monikers: make(map[string]string),
		db:       db,
		rs:       rs,
	}
}

func (l *Lobby) Init() {
	go l.allocateRooms()
	go l.printStats()

}

func (l *Lobby) allocateRooms() {
	for {
		if len(l.rooms) < l.c.Lobby.RoomCount {
			l.mu.Lock()
			roomId := len(l.rooms) + 1
			room := NewRoom(l.c, roomId, l.d, l.removeMoniker, l.rs, l.db)
			l.rooms[roomId] = room
			go room.Run()
			l.mu.Unlock()
		}
		time.Sleep(1 * time.Second)
	}
}

func (l *Lobby) validateMoniker(player *Player, moniker string) (bool, string) {

	// this is an existing user in our system - skip validation
	if player.id != "" && player.moniker != "" {
		l.monikers[strings.ToLower(player.moniker)] = player.id
		return true, ""
	}

	moniker = strings.TrimSpace(moniker)
	length := utf8.RuneCountInString(moniker)
	maxLength := l.c.Lobby.PlayerNameLengthMax

	if length < l.c.Lobby.PlayerNameLengthMin {
		return false, "moniker must be at least 3 characters long"
	}

	if length > maxLength {
		return false, fmt.Sprintf("moniker must be less than %d characters long", maxLength)
	}

	for _, r := range moniker {
		if unicode.IsControl(r) {
			return false, "moniker must contain only printable characters"
		}
	}

	inUseMessage := moniker + " is already in use. Please choose a different moniker."

	if moniker == l.c.Lobby.SystemMoniker {
		return false, inUseMessage
	}

	if _, ok := l.monikers[strings.ToLower(moniker)]; ok {
		return false, inUseMessage
	}

	// check if moniker belongs to any of the existing users
	user, _ := l.db.GetUserByUsername(moniker)

	if user.ID != "" && user.Moniker != "" {
		return false, inUseMessage
	}

	return true, ""
}

func toBytes(event any) ([]byte, error) {
	bytes, err := json.Marshal(event)
	if err != nil {
		return nil, errors.New("error marshalling event: " + err.Error())
	}
	return bytes, nil
}

func (l *Lobby) addMoniker(player *Player, moniker string) (*Player, error) {

	// this is an existing user in our system
	if player.id != "" && player.moniker != "" {
		l.monikers[strings.ToLower(player.moniker)] = player.id
		return player, nil
	}

	// for guest players
	uuid, err := helpers.NewUUIDV4()

	if err != nil {
		return nil, errors.New("error generating UUID")
	}

	l.monikers[strings.ToLower(moniker)] = uuid

	return &Player{
		moniker: moniker,
		id:      uuid,
	}, nil
}

func (l *Lobby) removeMoniker(moniker string) {
	delete(l.monikers, strings.ToLower(moniker))
}

func (l *Lobby) JoinRoom(player *Player, message []byte) ([]byte, error) {

	joinPayload := events.JoinRoomRequest{}
	err := json.Unmarshal(message, &joinPayload)

	if err != nil {
		return nil, err
	}

	ok, errorMessage := l.validateMoniker(player, joinPayload.Payload.PlayerName)

	if !ok {

		var joinRoomErrorEvent events.JoinRoomError
		joinRoomErrorEvent.Type = events.EventTypeJoinRoomError
		joinRoomErrorEvent.Payload.Message = errorMessage
		joinRoomErrorEvent.Payload.Timestamp = time.Now().UnixMilli()

		return toBytes(joinRoomErrorEvent)
	}

	basePlayer, err := l.addMoniker(player, joinPayload.Payload.PlayerName)

	if err != nil {
		return nil, err
	}

	var joinRoomOKEvent events.JoinRoomOK
	joinRoomOKEvent.Type = events.EventTypeJoinRoomOK
	joinRoomOKEvent.Payload.PlayerName = basePlayer.moniker
	joinRoomOKEvent.Payload.SystemMoniker = l.c.Lobby.SystemMoniker
	joinRoomOKEvent.Payload.PlayerId = basePlayer.id
	joinRoomOKEvent.Payload.Timestamp = time.Now().UnixMilli()
	joinRoomOKEvent.Payload.RoomId = 0

	// update player structs. This is redundant for existing users
	player.moniker = basePlayer.moniker
	player.id = basePlayer.id

	for i := 1; i <= len(l.rooms); i++ {

		room := l.rooms[i]

		if room.GetPlayerCount() < l.c.Lobby.MaxPlayersPerRoom {
			player.room = room
			room.registerChan <- player
			joinRoomOKEvent.Payload.RoomId = i
			break
		}
	}

	// cache moniker for both guests and existing players
	l.rs.CacheUserMoniker(context.Background(), player.id, player.moniker)

	return toBytes(joinRoomOKEvent)
}

func (l *Lobby) printStats() {
	for {

		totalPlayers := 0

		l.mu.RLock()
		for index, room := range l.rooms {
			count := room.GetPlayerCount()
			totalPlayers += count
			if count > 0 {
				fmt.Println("--------------------------------")
				fmt.Println("There are", count, "players in room #", index)
				fmt.Println("--------------------------------")
			}
		}
		l.mu.RUnlock()

		fmt.Println("--------------------------------")
		fmt.Println("There are", totalPlayers, "currently active players on the server")
		fmt.Println("--------------------------------")

		time.Sleep(60 * time.Second)
	}
}
