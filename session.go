package main

import (
	"net/http"
	"golang.org/x/oauth2"
)

func getTokenFromSession(r *http.Request) oauth2.Token {
	session, _ := store.Get(r, "cookie-name")
	token := session.Values["token"].(Oauth2Token).Token
	return token
}