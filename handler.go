package main

import (
	"fmt"
	"net/http"
	"log"
	//"os"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := auth.AuthURL(state)
	http.Redirect(w, r, string(url), 301) 
}

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
	
	//fmt.Fprintf(os.Stdout, "session : %#v\n", session.Values["token"].(Oauth2Token).flag)
	session, _ := store.Get(r, session_name)
	session.Values["token"] = Oauth2Token{*token}
	err = session.Save(r, w)
	//fmt.Fprintf(os.Stdout, "session : %#v\n", session.Values["token"].(Oauth2Token))
	//fmt.Println(session)
	fmt.Println("err: ", err)

	http.Redirect(w, r, "/home", 301)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromSession(r)
	client := auth.NewClient(&token)
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal("home", err)
	}
	fmt.Println("You are logged in as:", user.ID)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, session_name)
	session.Options.MaxAge = -1
	err := session.Save(r, w)
	if err != nil {
		log.Fatal("failed to delete session", err)
	}
	http.Redirect(w, r, "/login", 301)
}