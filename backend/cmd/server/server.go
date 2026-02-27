package main

import (
	"log"
	"net/http"

	"github.com/lexyblazy/gowords/internal/config"

	"github.com/gorilla/websocket"
	"github.com/lexyblazy/gowords/internal/lobby"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	c := config.New("config.json")
	l := lobby.New(c)
	l.Init()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		player := lobby.NewPlayer(conn, nil, "", "", l)
		go player.ReadPump()
		player.WritePump()

	})

	log.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
