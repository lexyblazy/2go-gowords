# GoWords (2go Remake)

A modern remake of **GoWords**, the classic word game that used to run on the **2go mobile app**.

This project recreates the gameplay using a modern web stack while keeping the spirit and simplicity of the original game.

> Originally played by thousands of users on feature phones through the 2go messenger platform.

[Original 2go GoWords screenshot](./2go-example.jpg)

*(Original gameplay screenshot from the 2go GoWords game.)*

---

# What is GoWords?

GoWords is a fast-paced multiplayer word game.

Players are given a set of letters and must quickly submit valid words that can be formed using those letters before the round timer ends.

Typical round flow:

1. A set of letters is revealed.
2. Players submit words using those letters.
3. Valid words score points.
4. At the end of the round, the winner is announced.
5. A new round begins.

The game rewards:

* vocabulary
* pattern recognition
* typing speed
* quick thinking

---

# Why This Project Exists

This project was created primarily as a **fun technical challenge and nostalgia project**.

However, rebuilding GoWords on the modern web turned out to be **non-trivial** due to several factors:

* real-time multiplayer communication
* authoritative game server
* WebSocket event handling
* consistent round state synchronization
* mobile-first UI constraints
* high-frequency activity feed updates

Even though the UI appears simple, the system requires careful coordination between the frontend and backend.

---

# Architecture

The system follows a **server-authoritative real-time architecture**.

```
Client (React)
      │
      │ WebSocket
      ▼
Go Backend (authoritative game engine)
```

### Backend

* Go
* WebSocket event system
* authoritative game state
* round timers managed by the server
* event-driven state updates

The server is responsible for:

* validating word submissions
* tracking round timers
* broadcasting game events
* maintaining the canonical game state

### Frontend

* React
* TypeScript
* Tailwind CSS
* WebSocket client

The frontend acts primarily as a **state renderer**, reacting to server events rather than owning game logic.

---

# Key Features

* Real-time multiplayer gameplay
* WebSocket-based event system
* Authoritative backend game engine
* Activity feed showing player actions
* Round timers synchronized with the server
* Mobile-friendly UI
* Sound effects for gameplay feedback
* Dark mode support

---

# Tech Stack

Frontend

* React
* TypeScript
* Tailwind CSS v4
* Vite
* pnpm

Backend

* Go
* WebSockets
* Event-driven state management

Infrastructure

* Docker
* Docker Compose
* Caddy reverse proxy

---

# Running Locally

The project provides separate Docker Compose configurations for development and production.

### Development

```
docker compose -f deploy/docker-compose.dev.yml up --build
```

Frontend will be available at:

```
http://localhost:5173
```

---

# Production

```
docker compose -f deploy/docker-compose.prod.yml up -d --build
```

Caddy handles TLS and routing for the production deployment.

---

# Project Goals

The goal of this project is not to perfectly reproduce every aspect of the original 2go game, but rather to:

* recreate the gameplay experience
* modernize the architecture
* explore real-time multiplayer design on the web

---

# Status

The game is playable and deployed.

Future improvements may include:

* better word validation dictionary
* leaderboard support
* player rooms
* reconnect handling
* analytics and telemetry

---

# Acknowledgements

Inspired by the **original GoWords game on the 2go platform**.

This project is an independent remake created for learning and experimentation.

---

# License

MIT
