package main

import (
	"fmt"
	"net/http"
	"log"
	//"os"
)

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

	// the client can now be used to make authenticated requests
	//ch <- &client
	session, _ := store.Get(r, "cookie-name")
	session.Values["token"] = Oauth2Token{*token}
	//fmt.Fprintf(os.Stdout, "session : %#v\n", session.Values["token"].(Oauth2Token).flag)
	err = session.Save(r, w)
	fmt.Println("err: ", err)

	http.Redirect(w, r, "/home", 301)
}

