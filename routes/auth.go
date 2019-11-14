package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// GitHubAuth is start of OAuth webflow
func GitHubAuth(w http.ResponseWriter, r *http.Request) {
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

// GitHubAppInstall redirects to GitHub for app installation on repos
func GitHubAppInstall(w http.ResponseWriter, r *http.Request) {
	installAppFlow(w, r)
}

// GitHubAppPostInstall handles post GitHub App installation, should get install_id
func GitHubAppPostInstall(w http.ResponseWriter, r *http.Request) {
	// ?installation_id=5164770&setup_action=install&state=99234ba4-5f2b-40b2-9977-89172e98cd30
	action := r.FormValue("setup_action")
	id := r.FormValue("installation_id")
	log.Println("action: ", action, " id: ", id)

	if action == "install" || action == "update" {
		installID, _ := strconv.ParseInt(id, 10, 64)
		sessionData.InstallID = &installID
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// GitHubAppWebhook handles webhook events from GitHub
func GitHubAppWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%#v", err)
		return
	}
	defer r.Body.Close()

	t := github.WebHookType(r)
	switch t {
	case "installation":
		handleGitHubAppInstallation(payload)
	default:
		log.Println("unhandled webhook: ", string(payload))
	}
}

// GitHubAppInstallationEvent is payload for installation webhook event
type GitHubAppInstallationEvent struct {
	Installation struct {
		ID *int64 `json:"id"`
	} `json:"installation"`
	Action       *string              `json:"action"`
	Repositories []*github.Repository `json:"repositories"`
	Sender       *github.User         `json:"sender"`
}

func handleGitHubAppInstallation(payload []byte) {
	var ie GitHubAppInstallationEvent
	err := json.Unmarshal(payload, &ie)
	if err != nil {
		log.Println("!! handleGitHubAppInstallation() payload unmarshal error: ", err.Error())
		return
	}

	switch *ie.Action {
	case "created", "update":
		sessionData.InstallID = ie.Installation.ID
	case "deleted":
		sessionData.InstallID = nil
	default:
		log.Println("!! handleGitHubAppInstallation() unhandled github installation event: ", ie)
	}
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
