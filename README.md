# Auth N Play

Simple project to exercise `OAuth` web flows.

## Changelog

Starts server and listens on port `8888`, returns simple string for `/` path. It's the only path atm.

## To Run

```bash
cd test
GITHUB_CLIENT_ID=<client_id> GITHUB_CLIENT_SECRET=<client_secret> go run test_server.go
```

`CTRL-C` to kill the server, it captures `SIGINT` and cleans up.
