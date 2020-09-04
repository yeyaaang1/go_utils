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

func RandPassword(n int) string {
	doInit()
	b := make([]rune, n)
	for i := range b {
		b[i] = passwordRunes[mrand.Intn(len(verifyCodeRunes))]
	}
	return string(b)
}
