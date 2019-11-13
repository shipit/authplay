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

	runTests()

	err = t.ExecuteTemplate(w, "index", sessionData)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		log.Printf("%#v\n", err)
		return
	}
}

func runTests() {
	if len(sessionData.AccessToken) > 0 {
		user, err := api.GetUser(sessionData.AccessToken)
		if err != nil {
			log.Printf("user error: %#v\n", err)
			return
		}
		sessionData.User = user

		buf, err := json.Marshal(user)
		if err != nil {
			log.Printf("marshal error: %#v\n", err)
			return
		}

		str := string(buf)
		sessionData.UserJSON = &str
	}
}
