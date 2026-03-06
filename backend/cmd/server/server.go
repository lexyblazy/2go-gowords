package main

import (
	"github.com/lexyblazy/gowords/internal/config"
	"github.com/lexyblazy/gowords/internal/lobby"
	"github.com/lexyblazy/gowords/internal/server"
	"github.com/lexyblazy/gowords/internal/store"
)

func main() {
	c := config.New("config.json")
	l := lobby.New(c)
	l.Init()
	db := store.NewSqlDB(c.Db.FileName)

	s := server.New(db, c.Server.Port, l)
	s.Start()

}
