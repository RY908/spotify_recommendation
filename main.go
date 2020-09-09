package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"net/url"

	"github.com/zmb3/spotify"
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
	client spotify.Client
)


func main() {
	auth.SetAuthInfo(clientID, secretKey)
	u := auth.AuthURL(state)
	fmt.Println(url.Parse(u))
	// first start an HTTP server
	uu, _ := url.Parse(u)
	fmt.Printf("URL: %s\n", uu.String())
  	fmt.Printf("Scheme: %s\n", uu.Scheme)
  	fmt.Printf("Opaque: %s\n", uu.Opaque)
  	fmt.Printf("User: %s\n", uu.User)
  	fmt.Printf("Host: %s\n", uu.Host)
  	fmt.Printf("Hostname(): %s\n", uu.Hostname())
  	fmt.Printf("Path: %s\n", uu.Path)
  	fmt.Printf("RawPath: %s\n", uu.RawPath)
  	fmt.Printf("RawQuery: %s\n", uu.RawQuery)
  	fmt.Printf("Fragment: %s\n", uu.Fragment)

	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", u)

	// wait for auth to complete
	//client := <-ch

	//http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/callback", redirectHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String(), client.Token)
		//log.Println(client.Token.accessToken)
		//fmt.Fprintf(w, "redirect")
		//fmt.Println(w)
		http.Redirect(w, r, string(u), 301) 
	})
	http.ListenAndServe(":8080", nil)

	fmt.Println(client.Token())

}

/*
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// use the same state string here that you used to generate the URL
	token, err := auth.Token(state, r)
	if err != nil {
		  http.Error(w, "Couldn't get token", http.StatusNotFound)
		  return
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// create a client using the specified token
	client = auth.NewClient(token)
	fmt.Println(client)
	// the client can now be used to make authenticated requests
	http.Redirect(w, r, "/", 301)
}*/
