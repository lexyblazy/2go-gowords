package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lexyblazy/gowords/internal/events"
)

type Client struct {
	conn      *websocket.Conn
	sendMsgCh chan []byte
	moniker   string
	playerId  string
}

func (c *Client) Run() error {

	log.Println("Connecting...")

	if c.conn == nil {

		conn, err := createConnection()

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
	}

	retryDelay := 1 * time.Second
	maxDelay := 30 * time.Second

	for {

		if err := c.ReceiveMessages(); err != nil {

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
		log.Println(string(message))
	}

	return nil

}

func (c *Client) ReceiveMessages() error {

	for {

		messageType, message, err := c.conn.ReadMessage()

		if err != nil {
			log.Println("ReadMessages:", err)

			if strings.Contains(err.Error(), "connection reset") {
				log.Println("Connection reset: Handle reconnection")
			} else {
				c.Close()
			}

			return err

		}

		if messageType == websocket.TextMessage {

			c.handleMessage(message)
		} else {
			log.Println("Unknown message type:", messageType, websocket.PongMessage, websocket.PingMessage)
		}

	}

}

func createConnection() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func main() {

	moniker := flag.String("moniker", "", "Moniker for the client")
	flag.Parse()

	if moniker == nil || *moniker == "" {
		log.Println("Moniker is required")
		return
	}

	client := &Client{
		conn:      nil,
		sendMsgCh: make(chan []byte, 1024),
		moniker:   *moniker,
	}

	go func() {
		if err := client.Run(); err != nil {
			log.Println("Error running client:", err)
		}

	}()

	// take input from stdin and send it to the server
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			err := scanner.Err()
			if err != nil {
				log.Println("Error reading input:", err)
				continue
			}

			var playerWordSubmissionEvent events.PlayerWordSubmissionEvent
			playerWordSubmissionEvent.Type = events.PlayerWordSubmission
			playerWordSubmissionEvent.Payload.PlayerId = client.playerId
			playerWordSubmissionEvent.Payload.Word = scanner.Text()

			playerWordSubmissionMessage, err := json.Marshal(playerWordSubmissionEvent)
			if err != nil {
				log.Println("Error marshalling player word submission message payload:", err)
				continue
			}

			client.sendMsgCh <- playerWordSubmissionMessage
		}
	}()

	for message := range client.sendMsgCh {
		if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error writing message:", err)
			continue
		}
	}

}
