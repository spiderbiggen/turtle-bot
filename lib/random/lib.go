package random

import (
	"math/rand"
	"time"
)

func Intn(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max)
}
