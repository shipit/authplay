package routes

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/shipit/authplay/api"
)

// Index is for /
func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index").Parse(sessionTemplate)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		log.Printf("%#v\n", err)
		return
	}

	if len(sessionData.AccessToken) > 0 {
		runTests(sessionData.AccessToken)
	}

	err = t.ExecuteTemplate(w, "index", sessionData)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		log.Printf("%#v\n", err)
		return
	}
}

type gqlResponse struct {
	User *api.GQLUser `json:"user"`
}

func GraphQL(w http.ResponseWriter, r *http.Request) {
	user, err := api.QueryUser(sessionData.AccessToken)
	if err != nil {
		log.Printf("user error: %#v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp := gqlResponse{
		User: user,
	}

	buf, err := json.Marshal(resp)
	if err != nil {
		log.Printf("marshal err: %#v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(buf)
}

func runTests(token string) {
	sessionData.UserJSON = runUserTest(token)
	sessionData.ReposJSON = runReposTest(token)
}
