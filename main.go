package main

import (
	"fmt"
	"log"
	"os"
	"net/http"

	"github.com/zmb3/spotify"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserFollowRead)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	clientID := os.Getenv("SPOTIFY_ID_3")
	secretKey := os.Getenv("SPOTIFY_SECRET_KEY_3")
	auth.SetAuthInfo(clientID, secretKey)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch

	follow, _ := client.CurrentUsersFollowedArtistsOpt(50, "")
	//fmt.Println(follow.Artists)
	for _, f := range follow.Artists {
		fmt.Println(f)
	}
	last := follow.Artists[len(follow.Artists)-1].SimpleArtist.ID
	fmt.Printf("%T\n", last)
	follow, _ = client.CurrentUsersFollowedArtistsOpt(50, last.String())
	//fmt.Println(follow.Artists)
	for _, f := range follow.Artists {
		fmt.Println(f)
	}
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}