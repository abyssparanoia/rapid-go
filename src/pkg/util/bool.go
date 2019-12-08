package util

import (
	"math/rand"
)

// BoolRand ... generate random bool value
func BoolRand(rate float32) bool {
	seed()
	return rand.Float32()*100 < rate
}
