package api

import (
	"context"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// GQLUser is the user object with fields we are interested in
type GQLUser struct {
	DatabaseID int64  `json:"id"`
	AvatarURL  string `json:"avatar"`
	Login      string `json:"login"`
	Name       string `json:"name"`
}

// QueryUser runs graphql query for user object
func QueryUser(token string) (*GQLUser, error) {
	var query struct {
		Viewer struct {
			GQLUser
		}
	}

	ctx := context.Background()
	client, _ := newGQLClient(token)
	err := client.Query(ctx, &query, nil)
	if err != nil {
		return nil, err
	}

	return &query.Viewer.GQLUser, nil
}

func newGQLClient(token string) (*githubv4.Client, error) {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	http := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(http)

	return client, nil
}
