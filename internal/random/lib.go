package random

import (
	"math/rand"
	"time"
)

func Random() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func Intn(max int) int {
	return Random().Intn(max)
}
