package util

import (
	"math/rand"
)

// IntRand ... 指定範囲の乱数を生成する
func IntRand(min int, max int) int {
	seed()
	return rand.Intn(max-min) + max
}
