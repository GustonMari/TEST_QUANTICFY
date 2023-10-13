package main

import (
	"fmt"

)

const (
	chars     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	domain    = "gmail.com" // Change this to your desired domain
	emailSize = 10            // Change this to the desired length of the username part
)


func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[random.Intn(len(chars))]
	}
	return string(b)
}

func randomEmail() string {
	username := randomString(emailSize)
	return fmt.Sprintf("%s@%s", username, domain)
}
