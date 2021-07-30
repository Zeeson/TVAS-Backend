package utils

import (
    "time"
    "math/rand"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$&")

func RandSeq() string {
	rand.Seed(time.Now().UnixNano())
    b := make([]rune, 10)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }    
	return string(b)
}