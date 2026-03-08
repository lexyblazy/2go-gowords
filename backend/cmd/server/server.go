package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lexyblazy/gowords/internal/config"
	"github.com/lexyblazy/gowords/internal/lobby"
	"github.com/lexyblazy/gowords/internal/server"
	"github.com/lexyblazy/gowords/internal/store"
)

func mustRedis(c *config.Config) *store.RedisStore {
	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		redisUrl = c.Redis.URL
	}
	rs, err := store.NewRedisStore(redisUrl)
	if err != nil {
		log.Fatal("Redis err", err)
	}

	return rs
}

func mustDb(c *config.Config) *store.SqlDb {
	dbDsn := os.Getenv("DB_DSN")

	if dbDsn == "" {
		dbDsn = c.Db.DSN
	}

	db, err := store.NewSqlDB(dbDsn)

	if err != nil {
		log.Fatal("Db err", err)
	}

	return db

}

func main() {
	c := config.New("config.json")
	rs := mustRedis(c)
	db := mustDb(c)

	l := lobby.New(c, db, rs)
	l.Init()

	s := server.New(db, rs, c.Server.Port, l)
	go s.Start()

	waitForShutdown(db, rs)

}

func waitForShutdown(db *store.SqlDb, rs *store.RedisStore) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-sig

	db.Close()
	rs.Close()

	os.Exit(0)
}
