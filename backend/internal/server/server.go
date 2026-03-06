package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"

	"github.com/lexyblazy/gowords/internal/store"
	"github.com/lexyblazy/gowords/internal/lobby"
)

type Server struct {
	db   *store.SqlDb
	port int
	lobby *lobby.Lobby
}

func New(db *store.SqlDb, port int, l *lobby.Lobby) *Server {
	return &Server{
		db:   db,
		port: port,
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

	router.HandleFunc("/login", s.jsonHandler(s.login))

	router.HandleFunc("/register", s.jsonHandler(s.register))

	router.HandleFunc("/logout", s.jsonHandler(s.logout))

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		player := lobby.NewPlayer(conn, nil, "", "", s.lobby)
		go player.ReadPump()
		player.WritePump()

	})
}

func (s *Server) jsonHandler(fn func(r *http.Request) (any, int, error)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, statusCode, err := fn(r)

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

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
