package api

import (
	"context"

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

func newClient(token string) (*github.Client, error) {
	ctx := context.Background()
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	http := oauth2.NewClient(ctx, src)

	return github.NewClient(http), nil
}
