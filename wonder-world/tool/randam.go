package tool

import (
	"math/rand"
	"time"
)

func Randam() int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	return r.Intn(100000)
}
