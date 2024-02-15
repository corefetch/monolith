package ws

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"corefetch/core/rest"

	"github.com/gorilla/websocket"
)

const (
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 16384
)

type Handler func(Message)
type ConnectionHandler func(*Connection)

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

var handlers = make([]Handler, 0)

var conns = make(map[string]*Connection)

var mux sync.Mutex

var connectHandler ConnectionHandler = func(c *Connection) {}

var closeHandler ConnectionHandler = func(c *Connection) {}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  16384,
	WriteBufferSize: 16384,
	CheckOrigin:     checkOrigin,
}

func Service() *rest.Service {
	srv := rest.NewService("ws", "0.0.0")
	srv.Get("/", rest.GuardAuth(connectionHandler))
	return srv
}

func (c *Connection) Context() *rest.Context {
	return c.ctx
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

		for _, h := range handlers {
			h(*message)
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
		closeHandler(c)
		return nil
	})

	conns[ctx.User()] = c
	connectHandler(c)
	c.watch()
}

func checkOrigin(r *http.Request) bool {
	return true
}

func SetOnConnectHandler(h ConnectionHandler) {
	connectHandler = h
}

func SetOnCloseHandler(h ConnectionHandler) {
	closeHandler = h
}

func OnMessage(h Handler) {
	handlers = append(handlers, h)
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
