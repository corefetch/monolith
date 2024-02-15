package service

import (
	"corefetch/modules/vcr/store"
	"corefetch/modules/ws"
)

type Participant struct {
	user string
	conn *ws.Connection
}

type Room struct {
	session      *store.Session
	participants []*Participant
}

func NewRoom(session *store.Session) (room *Room) {
	return &Room{
		session:      session,
		participants: make([]*Participant, 0),
	}
}

func (r *Room) Join(user string) (err error) {

	conn, err := ws.UserConnection(user)

	if err != nil {
		return err
	}

	r.participants = append(r.participants, &Participant{
		user: user,
		conn: conn,
	})

	return nil
}

// verify is user is present in the room
func (r *Room) Present(user string) bool {
	for _, participant := range r.participants {
		if participant.user == user {
			return true
		}
	}
	return false
}
