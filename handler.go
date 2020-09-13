package main

import (
	"fmt"
	"net/http"
	"log"
	"html/template"
	//"os"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login")
	url := auth.AuthURL(state)
	t := template.Must(template.ParseFiles("templates/index.html"))
	err := t.Execute(w, url)
	if err != nil {
		fmt.Println(err)
	}
	//http.Redirect(w, r, string(url), 301) 
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// use the same state string here that you used to generate the URL
	fmt.Println("/handle")
	token, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	
	//fmt.Fprintf(os.Stdout, "session : %#v\n", session.Values["token"].(Oauth2Token).flag)
	session, _ := store.Get(r, session_name)
	session.Values["token"] = Oauth2Token{Token:*token}
	err = session.Save(r, w)
	//fmt.Fprintf(os.Stdout, "session : %#v\n", session.Values["token"].(Oauth2Token))
	//fmt.Println(session)

	http.Redirect(w, r, "/home", 301)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")

	token := getTokenFromSession(r)
	client := auth.NewClient(&token)
	artistInfo, _ := getFollowingArtists(client)

	t := template.Must(template.ParseFiles("templates/home.html"))
	err := t.Execute(w, artistInfo)
	if err != nil {
		fmt.Println(err)
	}

}

func resultHander(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromSession(r)
	client := auth.NewClient(&token)

	r.ParseForm()
	artistsId := r.Form["artist"]

	recommendIds := getRecommendationId(client, artistsId)

	_, allArtistsId := getFollowingArtists(client) // TODO: use redis to get allArtistsId
	recommendedArtists := getRecommendedArtists(client, recommendIds, allArtistsId)
	
	t := template.Must(template.ParseFiles("templates/result.html"))
	err := t.Execute(w, recommendedArtists)
	if err != nil {
		fmt.Println(err)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logout")
	session, _ := store.Get(r, session_name)
	session.Options.MaxAge = -1
	err := session.Save(r, w)
	if err != nil {
		log.Fatal("failed to delete session", err)
	}
	http.Redirect(w, r, "/", 301)
}