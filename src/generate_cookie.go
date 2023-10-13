package main

import (
	"math/rand"
	"time"
)

const (
	cookieSize  = 20 // Change this to the desired length of the Cookie
	cookieChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	randomSource = rand.NewSource(time.Now().UnixNano())
	random       = rand.New(randomSource)
)

func randomCookie() string {
	cookie := make([]byte, cookieSize)

	for i := range cookie {
		cookie[i] = cookieChars[random.Intn(len(cookieChars))]
	}

	return string(cookie)
}
