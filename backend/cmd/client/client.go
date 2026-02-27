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
}

type IncomingMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Status  string `json:"status,omitempty"`
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
		joinRoomRequestEvent.Payload.Moniker = c.moniker

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

	log.Println(string(message))

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
			client.sendMsgCh <- []byte(scanner.Text())

		}
	}()

	for message := range client.sendMsgCh {
		if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error writing message:", err)
			continue
		}
	}

}
