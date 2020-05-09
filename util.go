package main

import (
	"math/rand"
)

func searchSlice(k string, s *[]string) bool {
	for _, vk := range *s {
		if vk == k {
			return true
		}
	}
	return false
}

const randSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = randSource[rand.Intn(len(randSource))]
	}
	return string(b)
}
