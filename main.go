package main

import (
	//"fmt"
	//"log"
	"os"
	"net/http"
	//"net/url"
	"encoding/gob"
	"golang.org/x/oauth2"
	"github.com/zmb3/spotify"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
	clientID = os.Getenv("SPOTIFY_ID_3")
	secretKey = os.Getenv("SPOTIFY_SECRET_3")
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserFollowRead)
	ch    = make(chan *spotify.Client)
	state = "abc123"
	//client spotify.Client
	session_name = "spotify_access_token"
	//key = []byte(session_name)
	//store = sessions.NewCookieStore(key)
	store *sessions.CookieStore
	session *sessions.Session
	conn = Connection()
)

type Oauth2Token struct {
	Token oauth2.Token
}

func main() {
	gob.Register(Oauth2Token{})
	defer conn.Close()

	// セッション初期処理
	sessionInit()

	auth.SetAuthInfo(clientID, secretKey)

	//fmt.Println("Please log in to Spotify by visiting the following page in your browser:", u)

	r := mux.NewRouter()
	r.HandleFunc("/callback", redirectHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/home", homeHandler)
	r.HandleFunc("/logout", logoutHandler)
	// rを割当
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)

}
