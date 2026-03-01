# GoWords Backend

Real-time multiplayer word game engine built in Go.

The server is authoritative.
Clients connect via WebSocket and receive structured JSON events.

Designed for:

* Deterministic gameplay
* Room isolation
* Configurable behavior
* Low memory footprint
* Clean concurrency model
* Real-time multiplayer

---

# 🏗 Architecture Overview

```text
Lobby
 └── Room
       └── GameState
              ├── Round Generator (buffered)
              ├── Round Loop
              └── GameRound
```

## Core Principles

* Server authoritative for all gameplay logic
* Stateless regarding activity feed
* Event-driven communication
* Per-room isolation
* Pre-generated round buffering
* Configurable via JSON config

---

# 🔌 Networking

* Transport: WebSocket
* Library: `gorilla/websocket`
* Endpoint: `/ws`
* JSON event protocol
* Fire-and-forget broadcast model

The server does **not** store feed history.
Clients construct their own activity feed from incoming events.

---

# 🎮 Game Flow

## Room Lifecycle

* Rooms are preallocated at service start
* Each room starts its GameState immediately
* Game runs continuously
* Players may join mid-round
* No host / ready-state coordination

Rooms behave as persistent live arenas.

---

## Round Lifecycle

Each round:

1. Two base words selected from dictionary
2. Valid word set precomputed
3. Round runs for configured duration
4. Expansion letters added near end
5. Word submissions validated in real time
6. Scores computed at round end
7. Winner announced
8. Countdown before next round

---

## GameRound Structure

```go
type GameRound struct {
    words                   []dictionary.Word
    expansionWord           dictionary.Word
    validWords              map[string]struct{}
    validWordsWithExpansion map[string]struct{}
    seenWords               map[string]struct{}
}
```

* `validWords` used during base phase
* `validWordsWithExpansion` used after expansion
* `seenWords` prevents duplicate scoring
* Maps are reference types (no deep copies)

Expired rounds are garbage collected.

---

# 🔄 Round Generation

Each `GameState` maintains a buffered round pool.

```go
func (gs *GameState) RefillRounds() {
    for {
        r := gs.NewRound()

        if gs.box.GetDistinctCharacterCount(r.words...) >= threshold {
            gs.rounds <- r // blocks automatically when full
        }
    }
}
```

## Design Characteristics

* Default buffer size configurable (e.g., 50)
* Backpressure handled via channel capacity
* No polling required
* Memory footprint scales with buffer size

Reducing `roundCount` reduces memory usage linearly.

---

# 📚 Dictionary

* Source: SCOWL word list
* ASCII words only
* Minimum word length: 3
* Loaded at startup
* Preprocessed into:

```go
type Word struct {
    Text    string
    FreqMap [26]uint8
}
```

All submissions are validated against dictionary and letter frequency.

---

# 📡 Event System

All outbound communication uses structured events.

## Event Interface

```go
type EnrichableEvent interface {
    GetType() EventType
    GetDestination() EventDestination
    GetPlayerID() string
    Enrich(moniker string)
}
```

## Event Destinations

```go
const (
    EventDestinationAll
    EventDestinationPlayer
    EventDestinationOtherPlayers
)
```

### Responsibilities

* Game logic constructs events
* Room enriches events (moniker + timestamp)
* Room handles routing only
* Player handles WebSocket transport

---

# 👤 Player Model

```go
type Player struct {
    id       string
    moniker  string
    conn     *websocket.Conn
    sendCh   chan []byte
    lobby    *Lobby
    room     *Room
}
```

## Player Concurrency

Each player runs:

* `readPump()` – inbound messages
* `writePump()` – outbound messages + heartbeat ping

Dead connections are cleaned up automatically.

---

# 🏟 Lobby

Lobby:

* Owns rooms
* Assigns players to available room
* Enforces moniker uniqueness
* Enforces max players per room
* Validates moniker length

Monikers may include emojis.
Player IDs are server-generated UUIDs.

---

# ⏱ Time Model

Server is authoritative.

Server sends absolute timestamps:

```json
{
  "endsAt": 1709135000000
}
```

Clients compute countdown locally from absolute timestamp.

Server may rebroadcast round info periodically to re-anchor client clocks.

Clients never determine scoring or round completion.

---

# ⚙️ Configuration

All behavior is configurable via JSON file.

Example: `config.example.json`

```json
{
  "server": {
    "port": 8080
  },
  "lobby": {
    "roomCount": 1,
    "maxPlayersPerRoom": 100,
    "playerNameLengthMin": 3,
    "playerNameLengthMax": 20,
    "systemMoniker": "🤖Game Master🤖"
  },
  "game": {
    "roundDurationSeconds": 120,
    "roundIntervalSeconds": 20,
    "printRoundIntervalSeconds": 5,
    "rules": "You are given a set of letters...",
    "wordLength": 7,
    "wordCount": 2,
    "distinctCharacterCount": 10,
    "roundCount": 50
  },
  "dictionary": {
    "fileName": "dictionary.txt"
  }
}
```

---

## Configuration Highlights

* `roomCount` scales horizontal capacity
* `roundCount` directly impacts memory footprint
* `distinctCharacterCount` affects round difficulty
* `printRoundIntervalSeconds` affects client drift correction
* `systemMoniker` controls system event identity

No recompilation required for gameplay tuning.

---

# 🎮 Play Via CLI

In addition to the Web UI, a lightweight CLI client is provided.

Located in:

```text
cmd/client/
```

## Running the CLI Client

Start the server first:

```bash
go run cmd/server/main.go
```

Then in another terminal:

```bash
go run cmd/client/main.go
```

The CLI client:

* Connects via WebSocket
* Prompts for moniker
* Displays round letters
* Prints event stream
* Accepts word submissions via stdin

The CLI client is useful for:

* Backend testing
* Load testing
* Debugging event flow
* Verifying round logic without UI

---

# 🧵 Concurrency Model

Per room:

* 1 GameState loop goroutine
* 1 round refill goroutine
* 1 readPump + 1 writePump per player

Communication via channels:

* Round buffering
* Event routing
* Player send queues

No shared mutable state across rooms.

---

# 🧠 Memory Management

Optimizations:

* Round buffer limited (default: 50)
* No server-side feed storage
* Map reassignment (no deep copies)
* Garbage collection of expired rounds

Typical footprint:
~30MB with 50 cached rounds.

Memory scales with:

* `roundCount`
* `roomCount`
* Dictionary size

---

# 🚀 Running

```bash
go run cmd/server/main.go
```

WebSocket endpoint:

```
ws://localhost:8080/ws
```


---

# 🔒 Authoritative Guarantees

Server validates:

* Word length
* Dictionary existence
* Letter frequency
* Duplicate submissions
* Round timing
* Score computation

Clients cannot influence scoring logic.

---

# 🎯 Design Philosophy

* Deterministic backend
* Event-driven architecture
* Stateless feed handling
* Competitive tension via activity stream
* Configurable gameplay
* Memory-conscious buffering
* Room-level isolation
* Clean concurrency model

