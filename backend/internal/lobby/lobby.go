package lobby

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lexyblazy/gowords/internal/config"
	"github.com/lexyblazy/gowords/internal/dictionary"
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

func (l *Lobby) makeEvent(status string, message string) []byte {

	eventBytes, err := json.Marshal(OutgoingMessage{
		Type:    "joinRoom",
		Status:  status,
		Message: message,
	})

	if err != nil {
		return nil
	}

	return eventBytes
}

func (l *Lobby) addMoniker(conn *websocket.Conn) ([]byte, error, *Player) {

	_, message, err := conn.ReadMessage()
	if err != nil {
		return nil, err, nil
	}

	joinPayload := JoinRoomPayload{}
	err = json.Unmarshal(message, &joinPayload)

	if err != nil {
		return nil, err, nil
	}

	moniker := strings.TrimSpace(joinPayload.Moniker)

	if len(moniker) > 16 {
		return l.makeEvent("error", "moniker must be less than 16 characters long"), nil, nil
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9-_]+$`).MatchString(moniker) {
		return l.makeEvent("error", "moniker must contain only letters and numbers"), nil, nil
	}

	if _, ok := l.monikers[moniker]; ok {
		return l.makeEvent("error", moniker+" is already in use. Please choose a different moniker."), nil, nil
	}

	uuid, err := newUUID()

	if err != nil {
		return l.makeEvent("error", "error generating UUID"), nil, nil
	}

	l.monikers[moniker] = uuid

	return l.makeEvent("success", "You joined the room as: "+moniker), nil, &Player{
		moniker: moniker,
		id:      uuid,
	}
}

func (l *Lobby) removeMoniker(moniker string) {
	delete(l.monikers, moniker)
}

func (l *Lobby) JoinRoom(conn *websocket.Conn) ([]byte, error) {

	joinMessage, err, basePlayer := l.addMoniker(conn)

	if basePlayer == nil {
		return joinMessage, err
	}

	for i := 1; i <= len(l.rooms); i++ {

		room := l.rooms[i]

		if room.GetPlayerCount() < l.c.Lobby.MaxPlayersPerRoom {
			room.registerChan <- NewPlayer(conn, room, basePlayer.moniker, basePlayer.id)
			return joinMessage, nil
		}
	}

	return l.makeEvent("error", "all rooms are full"), nil
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
