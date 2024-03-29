package api

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

// GQLNode is Node object
type GQLNode struct {
	Name string `json:"name"`
}

// GQLRepo is Repository object
type GQLRepo struct {
	GQLNode
	DatabaseID    int64  `json:"id"`
	NameWithOwner string `json:"nameWithOwner"`
	IsPrivate     bool   `json:"is_private"`
	IsFork        bool   `json:"is_fork"`
	Description   string `json:"description"`

	Owner struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	} `json:"owner"`

	Languages struct {
		Nodes []*GQLNode `json:"nodes"`
	} `graphql:"languages(first: 4)" json:"languages"`

	Collaborators struct {
		Nodes []GQLUser `json:"nodes"`
	} `json:"collaborators"`

	DefaultBranchRef struct {
		Name   string `json:"name"`
		Target struct {
			Commit struct {
				CommittedDate string `json:"committed_date"`
				OID           string `json:"commit_sha"`
				Message       string `json:"message"`
				Author        struct {
					User struct {
						Login string `json:"login"`
						Name  string `json:"name"`
					} `json:"user"`
				} `json:"author"`
				Tree struct {
					Entries []struct {
						GQLNode
						Type string `json:"type"`
					} `json:"entries"`
				} `json:"tree"`
			} `graphql:"... on Commit" json:"commit"`
		} `json:"target"`
	} `json:"default_branch"`
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

// QueryRepoByName fetches repo details, for an installation the repo must be added to it
func QueryRepoByName(token string, repo string) (*GQLRepo, error) {
	var query struct {
		Viewer struct {
			GQLRepo `graphql:"repository(name: $repoName)"`
		}
	}

	variables := map[string]interface{}{
		"repoName": githubv4.String(repo),
	}

	ctx := context.Background()
	client, _ := newGQLClient(token)
	err := client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	return &query.Viewer.GQLRepo, nil
}

// QueryRepos runs custom query for repos and their commits, collaborators, and root tree on default branch
func QueryRepos(token string) ([]GQLRepo, error) {
	var query struct {
		Viewer struct {
			Repositories struct {
				Nodes    []GQLRepo
				PageInfo struct {
					HasNextPage bool
					EndCursor   githubv4.String
				}
			} `graphql:"repositories(first: 100, after: $reposCursor, orderBy:{field: NAME, direction: ASC})"`
		}
	}

	variables := map[string]interface{}{
		"reposCursor": (*githubv4.String)(nil),
	}

	ctx := context.Background()
	client, _ := newGQLClient(token)

	var allRepos []GQLRepo
	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, err
		}

		for _, repo := range query.Viewer.Repositories.Nodes {
			allRepos = append(allRepos, repo)
		}

		if !query.Viewer.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["reposCursor"] = githubv4.NewString(query.Viewer.Repositories.PageInfo.EndCursor)
	}

	return allRepos, nil
}

type wrapped struct {
	base http.RoundTripper
}

func (w wrapped) RoundTrip(r *http.Request) (*http.Response, error) {
	log.Print(r.URL)
	resp, err := w.base.RoundTrip(r)
	if err != nil {
		return resp, err
	}

	reader := io.TeeReader(resp.Body, os.Stderr)
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return resp, err
	}
	log.Println(string(buf))

	resp.Body = ioutil.NopCloser(bytes.NewReader(buf))
	return w.base.RoundTrip(r)
}

func newGQLClient(token string) (*githubv4.Client, error) {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})

	http := oauth2.NewClient(context.Background(), src)
	http.Transport = wrapped{base: http.Transport}

	client := githubv4.NewClient(http)

	return client, nil
}
