package main

import (
	"fmt"
	"log"
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
	key = []byte("spotify_access_token")
    store = sessions.NewCookieStore(key)
)

type Oauth2Token struct {
	Token oauth2.Token
}


func main() {
	gob.Register(Oauth2Token{})

	// セッション初期処理
	//store.SessionInit()

	//var client *spotify.Client
	auth.SetAuthInfo(clientID, secretKey)
	u := auth.AuthURL(state)
	// first start an HTTP server

	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", u)

	// wait for auth to complete
	//client := <-ch

	r := mux.NewRouter()
	//http.HandleFunc("/callback", completeAuth)
	r.HandleFunc("/callback", redirectHandler)
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, string(u), 301) 
	})
	r.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		token := getTokenFromSession(r)
		client := auth.NewClient(&token)
		user, err := client.CurrentUser()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("You are logged in as:", user.ID)
	})
	// rを割当
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
	//client = <-ch
	//fmt.Println(client.Token())

}
