# GoWords Frontend

React + TypeScript + Tailwind v4 client for the GoWords real-time word game.

This frontend connects to an authoritative Go WebSocket backend.
The client is intentionally thin: it renders server events and avoids complex client-side game logic.

---

## Tech Stack

* React (functional components)
* TypeScript
* Tailwind CSS v4
* WebSocket (event-driven)
* Minimal local state
* Authoritative backend architecture

---

## Architecture Principles

### 1. Server Is Authoritative

The client:

* Does not validate words
* Does not compute scores
* Does not determine winners
* Does not track round progression

It only:

* Sends `PLAYER_WORD_SUBMISSION`
* Renders server events
* Maintains UI state (input, sound toggle, animations)

All game state is derived from server events.

---

### 2. Event-Driven State

All server messages are processed through:

```ts
getNewState(state, event)
```

The reducer:

* Produces a new state only when meaningful changes occur
* Keeps feed capped at `MAX_FEED_ITEMS`
* Avoids unnecessary object recreation where practical
* Returns previous state when no relevant changes occur

---

### 3. Component Structure

```
GameScreen
 ├── LetterBoard
 ├── ActivityFeed
 │    └── ActivityFeedItem
 └── WordInputBar
```

Design goals:

* Leaf components are memoized
* Feed rendering isolated
* Input local state only
* Minimal prop coupling

---

## Performance Decisions

### Feed Capped

```ts
feed.slice(-MAX_FEED_ITEMS)
```

Prevents unbounded growth and memory bloat.

---

### Memoization Strategy

Memoized:

* ActivityFeed
* ActivityFeedItem
* WordInputBar
* LetterBoard

Not memoized:

* GameScreen (layout container)

Reason:
Memo boundaries placed at heavy leaf components, not structural wrappers.

---

### Scroll Behavior

* Auto-scroll only depends on `feed.length`
* Scroll does not fight the user
* No virtualization (not needed at 500 max items)

---

### Rendering Reality

The app:

* Does not suffer from infinite re-renders
* Re-renders normally on 5s server round broadcasts
* Does not re-render feed on typing
* Shows smooth behavior on iPhone Safari

DevTools highlight may appear noisy, but console logs confirm actual render frequency is controlled.

---

## Sound System

Location:

```
/public/sounds/
```

Managed by:

```
src/lib/sound.ts
```

Features:

* Toggleable (default OFF)
* Word accepted
* Word rejected
* Final 10-second tick
* Round winner

Sound manager:

* Preloads audio
* Prevents overlap stacking
* No React state dependency

---

## Micro-Interactions

* Feed fade-up animation
* Final 10s urgency pulse
* Minimalist styling hierarchy
* Mobile-tight layout spacing

No animation libraries used.

---

## Mobile-First Design

Optimized for:

* iPhone 12 Pro viewport
* Reduced vertical waste
* Tight typography
* Controlled spacing
* Safe-area bottom padding

---

## Folder Structure (Simplified)

```
src/
  components/
    ActivityFeed.tsx
    ActivityFeedItem.tsx
    LetterBoard.tsx
    WordInputBar.tsx
  state/
    reducer.ts
    types.ts
  socket/
    socket.ts
  lib/
    sound.ts
  screens/
    GameScreen.tsx
```

---

## Local Development

Install:

```
npm install
```

Run:

```
npm run dev
```

Build:

```
npm run build
```

Preview production build:

```
npm run preview
```

---

## WebSocket Contract

The client expects server events such as:

* `ROUND_INFO`
* `ROUND_OVER`
* `NEXT_ROUND_COUNTDOWN`
* Word acceptance/rejection events

All timestamps use server `endsAt`.
Client applies drift correction with `addDriftTime`.

---

## State Update Philosophy

We follow:

> Simplicity > premature micro-optimization

We:

* Avoid unnecessary complexity
* Only optimize when performance issues appear
* Keep reducer readable
* Maintain predictable state transitions

---

## What This Frontend Does NOT Do

* No optimistic updates
* No client-side scoring
* No persistence (yet)
* No reconnection recovery logic (future enhancement)
* No virtualized list (not needed)

---

## Known Future Enhancements

* Reconnect handling
* Persistent username
* Latency indicator
* Round summary modal
* PWA support
* Haptic feedback on mobile
* Deployment hardening (Caddy + WebSocket proxy tuning)

---

## Production Readiness Status

* Stable event rendering
* Smooth mobile performance
* Controlled re-renders
* No render leaks
* Memory safe (feed capped)
* Clean component isolation
* Authoritative backend respected

Frontend polish phase complete.

