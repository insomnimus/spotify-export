package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/browser"
	"github.com/zmb3/spotify"
)

const state = "t0pk3k"

type Creds struct {
	Client *spotify.Client
	User   *spotify.PrivateUser
}

func authorize(addr string, auth *spotify.Authenticator, openBrowser bool) (*Creds, error) {
	ch := make(chan *spotify.Client)
	completeAuth := func(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/callback", completeAuth)
	go http.ListenAndServe(addr, nil)

	url := auth.AuthURL(state)
	if !openBrowser {
		fmt.Fprintln(os.Stderr, "Please visit this page on your browser to complete authentication:\n", url)
	} else {
		fmt.Fprintln(os.Stderr, "Opening below URL from your browser; if it does not open, visit it manually:\n", url)
		_ = browser.OpenURL(url)
	}
	// wait for auth to complete
	clt := <-ch
	// use the client to make calls that require authorization
	usr, err := clt.CurrentUser()
	if err != nil {
		return nil, fmt.Errorf("error auhorizing: %s", err)
	}

	return &Creds{
		Client: clt,
		User:   usr,
	}, nil
}

func login(c *Args) (*Creds, error) {
	auth := spotify.NewAuthenticator(c.redirect,
		spotify.ScopePlaylistReadPrivate,
		spotify.ScopePlaylistReadCollaborative,
		spotify.ScopeUserLibraryRead,
		spotify.ScopeUserReadPrivate,
		spotify.ScopeUserReadRecentlyPlayed,
		spotify.ScopeUserTopRead,
		spotify.ScopeStreaming,
	)
	auth.SetAuthInfo(c.id, c.secret)

	u, err := url.Parse(c.redirect)
	if err != nil {
		return nil, err
	}

	return authorize(u.Host, &auth, c.openBrowser)
}
