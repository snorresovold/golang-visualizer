package utils

import (
	"math/rand"
	"time"
)

func Rand32() float32 {
	rand.Seed(time.Now().UnixNano())
	min := -1.0
	max := 1.0
	x := min + rand.Float64()*(max-min)
	return float32(x)
}
