package util

import (
	"math/rand"
)

// IntRand ... generate random int value
func IntRand(min int, max int) int {
	seed()
	return rand.Intn(max-min) + max
}
