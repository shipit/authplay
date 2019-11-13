package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// GitHubAuth is start of OAuth webflow
func GitHubAuth(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("GHA") == "yes" {
		installAppFlow(w, r)
		return
	}
	oauthFlow(w, r)
}

// GitHubAuthCallback exchanges code for token
func GitHubAuthCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")
	log.Println("state: ", state, " requested scope: ", sessionData.StateMap[state])

	token, err := newGithubConf(nil).Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("token exchange failed: %#v\n", err)
		return
	}

	log.Println("github token: ", token, " granted scope: ", token.Extra("scope"))

	sessionData.AccessToken = token.AccessToken
	sessionData.Scope = token.Extra("scope").(string)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func newGithubConf(scope []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       scope,
		Endpoint:     githuboauth.Endpoint,
	}
}

func oauthFlow(w http.ResponseWriter, r *http.Request) {
	scopes := []string{"user:email", "repo"}
	githubConf := newGithubConf(scopes)

	state := "auth-" + uuid.Must(uuid.NewV4()).String()
	sessionData.StateMap[state] = strings.Join(scopes, ",")

	url := githubConf.AuthCodeURL(state, oauth2.AccessTypeOnline)
	log.Println("url: ", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func installAppFlow(w http.ResponseWriter, r *http.Request) {
	appName := os.Getenv("GITHUB_APP_NAME")
	state := uuid.Must(uuid.NewV4()).String()
	installURL := fmt.Sprintf("https://github.com/apps/%s/installations/new?state=%s", appName, state)
	log.Println("redirecting to: ", installURL)
	http.Redirect(w, r, installURL, http.StatusTemporaryRedirect)
}
