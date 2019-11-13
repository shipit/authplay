package api

import (
	"context"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GetUser returns GitHub user
func GetUser(token string) (*github.User, error) {
	client, err := newClient(token)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetRepos returns list of repos for the authed user
func GetRepos(token string) ([]*github.Repository, error) {
	client, err := newClient(token)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	repos, _, err := client.Repositories.List(ctx, "", &github.RepositoryListOptions{ListOptions: github.ListOptions{PerPage: 200}})
	if err != nil {
		log.Fatalf("%#v", err)
	}

	return repos, err
}

func newClient(token string) (*github.Client, error) {
	ctx := context.Background()
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	http := oauth2.NewClient(ctx, src)

	return github.NewClient(http), nil
}
