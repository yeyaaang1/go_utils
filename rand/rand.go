package rand

import (
	mrand "math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var initFlag = false

func doInit() {
	if !initFlag {
		mrand.Seed(time.Now().UnixNano())
	}
	initFlag = true
}

func RandRunes(n int) string {
	doInit()
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[mrand.Intn(len(letterRunes))]
	}
	return string(b)
}