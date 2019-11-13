# Auth N Play

[![Discord](https://img.shields.io/discord/586999333447270440.svg)](https://discord.gg/9AmwHdm)

Simple project to exercise `OAuth` web flows.

## Changelog

Starts server and listens on port `8888`

`/` returns login form or session data post OAuth
`/` also returns response for `/users` and `/user/repos/`

`/graphql` returns json encoded data for custom queries, only to be accessed after GitHub auth flow on `/`

## To Run

```bash
cd test
GITHUB_CLIENT_ID=<client_id> GITHUB_CLIENT_SECRET=<client_secret> go run test_server.go
```

`CTRL-C` to kill the server, it captures `SIGINT` and cleans up.
