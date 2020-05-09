package main

import (
	"time"
)

var (
	validAccessKey []string = make([]string, 0)
	issuedToken    []*Token = make([]*Token, 0)
)

const expireDuration time.Duration = time.Minute * 30

//Token is issued to approved users.
type Token struct {
	Token      string    //Token string
	IssuedTime time.Time //Issued datetime
}

func newToken() *Token {
	tok := Token{randString(32), time.Now()}
	issuedToken = append(issuedToken, &tok)
	return &tok
}

func checkExpired() {
	iteratedWhole := false
	now := time.Now()

	for !iteratedWhole {
		for i, t := range issuedToken {
			if now.Sub(t.IssuedTime) > expireDuration {
				issuedToken = append(issuedToken[0:i], issuedToken[i+1:]...)
				break
			}
			if i == len(issuedToken)-1 {
				iteratedWhole = true
			}
		}
	}
}

func validateToken(token string) bool {
	checkExpired()
	for _, t := range issuedToken {
		if token == t.Token {
			return true
		}
	}
	return false
}
