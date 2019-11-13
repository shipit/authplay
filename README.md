# Auth N Play

[![Discord](https://img.shields.io/discord/639591447569498133)](https://discord.gg/9AmwHdm)

Simple project to exercise `OAuth` web flows.

## Changelog

Starts server and listens on port `8888`

`/` returns login form or session data post OAuth
`/` also returns response for `/users` and `/user/repos/`

`/graphql` returns json encoded data for custom queries, only to be accessed after GitHub auth flow on `/`

## To Run

```bash
#!/bin/bash
cd test
GITHUB_CLIENT_ID=<client_id> GITHUB_CLIENT_SECRET=<client_secret> GITHUB_APP_ID=<app_id> GITHUB_APP_NAME=<app_name> GHA="yes|no" go run test_server.go
```

`GHA=yes` triggers GitHub App install and redirect to OAuth flow.

`CTRL-C` to kill the server, it captures `SIGINT` and cleans up.
