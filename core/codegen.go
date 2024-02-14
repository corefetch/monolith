package core

import (
	"github.com/dchest/uniuri"
)

func NewShortKey() string {
	var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	return uniuri.NewLenChars(6, StdChars)
}
