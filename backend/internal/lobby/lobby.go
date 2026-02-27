package lobby

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	// "github.com/gorilla/websocket"
	"github.com/lexyblazy/gowords/internal/config"
	"github.com/lexyblazy/gowords/internal/dictionary"
	"github.com/lexyblazy/gowords/internal/events"
)

type Lobby struct {
	c     *config.Config
	rooms map[int]*Room
	d     *dictionary.Dictionary
	mu    *sync.RWMutex

	monikers map[string]string
}

func New(c *config.Config) *Lobby {
	d := dictionary.NewDictionary(c.Dictionary.FileName)

	return &Lobby{
		c:        c,
		rooms:    make(map[int]*Room),
		d:        d,
		mu:       &sync.RWMutex{},
		monikers: make(map[string]string),
	}
}

func (l *Lobby) Init() {
	go l.allocateRooms()
	// go l.printStats()

}

func (l *Lobby) allocateRooms() {
	for {
		if len(l.rooms) < l.c.Lobby.RoomCount {
			l.mu.Lock()
			roomId := len(l.rooms) + 1
			room := NewRoom(l.c, roomId, l.d, l.removeMoniker)
			l.rooms[roomId] = room
			go room.Run()
			l.mu.Unlock()
		}
		time.Sleep(1 * time.Second)
	}
}

func (l *Lobby) validateMoniker(moniker string) (bool, string) {

	moniker = strings.TrimSpace(moniker)
	length := utf8.RuneCountInString(moniker)

	if length < 3 {
		return false, "moniker must be at least 3 characters long"
	}

	if length > 16 {
		return false, "moniker must be less than 16 characters long"
	}

	for _, r := range moniker {
		if unicode.IsControl(r) {
			return false, "moniker must contain only printable characters"
		}
	}

	if _, ok := l.monikers[moniker]; ok {
		return false, moniker + " is already in use. Please choose a different moniker."
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

func (l *Lobby) addMoniker(moniker string) (*Player, error) {

	uuid, err := newUUID()

	if err != nil {
		return nil, errors.New("error generating UUID")
	}

	l.monikers[moniker] = uuid

	return &Player{
		moniker: moniker,
		id:      uuid,
	}, nil
}

func (l *Lobby) removeMoniker(moniker string) {
	delete(l.monikers, moniker)
}

func (l *Lobby) JoinRoom(player *Player, message []byte) ([]byte, error) {

	joinPayload := events.JoinRoomRequest{}
	err := json.Unmarshal(message, &joinPayload)

	if err != nil {
		return nil, err
	}

	ok, errorMessage := l.validateMoniker(joinPayload.Payload.Moniker)

	if !ok {

		var joinRoomErrorEvent events.JoinRoomError
		joinRoomErrorEvent.Type = events.EventTypeJoinRoomError
		joinRoomErrorEvent.Payload.Message = errorMessage
		joinRoomErrorEvent.Payload.Timestamp = time.Now().UnixMilli()

		return toBytes(joinRoomErrorEvent)
	}

	basePlayer, err := l.addMoniker(joinPayload.Payload.Moniker)

	if err != nil {
		return nil, err
	}

	var joinRoomOKEvent events.JoinRoomOK
	joinRoomOKEvent.Type = events.EventTypeJoinRoomOK
	joinRoomOKEvent.Payload.Moniker = basePlayer.moniker
	joinRoomOKEvent.Payload.PlayerId = basePlayer.id
	joinRoomOKEvent.Payload.Timestamp = time.Now().UnixMilli()
	joinRoomOKEvent.Payload.RoomId = 0

	// update player structs
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

	return toBytes(joinRoomOKEvent)
}

func (l *Lobby) printStats() {
	for {

		l.mu.RLock()
		for index, room := range l.rooms {
			count := room.GetPlayerCount()
			if count > 0 {
				fmt.Println("--------------------------------")
				fmt.Println("There are", count, "players in room", index)
				fmt.Println("--------------------------------")
			}
		}
		l.mu.RUnlock()

		time.Sleep(10 * time.Second)
	}
}

func newUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Set version (4) and variant bits (RFC 4122)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:16],
	), nil
}
