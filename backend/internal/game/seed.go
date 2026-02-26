package game

import (
	"crypto/rand"
	"encoding/binary"
)

func NewSeed() int64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return int64(binary.LittleEndian.Uint64(b[:]))
}
