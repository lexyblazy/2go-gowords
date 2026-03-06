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
		RoomCount           int    `json:"roomCount"`
		MaxPlayersPerRoom   int    `json:"maxPlayersPerRoom"`
		PlayerNameLengthMin int    `json:"playerNameLengthMin"`
		PlayerNameLengthMax int    `json:"playerNameLengthMax"`
		SystemMoniker       string `json:"systemMoniker"`
	} `json:"lobby"`
	Game struct {
		Rules                     string `json:"rules"`
		RoundDurationSeconds      int    `json:"roundDurationSeconds"`
		RoundIntervalSeconds      int    `json:"roundIntervalSeconds"`
		PrintRoundIntervalSeconds int    `json:"printRoundIntervalSeconds"`
		WordLength                int    `json:"wordLength"`
		WordCount                 int    `json:"wordCount"`
		DistinctCharacterCount    int    `json:"distinctCharacterCount"`
		RoundCount                int    `json:"roundCount"` // number of rounds that must be available in the rounds channel
	} `json:"game"`
	Dictionary struct {
		FileName string `json:"fileName"`
	} `json:"dictionary"`
	Db struct {
		FileName string `json:"fileName"`
	} `json:"db"`
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
