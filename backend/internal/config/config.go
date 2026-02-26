package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
	Lobby struct {
		RoomCount         int `json:"roomCount"`
		MaxPlayersPerRoom int `json:"maxPlayersPerRoom"`
	} `json:"lobby"`
	Game struct {
		TimeLimit int `json:"timeLimit"`
		RoundIntervalSeconds int `json:"roundIntervalSeconds"`
		WordLength int `json:"wordLength"`
		WordCount int `json:"wordCount"`
		DistinctCharacterCount int `json:"distinctCharacterCount"`
	} `json:"game"`
	Dictionary struct {
		FileName string `json:"fileName"`
	} `json:"dictionary"`
}

func New(path string) *Config {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var config Config
	json.NewDecoder(file).Decode(&config)
	return &config
}
