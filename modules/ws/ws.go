package ws

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"corefetch/core/rest"

	"github.com/gorilla/websocket"
	em "github.com/olebedev/emitter"
)

const (
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 16384
)

type Message struct {
	Kind   string         `json:"kind"`
	Data   map[string]any `json:"data"`
	Time   int64          `json:"time"`
	Source *Connection    `json:"-"`
}

type Connection struct {
	ctx  *rest.Context
	conn *websocket.Conn
	send chan Message
}

var events = &em.Emitter{}

var conns = make(map[string]*Connection)

var mux sync.Mutex

func Events() *em.Emitter {
	return events
}

func Service() *rest.Service {
	srv := rest.NewService("ws", "0.0.0")
	srv.Get("/", rest.GuardAuth(connectionHandler))
	return srv
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  16384,
	WriteBufferSize: 16384,
	CheckOrigin:     checkOrigin,
}

func (c *Connection) watch() {
	go c.readin()
	c.writeout()
}

func (c *Connection) Write(e Message) {
	c.send <- e
}

func (c *Connection) writeout() {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case event, ok := <-c.send:

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("Participant watch write closed")
				return
			}

			if err := c.conn.WriteJSON(event); err != nil {
				log.Println("failed to write event:", err)
			}

		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *Connection) readin() {

	defer c.conn.Close()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {

		var message *Message

		if err := c.conn.ReadJSON(&message); err != nil {
			break
		}

		message.Source = c

		events.Emit("message", message)

		if message.Kind == "ping" {
			c.Write(Message{Kind: "pong"})
		}
	}
}

func connectionHandler(ctx *rest.Context) {

	conn, err := upgrader.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil)

	if err != nil {
		ctx.Write(err, http.StatusInternalServerError)
		return
	}

	mux.Lock()
	defer mux.Unlock()

	c := &Connection{
		ctx:  ctx,
		conn: conn,
		send: make(chan Message),
	}

	conn.SetCloseHandler(func(code int, text string) error {
		delete(conns, ctx.User())
		events.Emit("disconnect", c, code, text)
		return nil
	})

	conns[ctx.User()] = c

	events.Emit("connect", c)

	c.watch()
}

func checkOrigin(r *http.Request) bool {
	return true
}

func Broadcast(m Message) {
	for _, c := range conns {
		c.send <- m
	}
}

func GetConnection(user string) (c *Connection, err error) {

	c, connected := conns[user]

	if !connected {
		return nil, errors.New("connection not found")
	}

	return c, nil
}
