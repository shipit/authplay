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

func runTests(token string) {
	sessionData.UserJSON = runUserTest(token)
	sessionData.ReposJSON = runReposTest(token)
}

func runUserTest(token string) *string {
	user, err := api.GetUser(token)
	if err != nil {
		log.Printf("user error: %#v\n", err)
		return nil
	}
	sessionData.User = user
	return marshal(user)
}

func runReposTest(token string) *string {
	repos, err := api.GetRepos(token)
	if err != nil {
		log.Printf("repos error: %#v\n", err)
		return nil
	}

	sessionData.Repos = repos

	var public, private = 0, 0
	for _, repo := range repos {
		if *repo.Private {
			private++
		} else {
			public++
		}
	}

	sessionData.PublicRepos = public
	sessionData.PrivateRepos = private

	return marshal(repos)
}

func marshal(in interface{}) *string {
	buf, err := json.Marshal(in)
	if err != nil {
		log.Printf("marshal error: %#v", err)
		return nil
	}

	str := string(buf)
	return &str
}
