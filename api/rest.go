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

// GetInstalledRepos returns the repos the GitHub App is installed on
func GetInstalledRepos(token string) ([]*github.Repository, error) {
	ctx := context.Background()
	client, _ := newClient(token)
	list, _, err := client.Apps.ListInstallations(ctx, nil)
	if err != nil {
		log.Printf("%#v", err)
		return nil, err
	}

	var allRepos []*github.Repository
	for _, installation := range list {
		installID := installation.ID
		repos, _, err := client.Apps.ListUserRepos(ctx, *installID, nil)
		if err != nil {
			log.Printf("%#v", err)
			continue
		}
		allRepos = append(allRepos, repos...)
	}

	return allRepos, nil
}

func newClient(token string) (*github.Client, error) {
	ctx := context.Background()
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	http := oauth2.NewClient(ctx, src)

	return github.NewClient(http), nil
}
