package rand

import (
	mrand "math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var verifyCodeRunes = []rune("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var passwordRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
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

func RandVerifyCode(n int) string {
	doInit()
	b := make([]rune, n)
	for i := range b {
		b[i] = verifyCodeRunes[mrand.Intn(len(verifyCodeRunes))]
	}
	return string(b)
}

func RandPassword(min int, max ...int) string {
	doInit()
	n := min
	if len(max) > 0 {
		if max[0] > min {
			n = min + mrand.Intn(max[0]-min)
		}
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = passwordRunes[mrand.Intn(len(passwordRunes))]
	}
	return string(b)
}
