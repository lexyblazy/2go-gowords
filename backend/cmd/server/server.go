package main

import (
	"log"
	"os"

	"github.com/lexyblazy/gowords/internal/config"
	"github.com/lexyblazy/gowords/internal/lobby"
	"github.com/lexyblazy/gowords/internal/server"
	"github.com/lexyblazy/gowords/internal/store"
)

func main() {
	c := config.New("config.json")
	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		redisUrl = c.Redis.URL
	}
	rs, err := store.NewRedisStore(redisUrl)
	if err != nil {
		log.Fatal("Redis err", err)
	}
	dbDsn := os.Getenv("DB_DSN")

	if dbDsn == "" {
		dbDsn = c.Db.DSN
	}
	db, err := store.NewSqlDB(dbDsn)

	if err != nil {
		log.Fatal("Db err", err)
	}

	l := lobby.New(c, db)
	l.Init()

	s := server.New(db, rs, c.Server.Port, l)
	s.Start()

}
