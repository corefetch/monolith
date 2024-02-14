package route

import (
	"fmt"
	"net/http"

	"github.com/olebedev/emitter"
	"learnt.io/core/rest"
	"learnt.io/modules/vcr/service"
	"learnt.io/modules/vcr/store"
	"learnt.io/modules/ws"
)

var rooms map[string]*service.Room = make(map[string]*service.Room)

func init() {
	ws.Events().On("message", func(e *emitter.Event) {
		if m, v := e.Args[0].(ws.Message); v {
			fmt.Println(m)
		}
	})
}

func Create(c *rest.Context) {}

type JoinData struct {
	Session string `json:"session"`
}

func Join(c *rest.Context) {

	var data JoinData

	if err := c.Read(&data); err != nil {
		c.Write(err, http.StatusBadRequest)
		return
	}

	room, exists := rooms[data.Session]

	if !exists {

		session, err := store.GetSession(data.Session)

		if err != nil {
			c.Write(err, http.StatusNotFound)
			return
		}

		var participantOfRoom = false
		for _, participant := range session.Participants {
			if participant.Hex() == c.User() {
				participantOfRoom = true
			}
		}

		if !participantOfRoom {
			c.Write("not a room participant", http.StatusInternalServerError)
			return
		}

		room = service.NewRoom(session)

		rooms[data.Session] = room
	}

	if err := room.Join(c.User()); err != nil {
		c.Write(err, http.StatusInternalServerError)
		return
	}
}
