package routes

import (
	"html/template"
	"log"
	"net/http"
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

	err = t.ExecuteTemplate(w, "index", sessionData)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		log.Printf("%#v\n", err)
		return
	}
}
