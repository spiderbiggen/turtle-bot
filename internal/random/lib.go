package random

import (
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func Intn(max int) int {
	return r.Intn(max)
}
