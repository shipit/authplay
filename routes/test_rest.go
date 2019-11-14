package routes

import (
	"encoding/json"
	"log"

	"github.com/shipit/authplay/api"
)

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

func runInstalledReposTest(token string) *string {
	repos, err := api.GetInstalledRepos(token)
	if err != nil {
		log.Printf("installed repos error: %#v\n", err)
		return nil
	}

	sessionData.InstalledRepos = repos
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
