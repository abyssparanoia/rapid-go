package util

import (
	"math/rand"
	"time"
)

func seed() {
	rand.Seed(time.Now().UnixNano())
}
