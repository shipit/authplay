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
		{{range .Items}}
		<div>
			{{ . }}
		</div>
		{{ else }}
		<div>
			<strong>no items</strong>
		</div>
		{{ end }}
	</body>
</html>
`

// SessionData to transform index template
type SessionData struct {
	Title string
	Items []string

	User     *github.User
	UserJSON *string

	Repos        []*github.Repository
	ReposJSON    *string
	PublicRepos  int
	PrivateRepos int

	AccessToken string
	StateMap    map[string]string
	Scope       string
}

var sessionData = SessionData{
	Title: "Don't Panic",
	Items: []string{
		"One thing",
		"And the other thing",
	},
	StateMap: map[string]string{},
}
