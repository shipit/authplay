package routes

import "github.com/google/go-github/github"

const sessionTemplate = `
<!doctype html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{ .Title }}</title>
	</head>
	<body>
		{{ if .User }}
			<div>
				<h4>Session</h4>
				<ul>
					<li>Login: {{ .User.Login }}</li>
					<li>Name: {{ .User.Name }}</li>
					<li>Avatar URL: {{ .User.AvatarURL }}</li>
					<li>Token: {{ .AccessToken }}</li>
					<li>Scope: {{ .Scope }}</li>
				</ul>
				<div>
					<a href="/graphql" target="_blank">Run GraphQL</a>
				</div>
				<div>
					<a href="/install_app">Add Repos</a>
				</div>
			</div>
		{{ else }}
			<div>
				<form action="/auth" method="GET">
					<button>GitHub Login</button>
				</form>
			</div>
		{{end}}
		{{ if .UserJSON }}
			<h5>User</h5>
			<code>{{ .UserJSON }}</code>
		{{ end }}
		{{ if .InstalledReposJSON }}
			<h5>Installed Repos</h5>
			<code>{{ .InstalledReposJSON }}</code>
		{{ end }}
		{{ if .ReposJSON }}
			<h5>Repos</h5>
			{{ if .Repos }}
				<ul>
					<li>Count: {{ len .Repos }}</li>
					<li>Private: {{ .PrivateRepos }}</li>
					<li>Public: {{ .PublicRepos }}</li>
				</ul>
			{{ end }}
			<code>{{ .ReposJSON }}</code>
		{{ end }}
	</body>
</html>
`

// SessionData to transform index template
type SessionData struct {
	Title string `json:"-"`

	User     *github.User `json:"user"`
	UserJSON *string      `json:"-"`

	Repos        []*github.Repository `json:"repos"`
	ReposJSON    *string              `json:"-"`
	PublicRepos  int                  `json:"-"`
	PrivateRepos int                  `json:"-"`

	AccessToken string            `json:"access_token"`
	StateMap    map[string]string `json:"state_map"`
	Scope       string            `json:"scope"`

	InstallID          *int64               `json:"install_id"`
	InstalledRepos     []*github.Repository `json:"installed_repos"`
	InstalledReposJSON *string              `json:"-"`
}

var sessionData = SessionData{
	Title:    "Don't Panic",
	StateMap: map[string]string{},
}
