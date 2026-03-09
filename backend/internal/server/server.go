package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"

	"github.com/lexyblazy/gowords/internal/lobby"
	"github.com/lexyblazy/gowords/internal/store"
)

type Server struct {
	db    *store.SqlDb
	rs    *store.RedisStore
	port  int
	lobby *lobby.Lobby
}

func New(db *store.SqlDb, redis *store.RedisStore, port int, l *lobby.Lobby) *Server {
	return &Server{
		db:    db,
		rs:    redis,
		port:  port,
		lobby: l,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) loadRoutes(router *http.ServeMux) {

	router.HandleFunc("/api/login", s.jsonHandler(s.login))

	router.HandleFunc("/api/register", s.jsonHandler(s.register))

	router.HandleFunc("/api/logout", s.jsonHandler(s.logout))

	router.HandleFunc("/api/reset-password", s.jsonHandler(s.resetPassword))

	router.HandleFunc("/api/leaderboards", s.jsonHandler(s.getLeaderboards))

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		var (
			moniker  string
			playerId string
		)

		token := getCookie(r, "session")
		userId, err := s.rs.Get(r.Context(), fmt.Sprintf("sessions:%s", token))

		if len(userId) > 0 {
			user, err := s.db.GetUserById(userId)

			if err == nil {
				playerId = user.ID
				moniker = user.Moniker
			}
		}

		player := lobby.NewPlayer(conn, nil, moniker, playerId, s.lobby)
		go player.ReadPump()
		player.WritePump()

	})
}

func (s *Server) jsonHandler(fn func(r *http.Request, w http.ResponseWriter) (any, int, error)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, statusCode, err := fn(r, w)

		w.WriteHeader(statusCode)

		if err != nil {
			var errorResponse ErrorResponse
			errorResponse.Error = err.Error()
			log.Println("Route Handler Error:", r.Method, r.URL.Path, err, "Status Code:", statusCode)
			if encodeErr := json.NewEncoder(w).Encode(errorResponse); encodeErr != nil {
				log.Println("Error writing error response:", encodeErr)
			}

			return
		}

		if encodeErr := json.NewEncoder(w).Encode(data); encodeErr != nil {
			log.Println("Error writing response:", encodeErr)
		}
	}

}

func (s *Server) Start() {
	mux := http.NewServeMux()

	s.loadRoutes(mux)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	log.Println("Server is running on 🚀", s.port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
