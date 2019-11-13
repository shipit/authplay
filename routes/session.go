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
				<h4>User</h4>
				<ul>
					<li>Login: {{ .User.Login }}</li>
					<li>Name: {{ .User.Name }}</li>
					<li>Avatar URL: {{ .User.AvatarURL }}</li>
				</ul>
			</div>
			{{ if .UserJSON }}
				<code>{{ .UserJSON }}</code>
			{{ end }}
		{{ else }}
			<div>
				<form action="/auth" method="GET">
					<button>GitHub Login</button>
				</form>
			</div>
		{{end}}
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

	AccessToken string
	StateMap    map[string]string
}

var sessionData = SessionData{
	Title: "Don't Panic",
	Items: []string{
		"One thing",
		"And the other thing",
	},
	StateMap: map[string]string{},
}
