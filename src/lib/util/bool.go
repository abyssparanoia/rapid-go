package util

import (
	"math/rand"
)

// BoolRand ... 指定確率でbool値を生成する
func BoolRand(rate float32) bool {
	seed()
	return rand.Float32()*100 < rate
}
