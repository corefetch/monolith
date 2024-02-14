package core

import "learnt.io/core/rest"

type M map[string]any

type Coordinates struct {
	Lng float64 `json:"lng" bson:"x"`
	Lat float64 `json:"lat" bson:"y"`
}

type GeoLocation struct {
	Type        string       `json:"type" bson:"type"`
	Coordinates *Coordinates `json:"coordinates" bson:"coordinates"`
}

type UserRole byte

const (
	RoleRoot UserRole = iota
	RoleAdmin
	RoleTutor
	RoleStudent
	RoleAffiliate
)

func NotImplemented(c *rest.Context) {
	c.Write([]byte("Not Implemented"))
}
