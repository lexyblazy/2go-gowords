package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lexyblazy/gowords/internal/bot"
	"github.com/lexyblazy/gowords/internal/dictionary"
	"github.com/lexyblazy/gowords/internal/events"
)

const SERVER_URL = "localhost:8080"

type Client struct {
	conn      *websocket.Conn
	sendMsgCh chan []byte
	moniker   string
	playerId  string
	password  string
	bot       *bot.Bot
	headers   http.Header
}

func (c *Client) login() error {
	var params map[string]string = make(map[string]string)
	params["username"] = c.moniker
	params["password"] = c.password

	reqBody, _ := json.Marshal(params)

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s/login", SERVER_URL),
		bytes.NewBuffer(reqBody),
	)

	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, res.Body)
		return errors.New("Failed to login")
	}

	cookieHeader := ""
	for _, c := range res.Cookies() {
		cookieHeader += c.Name + "=" + c.Value + "; "
	}
	header := http.Header{}
	header.Set("Cookie", cookieHeader)

	c.headers = header

	return nil

}

func (c *Client) Run() error {

	// attempt a login
	log.Println("Logging in..")
	err := c.login()

	if err != nil {
		log.Fatal(err)
	}

	err = c.ConnectAndJoinRoom()

	if err != nil {
		log.Fatal(err)
	}

	retryDelay := 1 * time.Second
	maxDelay := 30 * time.Second

	for {

		if err := c.ReceiveMessages(); err != nil {

			log.Println("ReceiveMessagesErr:", err)
			if err = c.Reconnect(); err != nil {
				log.Println("ReconnectErr:", err)
			}
			time.Sleep(min(retryDelay, maxDelay))
			retryDelay *= 2
			continue
		}
	}

}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) handleMessage(message []byte) error {

	var incomingMessage events.JoinRoomOK
	err := json.Unmarshal(message, &incomingMessage)
	if err != nil {
		log.Println("Error unmarshalling incoming message:", err)
		return err
	}

	switch incomingMessage.Type {
	case events.EventTypeJoinRoomOK:
		c.playerId = incomingMessage.Payload.PlayerId
		log.Println("Joined Room #", incomingMessage.Payload.RoomId, "as", incomingMessage.Payload.PlayerName)
	default:
		if c.bot != nil {
			c.bot.HandleMessage(message)
		} else {
			log.Println(string(message))
		}
	}

	return nil

}

func (c *Client) ConnectAndJoinRoom() error {
	conn, err := createConnection(c.headers)

	if err != nil {
		log.Println("Error creating connection:", err)
		return err
	}

	c.conn = conn

	var joinRoomRequestEvent events.JoinRoomRequest
	joinRoomRequestEvent.Type = events.EventTypeJoinRoomRequest
	joinRoomRequestEvent.Payload.PlayerName = c.moniker

	joinMessage, err := json.Marshal(joinRoomRequestEvent)

	if err != nil {
		log.Println("Error marshalling join message payload:", err)
		return err
	}

	c.sendMsgCh <- joinMessage

	return nil
}

func (c *Client) Reconnect() error {
	log.Println("Reconnecting...")
	return c.ConnectAndJoinRoom()
}

func (c *Client) ReceiveMessages() error {

	for {

		messageType, message, err := c.conn.ReadMessage()

		if err != nil {
			log.Println("ReadMessages:", err)

			c.Close()

			return err

		}

		if messageType == websocket.TextMessage {

			c.handleMessage(message)
		} else {
			log.Println("Unknown message type:", messageType, websocket.PongMessage, websocket.PingMessage)
		}

	}

}

func createConnection(headers http.Header) (*websocket.Conn, error) {
	wsUrl := fmt.Sprintf("ws://%s/ws", SERVER_URL)
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, headers)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Client) submitWord(text string) {

	if c.playerId == "" {
		log.Println("Player ID is required")
		return
	}

	var playerWordSubmissionEvent events.PlayerWordSubmissionEvent
	playerWordSubmissionEvent.Type = events.PlayerWordSubmission
	playerWordSubmissionEvent.Payload.PlayerId = c.playerId
	playerWordSubmissionEvent.Payload.Word = text

	playerWordSubmissionMessage, err := json.Marshal(playerWordSubmissionEvent)
	if err != nil {
		log.Println("Error marshalling player word submission message:", err)
		return
	}
	c.sendMsgCh <- playerWordSubmissionMessage
}

func (c *Client) readInput() {

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			log.Println("Error reading input:", err)
			continue
		}

		c.submitWord(scanner.Text())
	}
}

func main() {

	moniker := flag.String("moniker", "", "Moniker for the client")
	botMode := flag.Bool("bot", false, "Run as a bot")
	password := flag.String("password", "", "password for existing users")
	flag.Parse()

	if moniker == nil || *moniker == "" {
		log.Println("Moniker is required")
		return
	}

	if password == nil || *password == "" {
		log.Println("password is required")
		return
	}

	client := &Client{
		conn:      nil,
		sendMsgCh: make(chan []byte, 1024),
		moniker:   *moniker,
		password:  *password,
	}

	if *botMode {
		bot := bot.NewBot(dictionary.NewDictionary("dictionary.txt"), client.submitWord)
		client.bot = bot
		go bot.PlayRound()
	} else {
		// take input from stdin and send it to the server
		go client.readInput()
	}

	go func() {
		if err := client.Run(); err != nil {
			log.Println("Error running client:", err)
		}

	}()

	for message := range client.sendMsgCh {
		if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error writing message:", err)
			continue
		}
	}

}
